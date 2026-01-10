package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv" // Import godotenv
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/transcriber"
	"yt-transcribe/pkg/uploader"
)

func main() {
	// Load environment variables from .env file.
	// This should be done as early as possible.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
	}

	// Define command-line flags
	videoURL := flag.String("url", "https://www.youtube.com/watch?v=rdWZo5PD9Ek", "Video URL to download audio from")
	outputDir := flag.String("output", os.TempDir(), "Directory to save downloaded audio")
	flag.Parse()

	fmt.Printf("Transcribing video from URL: %s\n", *videoURL)
	fmt.Printf("Output directory: %s\n", *outputDir)

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory %s: %v", *outputDir, err)
	}

	// --- Dependency Injection setup ---
	// Initialize the Video Downloader
	videoDownloader := downloader.NewYTDLPAudioDownloader()

	// Initialize the WhisperCPP Transcriber
	whisperModelPath := os.Getenv("WHISPER_MODEL_PATH")
	if whisperModelPath == "" {
		log.Println("WHISPER_MODEL_PATH environment variable not set.")
	}
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(whisperModelPath)

	// Initialize the Vercel Blob Uploader
	vercelBlobAPIURL := os.Getenv("VERCEL_BLOB_API_URL")
	if vercelBlobAPIURL == "" {
		log.Fatalf("VERCEL_BLOB_API_URL environment variable not set.")
	}
	vercelBlobAPIToken := os.Getenv("VERCEL_BLOB_API_TOKEN")
	if vercelBlobAPIToken == "" {
		log.Fatalf("VERCEL_BLOB_API_TOKEN environment variable not set.")
	}
	blobUploader := uploader.NewVercelBlobUploader(vercelBlobAPIURL, vercelBlobAPIToken, &http.Client{})

	// --- Main application logic ---
	// 1. Download the audio
	fmt.Println("Downloading audio...")
	audioFilePath, err := videoDownloader.DownloadAudio(*videoURL, *outputDir)
	if err != nil {
		log.Fatalf("Error downloading audio: %v", err)
	}
	fmt.Printf("Audio downloaded to: %s\n", audioFilePath)
	defer func() {
		if err := os.Remove(audioFilePath); err != nil {
			log.Printf("Warning: could not remove temporary audio file %s: %v", audioFilePath, err)
		}
		fmt.Printf("Removed temporary audio file: %s\n", audioFilePath)
	}()

	// 2. Transcribe the audio
	fmt.Println("Transcribing audio...")

	transcription, err := audioTranscriber.Transcribe(audioFilePath)
	if err != nil {
		log.Fatalf("Error transcribing audio: %v", err)
	}

	// 3. Save the transcription to a local file
	// 3. Save the transcription to a local file and prepare for upload
	platform := "other"
	if strings.Contains(*videoURL, "youtube.com") {
		platform = "youtube"
	} else if strings.Contains(*videoURL, "instagram.com") {
		platform = "instagram"
	}

	dateTimeNow := time.Now().Format("2006-01-02_15:04:05") // More standard format for filename

	var transcriptPath string
	if platform == "instagram" {
		transcriptPath = filepath.Join("njmtech", platform, dateTimeNow, "transript")
	} else {
		// all other platforms will be saved under the youtube directory
		transcriptPath = filepath.Join("njmtech", "youtube", dateTimeNow, "transcript")
	}

	// Create local directory for the transcript file
	if err := os.MkdirAll(filepath.Dir(transcriptPath), 0755); err != nil {
		log.Fatalf("Error creating transcript directory: %v", err)
	}

	err = os.WriteFile(transcriptPath, []byte(transcription), 0644)
	if err != nil {
		log.Fatalf("Error saving transcription to file: %v", err)
	}
	fmt.Printf("Transcription saved to: %s\n", transcriptPath)

	// 4. Upload the transcription
	fmt.Println("Uploading transcription...")
	// The upload path for vercel blob should use forward slashes, even on windows.
	uploadPath := strings.ReplaceAll(transcriptPath, string(filepath.Separator), "/")
	uploadResponse, err := blobUploader.Upload(transcription, uploadPath)
	if err != nil {
		log.Fatalf("Error uploading transcription: %v", err)
	}

	fmt.Println("\n--- Transcription Upload Complete ---")
	fmt.Println("Response from Vercel Blob API:")
	fmt.Println(uploadResponse)
}
