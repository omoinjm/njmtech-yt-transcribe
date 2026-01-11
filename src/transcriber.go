
package src

import (
	"fmt"
)

// TranscriptionProcessingService handles the audio transcription process.
type TranscriptionProcessingService struct {
	Transcriber Transcriber
}

// NewTranscriptionProcessingService creates a new TranscriptionProcessingService.
func NewTranscriptionProcessingService(transcriber Transcriber) *TranscriptionProcessingService {
	return &TranscriptionProcessingService{
		Transcriber: transcriber,
	}
}

// ProcessTranscription transcribes the audio from the given file path.
func (s *TranscriptionProcessingService) ProcessTranscription(audioFilePath string) (string, error) {
	fmt.Println("Transcribing audio...")
	transcription, err := s.Transcriber.Transcribe(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("error transcribing audio: %w", err)
	}
	return transcription, nil
}
