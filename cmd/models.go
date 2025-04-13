package main

import "github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/uploader"

type FileType int

const (
	VIDEOTYPE FileType = iota
	PDFTYPE
)

type ProcessFileRequest struct {
	Grade    string
	FilePath string
	FileType FileType
	FileId   string
	Uploader uploader.FileUploader
}
