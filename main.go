package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv" // Import godotenv
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/transcriber" // Import the transcriber package
)

func main() {
	// Load environment variables from .env file.
	// This should be done as early as possible.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
	}

	// Define command-line flags
	videoURL := flag.String("url", "", "YouTube video URL to transcribe")
	outputDir := flag.String("output", os.TempDir(), "Directory to save downloaded audio and transcription")
	flag.Parse()

	// Validate the video URL
	if *videoURL == "" {
		fmt.Println("Error: YouTube video URL is required.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Transcribing YouTube video from URL: %s\n", *videoURL)
	fmt.Printf("Output directory: %s\n", *outputDir)

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory %s: %v", *outputDir, err)
	}

	// --- Dependency Injection setup ---
	// Initialize the YouTube Downloader
	ytDownloader := downloader.NewYTDLPAudioDownloader()

	// Initialize the API Key Provider
	// This uses an environment variable, adhering to secure credential handling.
	apiKeyProvider := transcriber.NewEnvAPIKeyProvider("OPENAI_API_KEY")

	// Initialize the Transcriber with the injected API Key Provider.
	audioTranscriber := transcriber.NewOpenAITranscriber(apiKeyProvider)

	// --- Main application logic ---
	// 1. Download the audio
	fmt.Println("Downloading audio...")
	audioFilePath, err := ytDownloader.DownloadAudio(*videoURL, *outputDir)
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

	// 3. Output the transcription
	outputFileName := filepath.Join(*outputDir, fmt.Sprintf("%s.txt", sanitizeFilename(filepath.Base(*videoURL))))
	if err := os.WriteFile(outputFileName, []byte(transcription), 0644); err != nil {
		log.Fatalf("Error writing transcription to file %s: %v", outputFileName, err)
	}

	fmt.Println("\n--- Transcription Complete ---")
	fmt.Printf("Transcription saved to: %s\n", outputFileName)
	fmt.Println("Content:")
	fmt.Println(transcription)
}

// sanitizeFilename removes characters that are not safe for filenames.
func sanitizeFilename(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "*", "_")
	s = strings.ReplaceAll(s, "?", "_")
	s = strings.ReplaceAll(s, "\"", "_")
	s = strings.ReplaceAll(s, "<", "_")
	s = strings.ReplaceAll(s, ">", "_")
	s = strings.ReplaceAll(s, "|", "_")
	return s
}
