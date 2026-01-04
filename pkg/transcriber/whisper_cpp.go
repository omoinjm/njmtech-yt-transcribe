package transcriber

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// WhisperCPPTranscriber implements the Transcriber interface using whisper.cpp.
type WhisperCPPTranscriber struct {
	ModelPath string
}

// NewWhisperCPPTranscriber creates a new WhisperCPPTranscriber.
func NewWhisperCPPTranscriber(modelPath string) *WhisperCPPTranscriber {
	return &WhisperCPPTranscriber{
		ModelPath: modelPath,
	}
}

// Transcribe transcribes the given audio file using whisper.cpp.
func (t *WhisperCPPTranscriber) Transcribe(audioFilePath string) (string, error) {
	// Create a temporary directory for output files
	tmpDir, err := os.MkdirTemp("", "whisper-transcript-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	outputPrefix := filepath.Join(tmpDir, "transcript")
	outputFilePath := outputPrefix + ".txt" // whisper-cli adds .txt extension

	// Construct the command
	cmdArgs := []string{
		"-m", t.ModelPath,
		"-f", audioFilePath,
		"--output-txt",
		"--output-file", outputPrefix,
		"--no-prints",
	}

	cmd := exec.Command("whisper-cli", cmdArgs...)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute whisper-cli: %w\nOutput: %s", err, output)
	}

	// Read the transcribed text from the output file
	transcriptBytes, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file %s: %w", outputFilePath, err)
	}

	return string(transcriptBytes), nil
}
