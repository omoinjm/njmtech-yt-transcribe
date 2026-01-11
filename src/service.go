package src

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	APP_NAME           = "yt-transcribe"
	PLATFORM_OTHER     = "other"
	PLATFORM_YOUTUBE   = "youtube"
	PLATFORM_INSTAGRAM = "instagram"
	TRANSCRIPT_DIR     = "/tmp/njmtech-yt-transcribe"
)

// TranscriptionServiceImpl is the implementation of the TranscriptionService interface.
type TranscriptionServiceImpl struct {
	Downloader       VideoDownloader
	Transcriber      Transcriber
	Uploader         Uploader
	WhisperModelPath string
}

// NewTranscriptionService creates a new TranscriptionServiceImpl.
func NewTranscriptionService(downloader VideoDownloader, transcriber Transcriber, uploader Uploader, whisperModelPath string) TranscriptionService {
	return &TranscriptionServiceImpl{
		Downloader:       downloader,
		Transcriber:      transcriber,
		Uploader:         uploader,
		WhisperModelPath: whisperModelPath,
	}
}

// Execute orchestrates the download, transcription, and upload processes.
func (s *TranscriptionServiceImpl) Execute(videoURL, outputDir string) error {
	// 1. Download the audio
	fmt.Println("Downloading audio...")
	audioFilePath, videoID, err := s.Downloader.DownloadAudio(videoURL, outputDir)
	if err != nil {
		return fmt.Errorf("error downloading audio: %w", err)
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
	transcription, err := s.Transcriber.Transcribe(audioFilePath)
	if err != nil {
		return fmt.Errorf("error transcribing audio: %w", err)
	}

	// 3. Determine platform for upload path
	platform := PLATFORM_OTHER
	if strings.Contains(videoURL, "youtube.com") {
		platform = PLATFORM_YOUTUBE
	} else if strings.Contains(videoURL, "instagram.com") {
		platform = PLATFORM_INSTAGRAM
	}

	// 4. Upload the transcription
	fmt.Println("Uploading transcription...")
	// The upload path for vercel blob should be in the format: <app_name>/<platform>/<videoID>
	uploadPath := fmt.Sprintf("%s/%s/%s", APP_NAME, platform, videoID)
	uploadResponse, err := s.Uploader.Upload(transcription, uploadPath)
	if err != nil {
		return fmt.Errorf("error uploading transcription: %w", err)
	}

	fmt.Println("\n--- Transcription Upload Complete ---")
	fmt.Println("Response from Vercel Blob API:")
	fmt.Println(uploadResponse)

	return nil
}
