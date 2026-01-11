
package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UploadService handles the transcription upload process.
type UploadService struct {
	Uploader Uploader
}

// NewUploadService creates a new UploadService.
func NewUploadService(uploader Uploader) *UploadService {
	return &UploadService{
		Uploader: uploader,
	}
}

// SaveAndUpload saves the transcription to a file and uploads it.
func (s *UploadService) SaveAndUpload(transcription, videoURL string) (string, error) {
	platform := "other"
	if strings.Contains(videoURL, "youtube.com") {
		platform = "youtube"
	} else if strings.Contains(videoURL, "instagram.com") {
		platform = "instagram"
	}

	var transcriptPath string
	if platform == "instagram" {
		transcriptPath = filepath.Join("/tmp/njmtech-yt-transcribe", platform, "transript.txt")
	} else {
		transcriptPath = filepath.Join("/tmp/njmtech-yt-transcribe", "youtube", "transcript.txt")
	}

	if err := os.MkdirAll(filepath.Dir(transcriptPath), 0755); err != nil {
		return "", fmt.Errorf("error creating transcript directory: %w", err)
	}

	err := os.WriteFile(transcriptPath, []byte(transcription), 0644)
	if err != nil {
		return "", fmt.Errorf("error saving transcription to file: %w", err)
	}
	fmt.Printf("Transcription saved to: %s\n", transcriptPath)

	fmt.Println("Uploading transcription...")
	uploadPath := strings.ReplaceAll(transcriptPath, string(filepath.Separator), "/")
	uploadResponse, err := s.Uploader.Upload(transcription, uploadPath)
	if err != nil {
		return "", fmt.Errorf("error uploading transcription: %w", err)
	}

	return uploadResponse, nil
}
