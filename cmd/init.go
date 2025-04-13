package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/config"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/logger"
	"github.com/MikhaelMounay/beta-lms-video-and-pdf-processor/internal/uploader"
	"github.com/joho/godotenv"
	"golang.org/x/term"
)

var (
	log       *logger.Logger
	fUploader uploader.FileUploader
)

func init() {
	log = logger.NewLogger(os.Stdout)

	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v\n", err))
	}

	// Get keys from env file
	config.SECRET_KEY = os.Getenv("ENC_SECRET_KEY_HEX")
	config.IV = os.Getenv("ENC_IV_HEX")
	config.VOD_KEY_ID = os.Getenv("VOD_KEY_ID")
	config.VOD_KEY = os.Getenv("VOD_KEY")
	config.INSTANCE_NAME = os.Getenv("INSTANCE_NAME")
	config.BASE_B2_FOLDER = os.Getenv("BASE_B2_FOLDER")
	config.ENDPOINT = os.Getenv("ENDPOINT")
	config.APP_KEY_ID = os.Getenv("B2_APP_KEY_ID")
	config.APP_KEY = os.Getenv("B2_APP_KEY")
	config.BUCKET_ID = os.Getenv("B2_BUCKET_ID")
	config.BUCKET_NAME = os.Getenv("B2_BUCKET_NAME")
	config.BUCKET_REGION = os.Getenv("B2_BUCKET_REGION")
	config.DATABASE_URL = os.Getenv("DATABASE_URL")

	// Manually initialize the variables
	// config.SECRET_KEY = ""
	// config.IV = ""
	// config.VOD_KEY_ID = ""
	// config.VOD_KEY = ""
	// config.INSTANCE_NAME = ""
	// config.BASE_B2_FOLDER =""
	// config.ENDPOINT = ""
	// config.APP_KEY_ID = ""
	// config.APP_KEY = ""
	// config.BUCKET_ID = ""
	// config.BUCKET_NAME = ""
	// config.BUCKET_REGION = ""
	// config.DATABASE_URL = ""

	fUploader = uploader.ConfigUploader(log)

	// Print welcome message
	printWelcomeMsg()
}

func printWelcomeMsg() {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback if we can't get the terminal size
		panic(err)
	}

	dashedLine := func() string {
		return strings.Repeat(" ", termWidth) // With white space, it's acting as a '\n'
	}

	dashedLineWithText := func(text string) string {
		text = "  " + text + "  "
		textLength := len(text)

		padding := (termWidth - textLength) / 2

		return strings.Repeat("-", padding) + text + strings.Repeat("-", padding+termWidth-textLength-2*padding)
	}

	fmt.Println(dashedLine())
	fmt.Println(dashedLineWithText("Beta LMS"))
	fmt.Println(dashedLine())
	fmt.Println(dashedLineWithText("The Only Secure LMS You Need"))
	fmt.Println(dashedLine())
	fmt.Println(dashedLineWithText("Welcome to Beta LMS Video & PDF Processor! (v2.1.1)"))
	fmt.Println(dashedLine())
	fmt.Println(dashedLineWithText("Made with <3 by Transcendea Software"))
	fmt.Println(dashedLine())
	fmt.Println(dashedLineWithText(fmt.Sprintf("Welcome to Beta LMS Instance: %s", config.INSTANCE_NAME)))
	fmt.Println(dashedLine())
	fmt.Println()
}
