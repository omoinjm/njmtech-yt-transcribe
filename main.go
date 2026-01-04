package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"

	"github.com/joho/godotenv" // Import godotenv
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/transcriber"
)

func main() {
	// Load environment variables from .env file.
	// This should be done as early as possible.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
	}

	// Define command-line flags
	videoURL := flag.String("url", "https://www.youtube.com/watch?v=rdWZo5PD9Ek", "YouTube video URL to download audio from")
	outputDir := flag.String("output", os.TempDir(), "Directory to save downloaded audio")
	flag.Parse()

	fmt.Printf("Transcribing YouTube video from URL: %s\n", *videoURL)
	fmt.Printf("Output directory: %s\n", *outputDir)

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory %s: %v", *outputDir, err)
	}

	// --- Dependency Injection setup ---
	// Initialize the YouTube Downloader
	ytDownloader := downloader.NewYTDLPAudioDownloader()

	// Initialize the WhisperCPP Transcriber
	whisperModelPath := os.Getenv("WHISPER_MODEL_PATH")
	if whisperModelPath == "" {
		log.Println("WHISPER_MODEL_PATH environment variable not set. Using default path: /whisper.cpp/models/ggml-base.en.bin")
		whisperModelPath = "/whisper.cpp/models/ggml-base.en.bin"
	}
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(whisperModelPath)

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
