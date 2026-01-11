
package src

import (
	"fmt"
	"log"
	"os"
)

// DownloadService handles the audio download process.
type DownloadService struct {
	Downloader VideoDownloader
}

// NewDownloadService creates a new DownloadService.
func NewDownloadService(downloader VideoDownloader) *DownloadService {
	return &DownloadService{
		Downloader: downloader,
	}
}

// DownloadAudio downloads the audio from the given URL.
func (s *DownloadService) DownloadAudio(videoURL, outputDir string) (string, error) {
	fmt.Println("Downloading audio...")
	audioFilePath, err := s.Downloader.DownloadAudio(videoURL, outputDir)
	if err != nil {
		return "", fmt.Errorf("error downloading audio: %w", err)
	}
	fmt.Printf("Audio downloaded to: %s\n", audioFilePath)

	// Defer removal of the temporary audio file
	defer func() {
		if err := os.Remove(audioFilePath); err != nil {
			log.Printf("Warning: could not remove temporary audio file %s: %v", audioFilePath, err)
		}
		fmt.Printf("Removed temporary audio file: %s\n", audioFilePath)
	}()

	return audioFilePath, nil
}

