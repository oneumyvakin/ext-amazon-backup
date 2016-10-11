// Copyright 1999-2016. Parallels IP Holdings GmbH.

package main

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/oneumyvakin/amazons3"
	"github.com/oneumyvakin/ntp"
	"github.com/oneumyvakin/osext"
	"github.com/oneumyvakin/plesk"
)

var binaryDir string
var logger *log.Logger

const (
	backupOnAmazon = "backup-on-amazon"
)

type BackupToAmazon struct {
	Log               *log.Logger
	Plesk             plesk.Plesk
	AwsS3             amazons3.AmazonS3
	Backups           map[string]Backup
	UseGzip           bool
	SubFolder         string
	HttpStatus        map[string]string
	HttpStatusSetter  bool
	HttpHost          string
	HttpPort          string
	UploadDir         string
	DownloadDir       string
	BackupFileExt     string
	SkipIncremental   bool
	DeleteAfterUpload bool
}

type Backup struct {
	Name              string
	InProgress        bool
	IsRemote          bool
	RemotePath        string
	IsPartialUpload   bool
	IsPartialDownload bool
	UploadId          string
	IsLocal           bool
	IsLocalInvalid    bool
	Backup            plesk.Dump
}

type UnfinishedUpload struct {
	RemotePath   string
	UploadId     string
	Size         int64
	LastModified time.Time
}

type ErrCoder interface {
	Code() string
}

// Application level error
type Err struct {
	IsError    bool              `json:"is_error"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	UploadID   string            `json:"upload_id"`
	LocaleKey  string            `json:"locale_key"`
	LocaleArgs map[string]string `json:"locale_args"`
	OriginErr  error             `json:"origin_err"`
}

func (e Err) Error() string {
	return e.Message
}

type AwsCustomLogger struct {
	logger *log.Logger
}

func (awc AwsCustomLogger) Log(args ...interface{}) {
	awc.logger.Println(append([]interface{}{"AWS-Go-SDK:"}, args...)...)
}

func init() {
	var err error
	binaryDir, err = osext.ExecutableFolder()
	if err != nil {
		errJson(nil, fmt.Errorf("Failed to get binary folder: %s", err))
	}

	logger = getLogger()
	logger.Println(os.Args)
}

func main() {

	test := flag.Bool("test", false, "Test Amazon S3 connection")
	uploadAllBackups := flag.Bool("upload-all-backups", false, "Upload all local backups. Example: -upload-all-backups")
	backupNameToUpload := flag.String("upload-backup", "", "Exports and uploads backup. Example: -upload-backup backup_info_1606061802.xml")
	uploadFile := flag.String("upload-file", "", "File upload. Example: -upload-file /tmt/backup.tar")

	resumeUploadBackup := flag.String("resume-upload-backup", "", "Resume backup upload. Example: -resume-upload-backup backup_info_1606061802.xml")
	resumeUploadFile := flag.String("resume-upload-file", "", "Resume file upload. Example: -resume-upload-file /tmp/backup.tar")
	resumeUploadId := flag.String("resume-upload-id", "", "Use with -resume-upload-file. Example: -resume-upload-id <...upload ID here...>")
	cancelUpload := flag.String("cancel-upload", "", "Cancel unfinished upload backup. Example: -cancel-upload backup_info_1606061802.xml")

	cancelDownload := flag.String("cancel-download", "", "Cancel unfinished download backup. Example: -cancel-download backup_info_1606061802.xml")
	downloadFile := flag.String("download-file", "", "Download file. Example: -download-file path/file.tar")
	resumeDownloadBackup := flag.String("resume-download-backup", "", "Resume download and import backup. Example: -resume-download-backup path/file.tar")
	resumeDownloadFile := flag.String("resume-download-file", "", "Resume download file. Example: -resume-download-file path/file.tar")
	downloadFileDst := flag.String("file-destination", "", "Destination of download file. Example: -file-destination /tmp/file.tar")
	downloadBackup := flag.String("download-backup", "", "Download and import backup. Example: -download-backup subscription.com/backup_info_1607161246_1607161349.xml.tar.gz")
	listStorage := flag.Bool("list", false, "List backups. Example: -list")
	listLocalStorage := flag.Bool("list-local-storage", false, "List local backups. Example: -list-local-storage")
	listRemoteStorage := flag.String("list-remote-storage", "", "List backups on storage. Example: -list-remote-storage amazon")
	listUnfinishedUploads := flag.Bool("list-unfinished-uploads", false, "List unfinished uploads. Example: -list-unfinished-uploads")
	remoteStorageSubFolder := flag.String("sub-folder", "", "List backups on storage in specific folder. Example: -sub-folder bla-bla")
	remoteFileInfo := flag.String("remote-file-info", "", "Show info about file on storage. Example: -remote-file-info path/file.tar")
	backupPassword := flag.String("backup-password", "", "Example: -backup-password s3cReT")
	skipIncremental := flag.Bool("skip-incremental", false, "Skip incremental backups. Example: -skip-incremental")
	useGzip := flag.Bool("use-gzip", false, "Gzip file before uploading. Example: -use-gzip")
	checkSign := flag.Bool("check-sign", false, "Check backup signature. Example: -check-sign")
	deleteAfterUpload := flag.Bool("delete-after-upload", false, "Delete source backup after upload. Example: -delete-after-upload")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")
	flag.Parse()

	err := selfCheck()
	if err != nil {
		errJson(logger, err)
	}

	p, err := plesk.New(logger)
	if err != nil {
		errJson(logger, err)
	}

	if dumpTmpD, ok := p.Config["DUMP_TMP_D"]; !ok || dumpTmpD == "" {
		errJson(logger, errors.New("DUMP_TMP_D is undefined"))
	}

	AmazonBucket := os.Getenv("AWS_BUCKET")

	os.Setenv("AWS_REGION", "eu-central-1") // set some random region
	// Create dummy AWS instance to determine bucket region
	dummyAwss3 := amazons3.AmazonS3{
		Log:    logger,
		Svc:    s3.New(session.New()),
		Region: "",
		Bucket: AmazonBucket,
	}
	// Set bucket region
	AmazonRegion, err := dummyAwss3.GetBucketLocation(AmazonBucket)
	if err != nil {
		logger.Printf("Failed to get bucket %s region: %s\n", AmazonBucket, err)
		errJson(logger, err)
		return
	}
	os.Setenv("AWS_REGION", AmazonRegion)

	awsConfig := &aws.Config{
		Region: aws.String(AmazonRegion),
		Logger: AwsCustomLogger{logger: logger},
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).Dial,
				IdleConnTimeout:       1 * time.Minute,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   128,
				TLSClientConfig: &tls.Config{
					ClientSessionCache: tls.NewLRUClientSessionCache(128),
				},
			},
		},
	}

	awss3 := amazons3.AmazonS3{
		Log:    logger,
		Svc:    s3.New(session.New(awsConfig.WithLogLevel(aws.LogDebug))),
		Region: AmazonRegion,
		Bucket: AmazonBucket,
	}

	backupFileExt := ".tar"
	if runtime.GOOS == "windows" {
		backupFileExt = ".zip"
	}

	b2a := BackupToAmazon{
		Log:               logger,
		Plesk:             p,
		AwsS3:             awss3,
		Backups:           map[string]Backup{},
		UseGzip:           *useGzip,
		HttpStatus:        make(map[string]string),
		HttpStatusSetter:  false,
		HttpHost:          "localhost",
		HttpPort:          "64837",
		UploadDir:         filepath.Join(p.Config["DUMP_TMP_D"], "upload"),
		DownloadDir:       filepath.Join(p.Config["DUMP_TMP_D"], "download"),
		BackupFileExt:     backupFileExt,
		SkipIncremental:   *skipIncremental,
		DeleteAfterUpload: *deleteAfterUpload,
	}

	err = os.MkdirAll(b2a.UploadDir, 0640)
	if err != nil {
		errJson(b2a.Log, fmt.Errorf("Failed to create destination path %s: %s\n", b2a.UploadDir, err))
		return
	}

	err = os.MkdirAll(b2a.DownloadDir, 0640)
	if err != nil {
		errJson(b2a.Log, fmt.Errorf("Failed to create destination path %s: %s\n", b2a.DownloadDir, err))
		return
	}

	go b2a.StartHttp()
	time.Sleep(1 * time.Second) // Wait for server start

	if *remoteStorageSubFolder != "" {
		b2a.SubFolder = *remoteStorageSubFolder
		if !strings.HasSuffix(string(*remoteStorageSubFolder), "/") {
			b2a.SubFolder = *remoteStorageSubFolder + "/"
		}
		b2a.Log.Printf("Sub-Folder is %s\n", b2a.SubFolder)
	}
	if *test {
		err := testAmazonSettings(b2a.SubFolder)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
		return
	}
	if *uploadAllBackups {
		err := b2a.uploadAllBackups()
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
		return
	}
	if *listStorage {
		_, err := b2a.listStorage()
		if err != nil {
			errJson(b2a.Log, err)
		}
		b2a.PrintlnJson(b2a.Backups)
		return
	}
	if *listLocalStorage {
		serverBackups, err := b2a.getServerBackupList()
		if err != nil {
			errJson(b2a.Log, err)
		}
		b2a.PrintlnJson(serverBackups)
		return
	}
	if *listRemoteStorage == "amazon" {
		list, err := b2a.AwsS3.GetBucketFilesList(b2a.SubFolder)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
		b2a.PrintlnJson(list)
		return
	}
	if *listUnfinishedUploads {
		list, err := b2a.listUnfinishedUploads()
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
		b2a.PrintlnJson(list)
		return
	}
	if *backupNameToUpload != "" {

		err = b2a.uploadBackup(*backupNameToUpload)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *uploadFile != "" {
		err = b2a.uploadFile(*uploadFile, "") // "" means root
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *resumeUploadBackup != "" {

		err = b2a.resumeUploadBackup(*resumeUploadBackup)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *resumeUploadFile != "" && *resumeUploadId != "" {
		err = b2a.resumeUploadFile(*resumeUploadFile, *resumeUploadId)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *cancelDownload != "" {
		if err = b2a.cancelDownload(*cancelDownload); err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *cancelUpload != "" {
		if err = b2a.cancelUpload(*cancelUpload); err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *remoteFileInfo != "" {
		resp, err := b2a.AwsS3.GetFileInfo(*remoteFileInfo)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
		b2a.PrintlnJson(resp)
		return
	}
	if *downloadBackup != "" {

		err = b2a.downloadBackup(*downloadBackup, *backupPassword, *checkSign)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *resumeDownloadBackup != "" {

		err = b2a.resumeDownloadBackup(*resumeDownloadBackup, *backupPassword, *checkSign)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *downloadFile != "" && *downloadFileDst != "" {
		err = b2a.downloadFile(*downloadFile, *downloadFileDst)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}
	if *resumeDownloadFile != "" && *downloadFileDst != "" {
		err = b2a.resumeDownloadFile(*resumeDownloadFile, *downloadFileDst)
		if err != nil {
			errJson(b2a.Log, err)
			return
		}
	}

	// keep at the end
	if *memprofile != "" {
		/*
			f, err := os.Create(*memprofile)
			if err != nil {
				errJson(b2a.Log, err)
				return
			}
			pprof.WriteHeapProfile(f)
			f.Close()
		*/
	}
}

func (self *BackupToAmazon) uploadAllBackups() (err error) {
	self.Log.Println("Start upload all backups")

	serverBackups, err := self.listStorage()
	if err != nil {
		return fmt.Errorf("Failed to upload all backups: %s", err)
	}

	for _, backup := range serverBackups {
		if backup.IsLocal && !backup.IsLocalInvalid && !backup.IsRemote && !backup.IsPartialDownload && !backup.IsPartialUpload {
			self.setHttpStatus(backup.Name, "upload")

			self.Log.Printf("Upload: %#v\n", backup)
			err = self.uploadBackup(backup.Name)
			if err != nil {
				return fmt.Errorf("Failed to upload all backups: %s", err)
			}
			if self.DeleteAfterUpload && backup.Backup.IncrementBaseFullname != "" { // Delete only incremental backups
				err = self.Plesk.DeleteBackupFromLocalStorage(backup.Name)
				if err != nil {
					errMessage := fmt.Sprintf("Failed to upload all backups. Failed to delete source backup %s after success upload: %s", backup.Name, err)
					if appErr, ok := err.(ErrCoder); ok {
						return Err{
							IsError:   true,
							Code:      appErr.Code(),
							LocaleKey: appErr.Code(),
							LocaleArgs: map[string]string{
								"backupName": backup.Name,
								"error":      err.Error(),
							},
							Message:   errMessage,
							OriginErr: err,
						}
					}

					return errors.New(errMessage)
				}
			}

			self.unSetHttpStatus(backup.Name, "upload")
		} else {
			self.Log.Printf("Skip: %#v\n", backup)
		}
	}

	self.Log.Println("All backups are successfully uploaded")
	return
}

func (self *BackupToAmazon) listStorage() (map[string]Backup, error) {
	serverBackups, err := self.getServerBackupList()
	if err != nil {
		return nil, err
	}
	remoteFiles, err := self.AwsS3.GetBucketFilesList(self.SubFolder)
	if err != nil {
		return nil, err
	}
	uploads, err := self.listUnfinishedUploads()
	if err != nil {
		return nil, err
	}
	fsFiles, err := self.listFsStorage()
	if err != nil {
		return nil, err
	}
	currentHttpStatus := self.getHttpStatus()

	for _, backup := range serverBackups {
		b := Backup{
			Name:     backup.Fullname,
			IsRemote: false,
			IsLocal:  true,
			Backup:   backup,
		}

		if backup.DumpStatus.DumpStatus == "PARTIAL" || backup.DumpStatus.DumpStatus == "WRONG-FORMAT" || backup.DumpStatus.BackupProcessStatus == "ERROR" {
			b.IsLocalInvalid = true
			self.Log.Printf("Backup %s is invalid. DumpStatus = %s, BackupProcessStatus = %s\n", backup.Fullname, backup.DumpStatus.DumpStatus, backup.DumpStatus.BackupProcessStatus)
		}

		for _, remoteFile := range remoteFiles {
			if *remoteFile.Key == self.SubFolder+backup.Fullname+self.BackupFileExt {
				b.IsRemote = true
				b.RemotePath = *remoteFile.Key
			}
		}
		for _, fsFile := range fsFiles[self.DownloadDir] {
			self.Log.Printf("%s starts with %s\n", fsFile.Name(), backup.Name)
			if strings.HasPrefix(fsFile.Name(), backup.Name) {
				b.IsPartialDownload = true
			}
		}
		for _, fsFile := range fsFiles[self.UploadDir] {
			self.Log.Printf("%s starts with %s\n", fsFile.Name(), backup.Name)
			if strings.HasPrefix(fsFile.Name(), backup.Name) {
				b.IsPartialUpload = true
			}
		}
		self.Backups[backup.Fullname] = b

	}

	for _, upload := range uploads {
		remoteBackupName := strings.TrimPrefix(strings.TrimSuffix(upload.RemotePath, filepath.Ext(upload.RemotePath)), self.SubFolder)

		if backup, ok := self.Backups[remoteBackupName]; ok {
			backup.IsPartialUpload = true
			backup.UploadId = upload.UploadId
			self.Backups[remoteBackupName] = backup
		} else {
			self.Backups[remoteBackupName] = Backup{
				Name:            remoteBackupName,
				UploadId:        upload.UploadId,
				IsLocal:         false,
				IsRemote:        true,
				IsPartialUpload: true,
				RemotePath:      upload.RemotePath,
				Backup: plesk.Dump{
					Size:         upload.Size,
					CreationDate: upload.LastModified.In(time.Local).String(),
				},
			}
		}
	}

	for _, remoteFile := range remoteFiles {
		if strings.HasSuffix(*remoteFile.Key, "/") { // General folder
			continue // Ignore it
		}

		if strings.Contains(strings.TrimPrefix(*remoteFile.Key, self.SubFolder), "/") { // Content of underlying sub-folders
			continue // Ignore it
		}

		if !strings.HasSuffix(*remoteFile.Key, self.BackupFileExt) { // Not supported file type
			continue // Ignore it
		}

		remoteBackupName := strings.TrimPrefix(*remoteFile.Key, self.SubFolder)
		remoteBackupName = strings.TrimSuffix(remoteBackupName, self.BackupFileExt)

		if remoteBackupName == "" { // Sub-Folder itself
			continue
		}

		if _, ok := self.Backups[remoteBackupName]; !ok {

			isPartialDownload := false
			for _, fsFile := range fsFiles[self.DownloadDir] {
				self.Log.Printf("%s starts with %s\n", fsFile.Name(), remoteBackupName)
				if strings.HasPrefix(fsFile.Name(), remoteBackupName) {
					isPartialDownload = true
				}
			}

			self.Backups[remoteBackupName] = Backup{
				Name:              remoteBackupName,
				IsLocal:           false,
				IsRemote:          true,
				IsPartialDownload: isPartialDownload,
				RemotePath:        *remoteFile.Key,
				Backup: plesk.Dump{
					Size:         *remoteFile.Size,
					CreationDate: remoteFile.LastModified.In(time.Local).String(),
				},
			}
		}
	}

	for backupName, operation := range currentHttpStatus {
		if backup, ok := self.Backups[backupName]; ok {
			backup.InProgress = true
			self.Backups[backupName] = backup
			self.Log.Printf("Backup in progress %s: %s\n", operation, backupName)
		}
	}

	return self.Backups, nil
}

func (self *BackupToAmazon) getServerBackupList() (list []plesk.Dump, err error) {
	allBackups, err := self.Plesk.GetBackupListFromLocalStorage()
	if err != nil {
		return nil, err
	}
	for _, backup := range allBackups {
		if self.SkipIncremental && backup.IncrementBaseFullname != "" { // Skip incremental backups
			continue
		}
		if backup.DumpObject.Type != "server" { // Skip non-server backups
			continue
		}

		list = append(list, backup)
	}

	return
}

func (self *BackupToAmazon) listFsStorage() (files map[string][]os.FileInfo, err error) {
	files = map[string][]os.FileInfo{}

	unfinishedUploadFiles, err := ioutil.ReadDir(self.UploadDir)
	if err != nil {
		return nil, fmt.Errorf("Failed list directory %s: %s", self.UploadDir, err)
	}
	files[self.UploadDir] = unfinishedUploadFiles

	unfinishedDownloadFiles, err := ioutil.ReadDir(self.DownloadDir)
	if err != nil {
		return nil, fmt.Errorf("Failed list directory %s: %s", self.DownloadDir, err)
	}
	files[self.DownloadDir] = unfinishedDownloadFiles

	self.Log.Printf("List File Systems Storage: %s", files)
	return
}

func (self *BackupToAmazon) listUnfinishedUploads() (uploads []UnfinishedUpload, err error) {
	unfUps, err := self.AwsS3.ListUnfinishedUploads()
	if err != nil {
		return
	}

	for _, upload := range unfUps {
		var size int64
		var lastModified time.Time
		resp, err := self.AwsS3.ListParts(*upload.Key, *upload.UploadId)
		if err != nil {
			return nil, err
		}
		for _, part := range resp.Parts {
			size = size + *part.Size
			lastModified = *part.LastModified
		}

		uploads = append(uploads, UnfinishedUpload{
			RemotePath:   *upload.Key,
			UploadId:     *upload.UploadId,
			Size:         size,
			LastModified: lastModified,
		})
	}
	return
}
func (self BackupToAmazon) resumeDownloadFile(amazonPath, destinationPath string) (err error) {
	err = self.AwsS3.ResumeDownload(amazonPath, destinationPath)
	if err != nil {
		return fmt.Errorf("Failed to resume download file from %s to %s with error: %s\n", amazonPath, destinationPath, err)
	}
	self.Log.Printf("Successfully reume download file %s to %s\n", amazonPath, destinationPath)
	return nil
}

func (self *BackupToAmazon) cancelDownload(backupName string) error {
	tmpDownloadFile := filepath.Join(self.DownloadDir, backupName+self.BackupFileExt)
	self.Log.Printf("Cancel download %s. Delete temprorary file: %s\n", backupName, tmpDownloadFile)
	return os.Remove(tmpDownloadFile)
}

func (self BackupToAmazon) downloadFile(amazonPath, destinationPath string) (err error) {
	err = self.AwsS3.Download(amazonPath, destinationPath)
	if err != nil {
		return fmt.Errorf("Failed to download file from %s to %s with error: %s\n", amazonPath, destinationPath, err)
	}
	self.Log.Printf("Successfully download file %s to %s\n", amazonPath, destinationPath)
	return nil
}

func (self BackupToAmazon) remotePathToBackupName(remotePath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(remotePath, self.SubFolder), self.BackupFileExt)
}

func (self BackupToAmazon) pathToBackupName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func (self BackupToAmazon) downloadBackup(amazonFullPath, backupPassword string, checkSign bool) error {
	defer self.unSetHttpStatus(self.remotePathToBackupName(amazonFullPath), "download")
	self.setHttpStatus(self.remotePathToBackupName(amazonFullPath), "download")

	dstFilePath := filepath.Join(self.DownloadDir, strings.TrimPrefix(amazonFullPath, self.SubFolder))

	err := self.AwsS3.Download(amazonFullPath, dstFilePath)
	if err != nil {
		return fmt.Errorf("Failed to download backup from %s to %s with error: %s\n", amazonFullPath, dstFilePath, err)
	}

	err = self.Plesk.ImportBackupToLocalStorage(dstFilePath, backupPassword, checkSign)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to download backup: %s", err)
		if appErr, ok := err.(ErrCoder); ok {
			return Err{
				IsError:   true,
				Code:      appErr.Code(),
				LocaleKey: appErr.Code(),
				LocaleArgs: map[string]string{
					"backupName":  self.pathToBackupName(dstFilePath),
					"dstFilePath": dstFilePath,
					"checkSign":   fmt.Sprintf("%t", checkSign),
					"error":       err.Error(),
				},
				Message:   errMessage,
				OriginErr: err,
			}
		}

		return errors.New(errMessage)
	}

	err = os.Remove(dstFilePath)
	if err != nil {
		return fmt.Errorf("Failed to remove %s with error: %s\n", dstFilePath, err)
	}

	self.Log.Printf("Successfully download and import backup file: %s\n", amazonFullPath)
	return nil
}

func (self BackupToAmazon) resumeDownloadBackup(amazonFullPath, backupPassword string, checkSign bool) error {
	defer self.unSetHttpStatus(self.remotePathToBackupName(amazonFullPath), "download")
	self.setHttpStatus(self.remotePathToBackupName(amazonFullPath), "download")

	dstFilePath := filepath.Join(self.DownloadDir, strings.TrimPrefix(amazonFullPath, self.SubFolder))

	err := self.AwsS3.ResumeDownload(amazonFullPath, dstFilePath)
	if err != nil {
		return fmt.Errorf("Failed to resume download backup from %s to %s with error: %s\n", amazonFullPath, dstFilePath, err)
	}

	err = self.Plesk.ImportBackupToLocalStorage(dstFilePath, backupPassword, checkSign)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to resume download backup %s from %s with error: %s\n", amazonFullPath, dstFilePath, err)
		if appErr, ok := err.(ErrCoder); ok {
			return Err{
				IsError:   true,
				Code:      appErr.Code(),
				LocaleKey: appErr.Code(),
				LocaleArgs: map[string]string{
					"backupName":  self.pathToBackupName(dstFilePath),
					"dstFilePath": dstFilePath,
					"checkSign":   fmt.Sprintf("%t", checkSign),
					"error":       err.Error(),
				},
				Message:   errMessage,
				OriginErr: err,
			}
		}

		return errors.New(errMessage)
	}

	err = os.Remove(dstFilePath)
	if err != nil {
		return fmt.Errorf("Failed to remove %s with error: %s\n", dstFilePath, err)
	}

	self.Log.Printf("Successfully download and import backup file: %s\n", amazonFullPath)
	return nil
}

func (self BackupToAmazon) uploadBackup(backupFullName string) error {
	defer self.unSetHttpStatus(backupFullName, "upload")
	self.setHttpStatus(backupFullName, "upload")

	backupName := filepath.Base(backupFullName)
	amazonDstPath := backupName

	backupFound := false
	backupIncremental := true
	var backupSize int64

	list, err := self.Plesk.GetBackupListFromLocalStorage()
	if err != nil {
		return fmt.Errorf("Failed to upload backup: %s\n", err)
	}
	for _, backup := range list {
		if backup.Fullname != backupFullName {
			continue
		}

		backupFound = true
		backupSize = backup.Size
		if backup.IncrementBaseFullname != "" {
			backupIncremental = true
		} else {
			backupIncremental = false
		}

		if backup.DumpObject.Type == "server" {
			amazonDstPath = self.SubFolder
		} else {
			amazonDstPath = self.SubFolder + backup.DumpObject.Name + "/"
		}

		break
	}

	if !backupFound {
		return fmt.Errorf("Failed to found backup in local storage: %s\n", backupFullName)
	}

	freeSpace, err := getFreeDiskSpaceInPath(filepath.Clean(self.UploadDir))
	if err != nil {
		return fmt.Errorf("Failed to upload backup: Failed to get free disk space: %s\n", err)
	}

	if freeSpace <= uint64(backupSize) {
		err = errors.New(fmt.Sprintf("Failed to upload backup: Free disk space %d bytes in upload dir %s is not enough to export backup %s of size %d bytes\n", freeSpace, self.UploadDir, backupFullName, backupSize))
		return Err{
			IsError:   true,
			Message:   err.Error(),
			OriginErr: err,
			LocaleKey: "BackupUploadFailedNoDiskSpaceForExport",
			LocaleArgs: map[string]string{
				"uploadDir":  self.UploadDir,
				"backupName": backupFullName,
				"backupSize": string(backupSize),
				"freeSpace":  string(freeSpace),
			},
		}
	}

	self.Log.Printf("Free disk space %d is enough for backup %s size %d\n", freeSpace, backupFullName, backupSize)

	exportedBackupPath := filepath.Join(self.UploadDir, backupName) + self.BackupFileExt
	err = self.Plesk.ExportBackupFromLocalStorage(backupFullName, exportedBackupPath, false)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to upload backup: %s\n", err)
		if appErr, ok := err.(ErrCoder); ok {
			return Err{
				IsError:   true,
				Code:      appErr.Code(),
				LocaleKey: appErr.Code(),
				LocaleArgs: map[string]string{
					"backupFullName":     backupFullName,
					"exportedBackupPath": exportedBackupPath,
					"includeIcrements":   "false",
					"error":              err.Error(),
				},
				Message:   errMessage,
				OriginErr: err,
			}
		}

		return errors.New(errMessage)
	}

	err = self.uploadFile(exportedBackupPath, amazonDstPath)
	if err != nil {
		return err // Do not format this err because it can be amazons3.UploadErr
	}

	err = os.Remove(exportedBackupPath)
	if err != nil {
		return fmt.Errorf("Failed to remove exported backup %s after success upload: %s\n", exportedBackupPath, err)
	}

	if self.DeleteAfterUpload && backupIncremental { // Delete only incremental backups
		err = self.Plesk.DeleteBackupFromLocalStorage(backupName)
		if err != nil {
			errMessage := fmt.Sprintf("Failed to remove source backup %s after success upload: %s\n", backupName, err)
			if appErr, ok := err.(ErrCoder); ok {
				return Err{
					IsError:   true,
					Code:      appErr.Code(),
					LocaleKey: appErr.Code(),
					LocaleArgs: map[string]string{
						"backupName": backupName,
						"error":      err.Error(),
					},
					Message:   errMessage,
					OriginErr: err,
				}
			}

			return errors.New(errMessage)
		}
	}

	self.Log.Printf("Successfully upload backup: %s\n", backupFullName)
	return nil
}

func (self BackupToAmazon) uploadFile(uploadFilePath, dstPath string) (err error) {
	err = self.AwsS3.Upload(uploadFilePath, dstPath, self.UseGzip)
	if err != nil {
		if uploadErr, ok := err.(s3manager.MultiUploadFailure); ok {
			return uploadErr
		}
		return fmt.Errorf("Failed to upload file %s to %s: %s\n", uploadFilePath, dstPath, err)
	}

	return
}

func (self BackupToAmazon) resumeUploadBackup(backupName string) (err error) {
	defer self.unSetHttpStatus(backupName, "upload")
	self.setHttpStatus(backupName, "upload")

	fsFiles, err := self.listFsStorage()
	if err != nil {
		return fmt.Errorf("Failed to resume upload backup. Failed to list fs storage: %s\n", err)
	}
	var uploadFile string
	for _, upload := range fsFiles[self.UploadDir] {
		self.Log.Printf("Compare FS storage files: %s == %s\n", upload.Name(), backupName)
		if strings.HasPrefix(upload.Name(), backupName) {
			uploadFile = filepath.Join(self.UploadDir, upload.Name())
			self.Log.Printf("Found: %s starts with %s\n", upload.Name(), backupName)
			break
		}
	}
	if uploadFile == "" {
		return fmt.Errorf("Failed to resume upload backup. Failed to find upload file: %s in list %s\n", backupName+self.BackupFileExt, fsFiles[self.UploadDir])
	}

	uploads, err := self.AwsS3.ListUnfinishedUploads()
	if err != nil {
		return fmt.Errorf("Failed to resume upload backup. Failed to list unfinished uploads: %s\n", err)
	}
	var uploadId string
	for _, upload := range uploads {
		uploadKey := self.SubFolder + backupName + self.BackupFileExt
		self.Log.Printf("Compare remote key: %s == %s\n", *upload.Key, uploadKey)
		if self.SubFolder+backupName+self.BackupFileExt == *upload.Key {
			uploadId = *upload.UploadId
			break
		}

	}

	if uploadId == "" {
		self.Log.Printf("Unfinished upload not found on remote storage. Start new upload for %s", backupName)
		return self.uploadBackup(backupName)
	} else {
		err = self.resumeUploadFile(uploadFile, uploadId)
	}

	if err != nil {
		return fmt.Errorf("Failed to resume upload backup: %s\n", err)
	}

	err = os.Remove(uploadFile)
	if err != nil {
		return fmt.Errorf("Failed to remove temporary upload file %s after success upload: %s\n", uploadFile, err)
	}

	if self.DeleteAfterUpload {

		list, err := self.Plesk.GetBackupListFromLocalStorage()
		if err != nil {
			return fmt.Errorf("Failed to resume upload backup: %s\n", err)
		}

		backupFound := false
		backupIncremental := true
		for _, backup := range list {
			if backup.Fullname != backupName {
				continue
			}

			backupFound = true
			if backup.IncrementBaseFullname != "" {
				backupIncremental = true
			} else {
				backupIncremental = false
			}

			break
		}

		if backupFound && backupIncremental {
			err = self.Plesk.DeleteBackupFromLocalStorage(backupName)
			if err != nil {
				errMessage := fmt.Sprintf("Failed to remove source backup %s after success upload: %s\n", backupName, err)
				if appErr, ok := err.(ErrCoder); ok {
					return Err{
						IsError:   true,
						Code:      appErr.Code(),
						LocaleKey: appErr.Code(),
						LocaleArgs: map[string]string{
							"backupName": backupName,
							"error":      err.Error(),
						},
						Message:   errMessage,
						OriginErr: err,
					}
				}

				return errors.New(errMessage)
			}
		}
	}

	self.Log.Printf("Successfully resume upload backup: %s\n", backupName)
	return
}

func (self BackupToAmazon) resumeUploadFile(uploadFile, uploadId string) (err error) {
	uploads, err := self.AwsS3.ListUnfinishedUploads()
	if err != nil {
		return fmt.Errorf("Failed to resume upload file. Failed to list unfinished uploads: %s\n", err)
	}
	var key string
	for _, upload := range uploads {
		self.Log.Printf("%s == %s\n", *upload.UploadId, uploadId)
		if string(*upload.UploadId) == uploadId {
			key = *upload.Key
			break
		}
	}
	if key == "" {
		return fmt.Errorf("Failed to resume upload file. Failed to find upload id: %s in list %s\n", uploadId, uploads)
	}

	err = self.AwsS3.ResumeUpload(uploadFile, key, uploadId, self.UseGzip)
	if err != nil {
		return fmt.Errorf("Failed to resume upload file: %s\n", err)
	}

	return
}

func (self BackupToAmazon) cancelUpload(backupName string) (err error) {
	tmpUploadFile := filepath.Join(self.UploadDir, backupName+self.BackupFileExt)
	if _, err := os.Stat(tmpUploadFile); !os.IsNotExist(err) {
		self.Log.Printf("Cancel upload. Delete temporary file %s\n", tmpUploadFile)
		err = os.Remove(tmpUploadFile)
		if err != nil {
			return fmt.Errorf("Failed to cancel upload. Failed to delete %s: %s\n", tmpUploadFile, err)
		}
	}

	uploads, err := self.AwsS3.ListUnfinishedUploads()
	if err != nil {
		return fmt.Errorf("Failed to cancel upload. Failed to list unfinished uploads: %s\n", err)
	}
	var key string
	var uploadId string
	for _, upload := range uploads {
		self.Log.Printf("%s == %s\n", *upload.Key, backupName)
		if strings.Contains(*upload.Key, self.SubFolder+backupName) {
			key = *upload.Key
			uploadId = *upload.UploadId
			break
		}
	}
	if key == "" {
		self.Log.Println("Cancel upload. Unfinished upload not foum on remote storage")
		return
	}

	err = self.AwsS3.AbortUpload(key, uploadId)
	if err != nil {
		return fmt.Errorf("Failed to cancel upload for %s: %s\n", backupName, err)
	}

	return
}

func selfCheck() error {
	AmazonBucket := os.Getenv("AWS_BUCKET")
	AmazonKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	AmazonSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if AmazonBucket == "" {
		return errors.New("Environment variable AWS_BUCKET is not set")
	}
	if AmazonKeyId == "" {
		return errors.New("Environment variable AWS_ACCESS_KEY_ID is not set")
	}
	if AmazonSecretKey == "" {
		return errors.New("Environment variable AWS_SECRET_ACCESS_KEY is not set")
	}

	ntpServer := "0.pool.ntp.org"
	ntpTime, err := ntp.Query(ntpServer, 4)
	if err != nil {
		return nil // Do not block application if firewall enabled
	}
	if ntpTime.ClockOffset.Hours() > 1.0 { // Avoid Amazon RequestTimeTooSkewed error
		msg := fmt.Sprintf(
			"Gap in local and real time detected, local time is %s, real time is %s, gap in time is %s\n",
			time.Now().String(),
			ntpTime.Time.String(),
			ntpTime.ClockOffset.String(),
		)

		err := Err{
			IsError:   true,
			Message:   msg,
			LocaleKey: "gapInTime",
			LocaleArgs: map[string]string{
				"local": time.Now().String(),
				"real":  ntpTime.Time.String(),
				"diff":  ntpTime.ClockOffset.String(),
			},
		}

		return err
	}

	return nil
}

func getLogger() *log.Logger {
	logFilePath := path.Join(binaryDir, backupOnAmazon+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal("Failed to open log file", logFilePath, err)
	}

	return log.New(logFile, backupOnAmazon+" ", log.LstdFlags)
}

func createTestFile() string {
	testFilePath := path.Join(binaryDir, "test_"+getRandomString(10)+".txt")
	testFile, err := os.OpenFile(testFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal("Failed to create test file", testFilePath, err)
	}
	defer func() { _ = testFile.Close() }()

	_, err = testFile.Write([]byte("test amazon s3"))
	if err != nil {
		return ""
	}

	return testFilePath
}

func testAmazonSettings(subFolder string) (err error) {
	AmazonRegion := os.Getenv("AWS_REGION")
	AmazonBucket := os.Getenv("AWS_BUCKET")

	awss3 := amazons3.AmazonS3{
		Log:    getLogger(),
		Svc:    s3.New(session.New(&aws.Config{Region: aws.String(AmazonRegion)})),
		Region: AmazonRegion,
		Bucket: AmazonBucket,
	}

	err = awss3.IsRegionValid(AmazonRegion)
	if err != nil {
		return
	}

	bucketList, err := awss3.GetBucketFilesList("")
	if err != nil {
		return
	}
	if subFolder != "" {

		if !strings.HasSuffix(subFolder, "/") {
			subFolder = subFolder + "/"
		}
		subFolder = strings.Replace(subFolder, "//", "/", -1)

		subFolderExists := false

		for _, file := range bucketList {
			if *file.Key == subFolder || *file.Key == subFolder+"/" {
				subFolderExists = true
			}
		}
		if !subFolderExists {
			err = awss3.CreateFolder(subFolder)
			if err != nil {
				return
			}
		}
	}
	testFile := createTestFile()
	testFileDownloaded := testFile + ".downloaded"
	testFileDownloadedSuccess := false
	defer func() {
		err = os.Remove(testFile)
		if err != nil {
			return
		}
		if testFileDownloadedSuccess {
			err = os.Remove(testFileDownloaded)
			if err != nil {
				return
			}
		}
	}()

	err = awss3.Upload(testFile, subFolder, false)
	if err != nil {
		return
	}

	err = awss3.Download(subFolder+filepath.Base(testFile), testFileDownloaded)
	if err != nil {
		return
	}
	testFileDownloadedSuccess = true

	err = awss3.Delete(subFolder + filepath.Base(testFile))
	if err != nil {
		return
	}

	return
}

func execute(log *log.Logger, command string, args ...string) (output string, outputBytes []byte, code int, err error) {
	//log.Printf("%s %s", command, args)

	cmd := exec.Command(command, args...)
	var waitStatus syscall.WaitStatus

	if outputBytes, err = cmd.CombinedOutput(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			code = waitStatus.ExitStatus()
		}
	} else {
		// Command was successful
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		code = waitStatus.ExitStatus()
	}

	output = string(outputBytes)
	if err != nil {
		log.Println("output: ", output, "err: ", err, "code: ", code)
	}

	return
}

func (self BackupToAmazon) unpackGzFile(gzFilePath, dstFilePath string) (int64, error) {
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s for unpack: %s", gzFilePath, err)
	}
	defer self.IoClose(gzFile)

	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return 0, fmt.Errorf("Failed to create destination file %s for unpack: %s", dstFilePath, err)
	}
	defer self.IoClose(dstFile)

	ioReader, ioWriter := io.Pipe()

	go func() { // goroutine leak is possible here
		gzReader, _ := gzip.NewReader(gzFile)
		// it is important to close the writer or reading from the other end of the
		// pipe or io.copy() will never finish
		defer func() {
			self.IoClose(gzFile)
			self.IoClose(gzReader)
			self.IoClose(ioWriter)
		}()

		written, err := io.Copy(ioWriter, gzReader)
		if err != nil {
			self.Log.Printf("unpackGzFile io.Copy error: %s\n", err)
		}
		self.Log.Printf("unpackGzFile io.Copy bytes written: %s\n", string(written))
	}()

	written, err := io.Copy(dstFile, ioReader)
	if err != nil {
		return 0, err // goroutine leak is possible here
	}
	err = ioReader.Close()
	if err != nil {
		return 0, err // goroutine leak is possible here
	}

	return written, nil
}

func (self BackupToAmazon) IoClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		self.Log.Println(err)
	}
}

func (self BackupToAmazon) PrintlnJson(v interface{}) {
	jsonOutput, err := json.Marshal(v)
	if err != nil {
		self.Log.Fatal(err)
		return
	}
	self.Log.Println(string(jsonOutput))
	fmt.Println(string(jsonOutput))
	return
}

func errJson(logger *log.Logger, err error) {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}

	errStruct := Err{ // Default error
		IsError: true,
		Message: err.Error(),
	}

	if errStr, ok := err.(Err); ok {
		errStruct.IsError = errStr.IsError
		errStruct.Code = errStr.Code
		errStruct.Message = errStr.Message
		errStruct.UploadID = errStr.UploadID
		errStruct.LocaleKey = errStr.LocaleKey
		errStruct.LocaleArgs = errStr.LocaleArgs
	}

	if awsErr, ok := err.(awserr.Error); ok {
		errStruct.Code = awsErr.Code()
		errStruct.Message = awsErr.Message()
		errStruct.LocaleKey = awsErr.Code()
	}

	if uploadErr, ok := err.(s3manager.MultiUploadFailure); ok {
		errStruct.Code = uploadErr.Code()
		errStruct.Message = uploadErr.Message()
		errStruct.LocaleKey = uploadErr.Code()
		errStruct.UploadID = uploadErr.UploadID()
	}

	jsonOutput, err := json.Marshal(errStruct)
	if err != nil {
		jsonOutput, _ := json.Marshal(err)
		logger.Println(string(jsonOutput))
		println(string(jsonOutput))
		os.Exit(1)
		return
	}
	logger.Println(string(jsonOutput))
	println(string(jsonOutput))
	os.Exit(1)
	return
}

func getRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
