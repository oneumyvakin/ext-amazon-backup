package amazons3

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (self AmazonS3) Download(fileName, destinationPath string) error {
	fileInfo, err := self.GetFileInfo(fileName)
	if err != nil {
		return fmt.Errorf("Failed to get file %s: %s\n", fileName, err)
	}
	var file *os.File

	if _, err = os.Stat(destinationPath); os.IsNotExist(err) {
		file, err = self.createEmptyFile(destinationPath, *fileInfo.ContentLength)
		if err != nil {
			return fmt.Errorf("Failed to allocate %s bytes on disk for destination file %s: %s\n", *fileInfo.ContentLength, destinationPath, err)
		}
		self.IoClose(file)
		err = os.Remove(destinationPath)
		if err != nil {
			return fmt.Errorf("Failed to remove temporary file %s: %s\n", destinationPath, err)
		}
		file, err = os.Create(destinationPath)
		if err != nil {
			return fmt.Errorf("Failed to create destination file %s: %s\n", destinationPath, err)
		}

	} else {
		file, err = os.OpenFile(destinationPath, os.O_WRONLY, 666)
		if err != nil {
			return fmt.Errorf("Failed to create destination file %s: %s\n", destinationPath, err)
		}
	}
	defer self.IoClose(file)

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(self.Region)}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(self.Bucket),
			Key:    aws.String(fileName),
		})
	if err != nil {
		self.Log.Printf("Failed to download file %s to destination path %s with error: %s\n", fileName, destinationPath, err)
		return err
	}

	self.Log.Printf("Downloaded file %s with size %s\n", file.Name(), numBytes)
	return nil
}

func (self AmazonS3) ResumeDownload(fileName, destinationPath string) error {
	remoteFileInfo, err := self.GetFileInfo(fileName)
	if err != nil {
		self.Log.Printf("Failed to get file %s: %s\n", fileName, err)
		return err
	}

	file, err := os.OpenFile(destinationPath, os.O_WRONLY, 666)
	if err != nil {
		return fmt.Errorf("Failed to create destination file %s: %s\n", destinationPath, err)
	}
	defer self.IoClose(file)

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Failed to stat destination file %s: %s\n", destinationPath, err)
	}

	if *remoteFileInfo.ContentLength < stat.Size() {
		return fmt.Errorf("Failed to compare size of remote %s and destination file %s: %d <= %d\n", fileName, destinationPath, *remoteFileInfo.ContentLength, stat.Size())
	}

	if *remoteFileInfo.ContentLength == stat.Size() {
		self.Log.Printf("Size of remote %s and destination %s file match: %d == %d. Nothing to do.\n", fileName, destinationPath, *remoteFileInfo.ContentLength, stat.Size())
		return nil
	}

	d := downloader{
		AmazonS3:   self,
		File:       file,
		FileOffset: stat.Size(),
	}

	taskPartChan := make(chan filePart, s3manager.DefaultDownloadConcurrency)
	var wg sync.WaitGroup
	for i := 0; i < s3manager.DefaultUploadConcurrency; i++ {
		wg.Add(1)
		go d.asyncDownloadPart(taskPartChan, &wg)
	}

	partOffset := stat.Size()
	leftBytes := *remoteFileInfo.ContentLength - stat.Size()
	go func() {
		for {
			self.Log.Printf("Resume download: Left bytes %d\n", leftBytes)
			if leftBytes <= s3manager.DefaultDownloadPartSize {
				partRange := fmt.Sprintf("bytes=%d-%d", partOffset, partOffset+leftBytes-1)
				self.Log.Printf("Resume download: File range %s\n", partRange)
				taskPartChan <- filePart{
					Key:    fileName,
					Range:  partRange,
					Offset: partOffset,
					Length: leftBytes,
					Body:   make([]byte, leftBytes),
				}
				close(taskPartChan)
				self.Log.Println("Resume download: All parts send to download. Close channel.")
				return
			}
			fileRange := fmt.Sprintf("bytes=%d-%d", partOffset, partOffset+s3manager.DefaultDownloadPartSize-1)
			self.Log.Printf("Resume download: Part range %s\n", fileRange)
			self.Log.Printf("Resume download: Part offset %d\n", partOffset)

			taskPartChan <- filePart{
				Key:    fileName,
				Range:  fileRange,
				Offset: partOffset,
				Length: s3manager.DefaultDownloadPartSize,
			}
			partOffset = partOffset + s3manager.DefaultDownloadPartSize
			leftBytes = leftBytes - s3manager.DefaultDownloadPartSize
		}
	}()

	wg.Wait()

	return nil
}

func (self *downloader) asyncDownloadPart(taskPartChan <-chan filePart, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if part, ok := <-taskPartChan; ok {
			if self.Err != nil {
				self.Log.Printf("Failed to start download: %s\n", self.Err)
				return
			}
			self.Log.Printf("Start to download part for key %s: Range: %s, Offset: %d, Length: %d\n", part.Key, part.Range, part.Offset, part.Length)

			resp, err := self.Svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(self.Bucket), // Required
				Key:    aws.String(part.Key),    // Required
				Range:  aws.String(fmt.Sprintf("bytes=%d-%d", part.Offset, part.Offset+part.Length-1)),
			})
			self.Log.Printf("Request sent for %s range %s\n", part.Key, part.Range)
			if err != nil {
				self.Log.Printf("Failed to download file %s range %s: %s\n", part.Key, part.Range, err)
				return
			}
			self.Log.Printf("Response for %s range %s: %s\n", part.Key, part.Range, resp)
			self.Log.Printf("File offset: %d\n", self.FileOffset)
			self.Log.Printf("Part offset: %d\n", part.Offset)
			defer self.IoClose(resp.Body)

			for {
				if self.FileOffset == part.Offset {
					n, err := io.Copy(self.File, resp.Body)
					if err != nil {
						self.Err = err
						self.Log.Printf("Failed to write file %s range %s: %s\n", part.Key, part.Range, err)
						return
					}
					self.Log.Printf("Finish write %d bytes part range %s for key %s \n", n, part.Range, part.Key)

					self.FileOffset = part.Offset + part.Length
					self.Log.Printf("New file offset: %d\n", self.FileOffset)
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			self.Log.Println("Download channel closed. Return.")

			return
		}
	}
}

func (self AmazonS3) createEmptyFile(filePath string, size int64) (f *os.File, err error) {
	f, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}

	chunkSize := int64(1024 * 1024 * 25)

	self.Log.Printf("Start creating empty file %s with size %s\n", filePath, size)
	for {
		if size <= chunkSize {
			s := make([]byte, size)
			n, err := f.Write(s)
			self.Log.Printf("Bytes written %s\n", n)
			_, err = f.Seek(0, 0)
			return f, err
		}

		size = size - chunkSize

		s := make([]byte, chunkSize)

		n, err := f.Write(s)
		if err != nil {
			return nil, err
		}
		self.Log.Printf("Bytes written %s\n", n)
	}
}
