package pdfprocessor

import (
	"fmt"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/shared"
)

type PdfFile struct {
	shared.File
}

func (f *PdfFile) ProcessFile() error {
	if err := shared.EncryptFile(f, config.SECRET_KEY, config.IV); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	f.UploadPaths = append(f.UploadPaths, f.GetNewPath()+".enc")

	return nil
}
