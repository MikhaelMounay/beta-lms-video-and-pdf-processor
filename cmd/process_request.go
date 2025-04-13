package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/utils"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/logger"
	pdfprocessor "github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/pdf_processor"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/shared"
	videoprocessor "github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/video_processor"
	"golang.design/x/clipboard"
)

func PromptNewFile(log *logger.Logger) *ProcessFileRequest {
	scanner := bufio.NewScanner(os.Stdin)
	pFRequest := &ProcessFileRequest{}

	// File Path
	log.Prompt("Enter the path of the video/pdf file : ")
	scanner.Scan()
	pFRequest.FilePath = strings.ReplaceAll(scanner.Text(), "\"", "")
	if pFRequest.FilePath == "" {
		return nil
	}

	// File Type
	fileType := filepath.Ext(pFRequest.FilePath)

	switch fileType {
	case ".mp4":
		pFRequest.FileType = VIDEOTYPE
	case ".pdf":
		pFRequest.FileType = PDFTYPE
	default:
		log.Log("Accepted file types: '.mp4' , '.pdf'\n")
		return nil
	}

	_, err := os.Stat(pFRequest.FilePath)
	if os.IsNotExist(err) {
		log.Log("Error: File does not exist!\n")
		return nil
	}

	// Grade
	log.Prompt("Enter the grade (10-11-12) : ")
	scanner.Scan()
	pFRequest.Grade = strings.TrimSpace(scanner.Text())
	if !slices.Contains([]string{"10", "11", "12"}, pFRequest.Grade) {
		log.Log("Accepted grade values : '10' , '11' , '12'\n")
		return nil
	}

	// File ID
	log.Prompt("Enter file ID ([Enter/Return] Generate a new ID) : ")
	scanner.Scan()
	pFRequest.FileId = strings.TrimSpace(scanner.Text())

	return pFRequest
}

func ProcessFile(log *logger.Logger, pFRequest *ProcessFileRequest) error {
	var file shared.IFile

	// Set file type, path, and Id
	if pFRequest.FileType == VIDEOTYPE {
		file = &videoprocessor.VideoFile{}
	} else if pFRequest.FileType == PDFTYPE {
		file = &pdfprocessor.PdfFile{}
	} else {
		return fmt.Errorf("invalid file type. Please try again")
	}

	file.SetGrade(pFRequest.Grade)
	file.SetPath(pFRequest.FilePath)
	file.SetId(pFRequest.FileId)

	// Process (encrypt/package) file
	if err := file.ProcessFile(); err != nil {
		return err
	}

	// Compute and log hash
	fileHash, err := utils.ComputeSHA256Hash(file.GetNewPath() + ".enc")
	if err != nil {
		log.Fatal("Error computing file hash for %s: %s\n", file.GetNewPath()+".enc", err.Error())
	}
	log.Logf("File ID:\t\033[1;32m%s\033[0m\t (copied to clipboard!)\n", file.GetId())
	log.Logf("File Hash:\t\033[1;32m%s\033[0m\n", fileHash)

	// Copy FileId to clipboard
	if err := clipboard.Init(); err != nil {
		panic(fmt.Sprintf("Error initializing clipboard: %v\n", err))
	}
	clipboard.Write(clipboard.FmtText, []byte(file.GetId()))

	// Add to upload queue
	filesToUpload := file.GetUploadPaths()
	for _, v := range filesToUpload {
		if config.BASE_B2_FOLDER == "" {
			fUploader.AddToQueue(v, fmt.Sprintf("stem-g%s/%s", pFRequest.Grade, filepath.Base(v)), file)
		} else {
			fUploader.AddToQueue(v, fmt.Sprintf("%sstem-g%s/%s", config.BASE_B2_FOLDER, pFRequest.Grade, filepath.Base(v)), file)
		}
	}

	return nil
}
