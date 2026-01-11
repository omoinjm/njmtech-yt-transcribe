package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv" // Import godotenv
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/transcriber"
	"yt-transcribe/pkg/uploader"
	"yt-transcribe/src"
)

const (
	DEFAULT_VIDEO_URL = "https://www.youtube.com/watch?v=rdWZo5PD9Ek"
	URL_FLAG          = "url"
	OUTPUT_FLAG       = "output"
)

// handleFatalError logs a fatal error and exits the program.
func handleFatalError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
	log.Fatal(message)
}

func main() {
	// Load environment variables from .env file.
	// This should be done as early as possible.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
	}

	// Define command-line flags
	videoURL := flag.String(URL_FLAG, "", "Video URL to download audio from. Can also be provided as a positional argument.")
	outputDir := flag.String(OUTPUT_FLAG, os.TempDir(), "Directory to save downloaded audio")
	flag.Parse()

	// If no URL is provided via flag, check for a positional argument
	if *videoURL == "" {
		if len(flag.Args()) > 0 {
			*videoURL = flag.Args()[0]
		} else {
			*videoURL = DEFAULT_VIDEO_URL
		}
	}

	// Validate the video URL
	if _, err := url.ParseRequestURI(*videoURL); err != nil {
		handleFatalError(fmt.Sprintf("Error: Invalid video URL provided: %s", *videoURL), err)
	}

	fmt.Printf("Transcribing video from URL: %s\n", *videoURL)
	fmt.Printf("Output directory: %s\n", *outputDir)

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		handleFatalError(fmt.Sprintf("Error creating output directory %s", *outputDir), err)
	}

	// --- Dependency Injection setup ---
	videoDownloader := downloader.NewYTDLPAudioDownloader()
	whisperModelPath := os.Getenv("WHISPER_MODEL_PATH")
	if whisperModelPath == "" {
		log.Println("WHISPER_MODEL_PATH environment variable not set.")
	}
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(whisperModelPath)
	vercelBlobAPIURL := os.Getenv("VERCEL_BLOB_API_URL")
	if vercelBlobAPIURL == "" {
		handleFatalError("VERCEL_BLOB_API_URL environment variable not set", nil)
	}
	vercelBlobAPIToken := os.Getenv("VERCEL_BLOB_API_TOKEN")
	if vercelBlobAPIToken == "" {
		handleFatalError("VERCEL_BLOB_API_TOKEN environment variable not set", nil)
	}
	blobUploader := uploader.NewVercelBlobUploader(vercelBlobAPIURL, vercelBlobAPIToken, &http.Client{})

	// Initialize the transcription service
	transcriptionService := src.NewTranscriptionService(videoDownloader, audioTranscriber, blobUploader, whisperModelPath)

	// Execute the transcription service
	if err := transcriptionService.Execute(*videoURL, *outputDir); err != nil {
		handleFatalError("Error executing transcription service", err)
	}
}
