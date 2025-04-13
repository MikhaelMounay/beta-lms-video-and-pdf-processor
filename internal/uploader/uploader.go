package uploader

import (
	"context"
	"os"
	"path/filepath"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/history"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/logger"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/shared"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileUploader interface {
	AddToQueue(filePath string, destFilePath string, file shared.IFile)
	StartQueueUpload()
	GetQueueSize() int
}

type UploadRequest struct {
	FilePath     string
	DestFilePath string
	File         shared.IFile
}

type MinIOFileUploader struct {
	Log          *logger.Logger
	Queue        []UploadRequest
	Endpoint     string
	AppKeyId     string
	AppKey       string
	BucketId     string
	BucketName   string
	BucketRegion string
	BaseB2Folder string
	MinioClient  *minio.Client
}

func (u *MinIOFileUploader) AddToQueue(filePath string, destFilePath string, file shared.IFile) {
	u.Queue = append(u.Queue, UploadRequest{FilePath: filePath, DestFilePath: destFilePath, File: file})
}

func (u *MinIOFileUploader) StartQueueUpload() {
	if len(u.Queue) == 0 {
		u.Log.Log("No files to upload!\n")
		return
	}

	u.Log.Log("Starting Upload Queue...\n")

	ctx := context.Background()
	for i, v := range u.Queue {
		func() {
			// Prepare file and start uploading
			file, err := os.Open(v.FilePath)
			if err != nil {
				u.Log.Logf("Error opening file %s: %s\n", v.FilePath, err.Error())
				return
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				u.Log.Logf("Error getting file info for %s: %s\n", v.FilePath, err.Error())
				return
			}

			progressReader := NewProgressReader(u.Log, file, fileInfo.Size(), nil)

			u.Log.Logf("Queue (%d/%d): %s (%s)\n", i+1, len(u.Queue), filepath.Base(v.FilePath), utils.ByteCountSI(progressReader.TotalSize))

			_, err = u.MinioClient.PutObject(ctx, u.BucketName, v.DestFilePath, progressReader, progressReader.TotalSize, minio.PutObjectOptions{ContentType: "b2/x-auto"})
			if err != nil {
				u.Log.Logf("Error Uploading File: %v\n", err.Error())
				return
			}

			// Save the upload entry to DB
			fileHash, err := utils.ComputeSHA256Hash(v.FilePath)
			if err != nil {
				u.Log.Logf("Error computing file hash for %s: %s\n", v.FilePath, err.Error())
			}
			uploadEntry := history.NewUploadEntry(v.File.GetGrade(), v.File.GetName(), v.File.GetId(), filepath.Base(v.DestFilePath), fileHash)
			if err := uploadEntry.Save(); err != nil {
				u.Log.Logf("Error saving upload entry: %s\n", err.Error())
			}
		}()

		if err := os.Remove(v.FilePath); err != nil {
			u.Log.Logf("Error removing file %s: %s\n", v.FilePath, err.Error())
		}
	}
}

func (u *MinIOFileUploader) GetQueueSize() int {
	return len(u.Queue)
}

func ConfigUploader(log *logger.Logger) FileUploader {
	fUploader := &MinIOFileUploader{
		Log:          log,
		Endpoint:     config.ENDPOINT,
		AppKeyId:     config.APP_KEY_ID,
		AppKey:       config.APP_KEY,
		BucketId:     config.BUCKET_ID,
		BucketName:   config.BUCKET_NAME,
		BucketRegion: config.BUCKET_REGION,
		BaseB2Folder: config.BASE_B2_FOLDER,
	}

	minioClient, err := minio.New(fUploader.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(fUploader.AppKeyId, fUploader.AppKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal("Error: Couldn't connect to server : " + err.Error())
	}

	fUploader.MinioClient = minioClient

	return fUploader
}
