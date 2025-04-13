package videoprocessor

import (
	"fmt"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/shared"
)

type VideoFile struct {
	shared.File
}

func (f *VideoFile) ProcessFile() error {
	if err := PackageVideoFile(f, config.VOD_KEY_ID, config.VOD_KEY, config.IV); err != nil {
		return fmt.Errorf("error packaging video file: %w", err)
	}

	if err := shared.EncryptFile(f, config.SECRET_KEY, config.IV); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	f.UploadPaths = append(f.UploadPaths, f.GetNewPath()+".enc")
	f.UploadPaths = append(f.UploadPaths, f.GetNewPath()+"_v.mp4")
	f.UploadPaths = append(f.UploadPaths, f.GetNewPath()+"_a.mp4")
	f.UploadPaths = append(f.UploadPaths, f.GetNewPath()+"_m.mpd")

	return nil
}
