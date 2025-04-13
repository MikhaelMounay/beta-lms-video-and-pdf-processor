package uploader

import (
	"fmt"
	"io"
	"time"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/logger"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/utils"
)

type ProgressReader struct {
	Log           *logger.Logger
	Reader        io.Reader
	TotalSize     int64
	BytesRead     int64
	OnProgress    func(readBytes int64, totalBytes int64)
	lastTime      time.Time
	lastReadBytes int64
	startTime     time.Time
	fullyRead     bool
}

func NewProgressReader(log *logger.Logger, reader io.Reader, totalSize int64, onProgressCallback func(int64, int64)) *ProgressReader {
	return &ProgressReader{
		Log:           log,
		Reader:        reader,
		TotalSize:     totalSize,
		BytesRead:     0,
		OnProgress:    onProgressCallback,
		lastTime:      time.Now(),
		lastReadBytes: 0,
	}
}

func (pr *ProgressReader) PrintProgress(readBytes int64, totalBytes int64) {
	now := time.Now()
	elapsed := now.Sub(pr.lastTime).Seconds()
	bytesTransferred := readBytes - pr.lastReadBytes

	// Avoid division by zero
	if elapsed == 0 {
		elapsed = 0.001
	}

	speed := float64(bytesTransferred) / elapsed // bytes/sec
	remainingBytes := totalBytes - readBytes
	eta := float64(remainingBytes) / speed // seconds

	percent := float64(readBytes) / float64(totalBytes) * 100

	speedStr := utils.ByteCountSI(int64(speed)) + "/s"
	etaStr := time.Duration(eta * float64(time.Second)).Truncate(time.Second).String()

	if pr.BytesRead >= pr.TotalSize && !pr.fullyRead {
		pr.fullyRead = true
		fmt.Print("\r\033[K")
		pr.Log.Logf("Successfully Uploaded (Elapsed Time: %s)\n", now.Sub(pr.startTime).Truncate(time.Second).String())
	} else if !pr.fullyRead {
		fmt.Print("\r\033[K")
		fmt.Printf("Uploading... %.2f%% complete | %s | ETA: %s", percent, speedStr, etaStr)
	}
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)

	if pr.BytesRead == 0 {
		pr.startTime = time.Now()
	}

	pr.BytesRead += int64(n)
	if pr.OnProgress != nil {
		pr.OnProgress(pr.BytesRead, pr.TotalSize)
	} else {
		pr.PrintProgress(pr.BytesRead, pr.TotalSize)
	}

	return n, err
}
