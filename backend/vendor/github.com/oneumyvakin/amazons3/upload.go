package amazons3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"crypto/md5"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Upload filePath to destinationPath, where destinationPath contains only folders like /folder/folder2
func (self AmazonS3) Upload(filePath, destinationPath string, useGzip bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open file %s for upload: %s\n", filePath, err)
	}

	key := destinationPath + filepath.Base(filePath)

	// Not required, but you could zip the file before uploading it
	// using io.Pipe read/writer to stream gzip'd file contents.
	reader, writer := io.Pipe()

	if useGzip {
		go func() {
			gw := gzip.NewWriter(writer)
			written, err := io.Copy(gw, file)
			if err != nil {
				self.Log.Printf("AmazonS3 Upload gzip io.Copy error: %s\n", err)
			}
			self.Log.Printf("AmazonS3 Upload gzip io.Copy written: %s\n", written)

			self.IoClose(file)
			self.IoClose(gw)
			self.IoClose(writer)
		}()

		key = key + ".gz"
	} else {
		go func() {
			bw := bufio.NewWriter(writer)
			written, err := io.Copy(bw, file)
			if err != nil {
				self.Log.Printf("AmazonS3 Upload buffer io.Copy error: %s\n", err)
			}
			self.Log.Printf("AmazonS3 Upload buffer io.Copy written: %s\n", written)

			self.IoClose(file)
			err = bw.Flush()
			if err != nil {
				self.Log.Printf("bufio flush error: %s\n", err)
			}
			self.IoClose(writer)
		}()
	}

	self.Log.Printf("Upload %s to %s with Gzip: %t\n", filePath, key, useGzip)

	uploader := s3manager.NewUploader(
		session.New(
			&aws.Config{
				Region: aws.String(self.Region),
			},
		),
		func(u *s3manager.Uploader) {
			u.LeavePartsOnError = true // Leave good uploaded parts on Storage in case of failure
		},
	)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(self.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if multiErr, ok := err.(s3manager.MultiUploadFailure); ok {
			// Process error and its associated uploadID
			self.Log.Printf("Error code: %s, Message: %s, UploadID: %s\n", multiErr.Code(), multiErr.Message(), multiErr.UploadID())
			return multiErr
		}
		return fmt.Errorf("Failed upload file %s: %s\n", filePath, err)
	}

	self.Log.Println("Successfully uploaded to", result.Location)
	return nil
}

func (self AmazonS3) ResumeUpload(filePath, key, uploadId string, useGzip bool) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open file %s for upload: %s\n", filePath, err)
	}

	// Not required, but you could zip the file before uploading it
	// using io.Pipe read/writer to stream gzip'd file contents.
	pipeReader, writer := io.Pipe()

	if useGzip {
		go func() {
			gw := gzip.NewWriter(writer)
			written, err := io.Copy(gw, file)
			if err != nil {
				self.Log.Printf("AmazonS3 Upload gzip io.Copy error: %s\n", err)
			}
			self.Log.Printf("AmazonS3 Upload gzip io.Copy written: %s\n", written)

			self.IoClose(file)
			self.IoClose(gw)
			self.IoClose(writer)
		}()
	} else {
		go func() {
			bw := bufio.NewWriter(writer)
			written, err := io.Copy(bw, file)
			if err != nil {
				self.Log.Printf("AmazonS3 Upload buffer io.Copy error: %s\n", err)
			}
			self.Log.Printf("AmazonS3 Upload buffer io.Copy written: %s\n", written)

			self.IoClose(file)
			err = bw.Flush()
			if err != nil {
				self.Log.Printf("bufio flush error: %s\n", err)
			}
			self.IoClose(writer)
		}()
	}

	self.Log.Printf("Resume Upload %s to %s with Gzip: %t\n", filePath, key, useGzip)

	resp, err := self.ListParts(key, uploadId)
	if err != nil {
		return fmt.Errorf("Failed to list uploaded parts for key %s of upload id %s: %s\n", key, uploadId, err)
	}

	partQueue := make(chan filePart, s3manager.DefaultUploadConcurrency)
	var wg sync.WaitGroup

	for i := 0; i < s3manager.DefaultUploadConcurrency; i++ {
		wg.Add(1)
		go self.asyncUploadPart(key, uploadId, partQueue, &wg)
	}

	go self.getFileParts(partQueue, pipeReader, resp.Parts)

	self.Log.Println("Wait for all parts are uploading...")
	wg.Wait()

	err = self.CompleteUpload(key, uploadId)
	if err != nil {
		return fmt.Errorf("Failed to complete upload with key %s: %s\n", key, err)
	}

	self.Log.Println("Successfully resumed upload to", key)

	return nil
}

func (self AmazonS3) getFileParts(partChan chan<- filePart, reader io.Reader, uploadedParts []*s3.Part) {
	var lastPartNumber int64
	var offset int64
	lastPartNumber = 1
	offset = 0

	for {
		part := make([]byte, s3manager.DefaultUploadPartSize)
		partSize, errRead := io.ReadFull(reader, part)
		if errRead != nil && errRead != io.EOF && errRead != io.ErrUnexpectedEOF {
			self.Log.Fatalf("Failed to read part number %s from reader at offset %s: %s\n", lastPartNumber, offset, errRead)
		}
		self.Log.Printf("Read bytes %s for part number %d with size: %s\n", partSize, lastPartNumber, len(part))

		if int64(partSize) != s3manager.DefaultUploadPartSize {
			lastPart := make([]byte, partSize)
			copy(lastPart, part)
			part = lastPart
		}

		partEtag, err := self.getPartEtag(part)
		if err != nil {
			self.Log.Fatalf("Failed to get Etag for part number %s with size %s: %s\n", lastPartNumber, partSize, err)
		}

		self.Log.Printf("Part number %s size bytes %s has ETag: %s\n", lastPartNumber, len(part), partEtag)

		if true == self.needToUpload(uploadedParts, lastPartNumber, partEtag) {
			partChan <- filePart{
				Body:       part,
				PartNumber: lastPartNumber,
			}
		}

		offset = offset + int64(len(part))
		lastPartNumber = lastPartNumber + 1

		if errRead == io.EOF || errRead == io.ErrUnexpectedEOF {
			self.Log.Printf("EOF or ErrUnexpectedEOF. All parts are read and send to upload. Last part is %s, offset is %d", lastPartNumber, offset)
			close(partChan)
			return
		}
	}
}

func (self AmazonS3) needToUpload(uploadedParts []*s3.Part, partNumber int64, partEtag string) bool {
	for _, part := range uploadedParts {
		if *part.PartNumber == partNumber {
			self.Log.Printf("Part number %s with ETag %s found\n", *part.PartNumber, string(*part.ETag))

			if *part.ETag == partEtag {
				self.Log.Printf("Match Etag for part number %s with size %s ETag %s == %s.\n", *part.PartNumber, *part.Size, string(*part.ETag), partEtag)
				return false
			} else {
				self.Log.Printf("Mismatch Etag for part number %s with size %s ETag %s != %s. Reuploading...\n", *part.PartNumber, *part.Size, string(*part.ETag), partEtag)
				return true
			}
		}
	}
	self.Log.Printf("Part number %s not found\n", partNumber)

	return true
}

func (self AmazonS3) asyncUploadPart(key string, uploadId string, partChan <-chan filePart, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if part, ok := <-partChan; ok {
			self.Log.Printf("Start to upload part number %s for key %s\n", part.PartNumber, key)

			_, err := self.Svc.UploadPart(&s3.UploadPartInput{
				Bucket:     aws.String(self.Bucket),    // Required
				Key:        aws.String(key),            // Required
				PartNumber: aws.Int64(part.PartNumber), // Required
				UploadId:   aws.String(uploadId),       // Required
				Body:       bytes.NewReader(part.Body),
			})
			if err != nil {
				self.Log.Printf("Failed to upload part number %s for key %s: %s\n", part.PartNumber, key, err)
				return
			}
			self.Log.Printf("Finished upload part number %s for key %s\n", part.PartNumber, key)
		} else {
			self.Log.Println("Upload channel closed. Return.")

			return
		}
	}

	return
}

func (self AmazonS3) uploadPart(key string, partNumber int64, uploadId string, body []byte) (err error) {
	self.Log.Printf("Start upload part number %d of key %s for upload id %s\n", partNumber, key, uploadId)

	_, err = self.Svc.UploadPart(&s3.UploadPartInput{
		Bucket:     aws.String(self.Bucket), // Required
		Key:        aws.String(key),         // Required
		PartNumber: aws.Int64(partNumber),   // Required
		UploadId:   aws.String(uploadId),    // Required
		Body:       bytes.NewReader(body),
	})

	return
}

func (self AmazonS3) getPartEtag(part []byte) (etag string, err error) {
	hasher := md5.New()
	_, err = hasher.Write(part)
	if err != nil {
		self.Log.Printf("Failed to write part to hasher: %s", err)
		return
	}
	etag = fmt.Sprintf("\"%s\"", hex.EncodeToString(hasher.Sum(nil)))

	return
}
