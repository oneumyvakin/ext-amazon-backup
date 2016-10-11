package amazons3

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AmazonS3 struct {
	Log    *log.Logger
	Svc    *s3.S3
	Region string
	Bucket string
}

type downloader struct {
	AmazonS3

	File       *os.File
	FileOffset int64
	Err        error
}

type filePart struct {
	Key        string
	Range      string
	Etag       string
	Offset     int64
	Length     int64
	PartNumber int64
	Body       []byte
}

func (self AmazonS3) GetRegions() []string {
	var AwsRegions = []string{
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"eu-west-1",
		"eu-central-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"sa-east-1",
	}
	return AwsRegions
}

func (self AmazonS3) IsRegionValid(name string) error {
	regions := self.GetRegions()
	sort.Strings(regions)
	i := sort.SearchStrings(regions, name)
	if i < len(regions) && regions[i] == name {
		self.Log.Println("Region valid:", name)
		return nil
	}

	return fmt.Errorf("Failed to validate region: %s", name)
}

// Get Bucket's region
func (self AmazonS3) GetBucketLocation(name string) (region string, err error) {
	params := &s3.GetBucketLocationInput{
		Bucket: aws.String(name), // Required
	}
	self.Log.Printf("Get Bucket location: %s\n", name)
	resp, err := self.Svc.GetBucketLocation(params)
	if err != nil {
		return "", err
	}
	region = *resp.LocationConstraint

	self.Log.Printf("Bucket %s region: %s", name, region)
	return
}

func (self AmazonS3) CreateBucket(name string) error {
	buckets, err := self.GetBucketsList()
	if err != nil {
		return err
	}

	sort.Strings(buckets)
	i := sort.SearchStrings(buckets, name)
	if i < len(buckets) && buckets[i] == name {
		self.Log.Println("Bucket already exists:", name)
		return nil
	}

	_, err = self.Svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: &name,
	})
	if err != nil {
		self.Log.Printf("Failed to create bucket %s: %s", name, err)
		return err
	}

	if err = self.Svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &name}); err != nil {
		self.Log.Printf("Failed to wait for bucket to exist %s: %s\n", name, err)
		return err
	}

	self.Log.Println("Create bucket:", name)
	return nil
}

func (self AmazonS3) CreateFolder(path string) error {
	req := &s3.PutObjectInput{
		Bucket: aws.String(self.Bucket),
		Key:    aws.String(path + "/"),
	}
	_, err := self.Svc.PutObject(req)

	return err
}

// List available buckets
func (self AmazonS3) GetBucketsList() (list []string, err error) {
	result, err := self.Svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		self.Log.Printf("Failed to list buckets: %s\n", err)
		return
	}

	for _, bucket := range result.Buckets {
		list = append(list, *bucket.Name)
	}

	self.Log.Println("Get buckets:", list)
	return
}

// List files and folders.
// SubFolder can be ""
func (self AmazonS3) GetBucketFilesList(subFolder string) ([]*s3.Object, error) {
	subFolder = strings.TrimSuffix(subFolder, "/")
	result, err := self.Svc.ListObjects(&s3.ListObjectsInput{Bucket: &self.Bucket, Prefix: &subFolder})
	if err != nil {
		self.Log.Printf("Failed to list objects: %s\n", err)
		return nil, err
	}

	self.Log.Printf("Get bucket files in /%s:\n", subFolder, result.Contents)
	return result.Contents, nil
}

// Get file info
// Returns http://docs.aws.amazon.com/sdk-for-go/api/service/s3/#HeadObjectOutput
func (self AmazonS3) GetFileInfo(path string) (resp *s3.HeadObjectOutput, err error) {
	resp, err = self.Svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(self.Bucket), // Required
		Key:    aws.String(path),        // Required
	})

	if err != nil {
		self.Log.Printf("Failed to get file %s info error: %s\n", path, err)
		return
	}
	self.Log.Println("Get file info:", path, resp)

	//_, _ = self.GetFilePartInfo(path, fmt.Sprintf("bytes=%d-%d", 0, s3manager.DefaultUploadPartSize))
	return
}

// Get File Part
// partRange "bytes=0-100"
func (self AmazonS3) GetFilePart(path string, partRange string) (resp *s3.GetObjectOutput, err error) {
	resp, err = self.Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(self.Bucket), // Required
		Key:    aws.String(path),        // Required
		Range:  aws.String(partRange),
	})

	if err != nil {
		self.Log.Printf("Failed to get file %s part: %s\n", path, err)
		return
	}
	self.Log.Println("Get file part:", path, resp)
	return
}

func (self AmazonS3) Delete(path string) (err error) {
	resp, err := self.Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(self.Bucket), // Required
		Key:    aws.String(path),        // Required
	})

	if err != nil {
		self.Log.Println("Failed to delete:", path, err)
		return
	}
	self.Log.Println("Delete path:", path, resp)
	return
}

// List bucket's unfinished uploads
// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/#MultipartUpload
func (self AmazonS3) ListUnfinishedUploads() ([]*s3.MultipartUpload, error) {
	resp, err := self.Svc.ListMultipartUploads(&s3.ListMultipartUploadsInput{
		Bucket: aws.String(self.Bucket), // Required
	})
	if err != nil {
		self.Log.Printf("Failed list unfinised uploads: %s\n", err)
		return nil, err
	}
	self.Log.Println("List bucket's unfinished uploads", resp)

	return resp.Uploads, nil
}

// List parts of unfinished uploads
// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/#ListPartsOutput
// Parts []*Part - can be empty
// Part.PartNumber
// Part.Size
func (self AmazonS3) ListParts(key string, uploadId string) (resp *s3.ListPartsOutput, err error) {
	resp, err = self.Svc.ListParts(&s3.ListPartsInput{
		Bucket:   aws.String(self.Bucket), // Required
		Key:      aws.String(key),         // Required
		UploadId: aws.String(uploadId),    // Required
	})
	self.Log.Printf("List parts for key %s of upload id %s: %s\n", key, uploadId, resp)

	return
}

// Abort upload
func (self AmazonS3) AbortUpload(key string, uploadId string) (err error) {
	resp, err := self.Svc.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   aws.String(self.Bucket), // Required
		Key:      aws.String(key),         // Required
		UploadId: aws.String(uploadId),    // Required
	})
	self.Log.Printf("Abort upload for key %s of upload id %s: %s\n", key, uploadId, resp)

	return
}

// Complete upload
func (self AmazonS3) CompleteUpload(key string, uploadId string) (err error) {
	respParts, err := self.ListParts(key, uploadId) // Just for debug
	if err != nil {
		self.Log.Printf("Failed to complete upload: Failed to list parts for key %s of upload id %s: %s\n", key, uploadId, err)
		return
	}

	var completedParts []*s3.CompletedPart
	for _, part := range respParts.Parts {
		completedPart := &s3.CompletedPart{
			ETag:       part.ETag,
			PartNumber: part.PartNumber,
		}
		completedParts = append(completedParts, completedPart)
	}
	resp, err := self.Svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(self.Bucket), // Required
		Key:      aws.String(key),         // Required
		UploadId: aws.String(uploadId),    // Required
		MultipartUpload: &s3.CompletedMultipartUpload{ // Required
			Parts: completedParts,
		},
	})
	if err != nil {
		self.Log.Printf("Failed to complete upload for key %s of upload id %s: %s\n", key, uploadId, err)
		return
	}
	self.Log.Printf("Complete upload for key %s of upload id %s: %s\n", key, uploadId, resp)

	return
}

func (self AmazonS3) IoClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		self.Log.Println(err)
	}
}
