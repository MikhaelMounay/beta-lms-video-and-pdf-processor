package history

import (
	"context"
	"errors"
	"time"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UploadEntry struct {
	dbPool           *pgxpool.Pool
	Id               int
	InstanceName     string
	Grade            string
	OriginalFileName string
	FileId           string
	RemoteFileName   string
	FileHashSHA256   string
	UploadedAt       time.Time
}

func NewUploadEntry(grade string, originalFileName string, fileId string, remoteFileName string, fileHashSHA256 string) *UploadEntry {
	return &UploadEntry{
		dbPool:           config.DbPool,
		InstanceName:     config.INSTANCE_NAME,
		Grade:            grade,
		OriginalFileName: originalFileName,
		FileId:           fileId,
		RemoteFileName:   remoteFileName,
		FileHashSHA256:   fileHashSHA256,
	}
}

func (ue *UploadEntry) Save() error {
	if ue.dbPool == nil {
		return errors.New("no existing DB connection")
	}

	query := `
		INSERT INTO uploads (
            instance_name,
            grade,
            original_filename,
            file_id,
		    remote_filename,
            file_hash_sha256
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, uploaded_at;
	`

	err := ue.dbPool.QueryRow(context.Background(), query, ue.InstanceName, ue.Grade, ue.OriginalFileName, ue.FileId, ue.RemoteFileName, ue.FileHashSHA256).Scan(&ue.Id, &ue.UploadedAt)

	return err
}
