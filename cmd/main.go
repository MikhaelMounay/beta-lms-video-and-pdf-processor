package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/history"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	closeFunc := config.ConnectDB(log)
	defer closeFunc()

	i := 0
	for {
		i++

		// Prompt user to choose action
		log.Prompt("(1) Process a file\n")
		log.Promptf("(2) Upload files (%v in queue) and exit\n", fUploader.GetQueueSize())
		if i == 1 {
			log.Prompt("(3) Fetch History and exit\n")
		}
		log.Prompt("Enter required command number : ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			pFRequest := PromptNewFile(log)
			if pFRequest == nil {
				log.Log("Invalid input. Please try again.\n")
				continue
			}

			if err := ProcessFile(log, pFRequest); err != nil {
				log.Logf("Couldn't process the given file: %s\n", err.Error())
			}
		case "2":
			fUploader.StartQueueUpload()

			log.Log("Files Uploaded Successfully. Make sure to copy file \033[1;32mIDs\033[0m and \033[1;32mHashes\033[0m above!\n")

			fmt.Print("\nThank you for using BetaLMS! Made with <3 by Transcendea Software.\n\nPress [Enter/Return] to exit... ")
			scanner.Scan()
			return
		case "3":
			if i == 1 {
				log.Prompt("Enter desired CSV file output path ([Enter/Return] directly to use the default): ")
				scanner.Scan()
				csvFilePath := strings.TrimSpace(scanner.Text())
				// Validating/Sanitizing input
				if csvFilePath == "" {
					csvFilePath = "upload-history_" + config.INSTANCE_NAME + ".BetaLMS-TranscendeaSoftware_" + time.Now().Format("2006-01-02_03-04pm") + ".csv"
				} else if filepath.Ext(csvFilePath) != ".csv" {
					csvFilePath += ".csv"
				}
				csvFilePath = strings.ReplaceAll(csvFilePath, " ", "_")

				if err := history.SaveHistoryCSV(csvFilePath); err != nil {
					log.Fatal("Error processing upload history: %s\n", err.Error())
				}

				log.Logf("Upload history saved successfully to %s\n", csvFilePath)

				fmt.Print("\nThank you for using BetaLMS! Made with <3 by Transcendea Software.\n\nPress [Enter/Return] to exit... ")
				scanner.Scan()
				return
			}

			fallthrough
		default:
			log.Log("Invalid choice. Please try again.\n")
			continue
		}
	}
}
