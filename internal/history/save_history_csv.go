package history

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
)

func SaveHistoryCSV(outputFile string) error {
	// Querying DB
	query := `
		SELECT grade, original_filename, file_id, remote_filename, file_hash_sha256, uploaded_at
		FROM uploads
		WHERE instance_name = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := config.DbPool.Query(context.Background(), query, config.INSTANCE_NAME)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	err = writer.Write([]string{
		"Index", "Grade", "Original Filename", "File ID", "Remote Filename", "File Hash (SHA256)", "Uploaded At",
	})
	if err != nil {
		return fmt.Errorf("write header error: %w", err)
	}

	i := 0
	for rows.Next() {
		i++

		var entry UploadEntry
		err := rows.Scan(
			&entry.Grade,
			&entry.OriginalFileName,
			&entry.FileId,
			&entry.RemoteFileName,
			&entry.FileHashSHA256,
			&entry.UploadedAt,
		)
		if err != nil {
			return fmt.Errorf("row scan error: %w", err)
		}

		err = writer.Write([]string{
			fmt.Sprintf("%d", i), // 1-based indexing for entries in the CSV file
			entry.Grade,
			entry.OriginalFileName,
			entry.FileId,
			entry.RemoteFileName,
			entry.FileHashSHA256,
			entry.UploadedAt.Format("2006-01-02 15:04:05 EET"),
		})
		if err != nil {
			return fmt.Errorf("write row error: %w", err)
		}
	}

	return nil
}
