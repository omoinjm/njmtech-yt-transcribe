package transcriber

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

// OllamaTranscriber is a transcriber that uses the Ollama API.
type OllamaTranscriber struct {
	client    *api.Client
	modelName string
}

// NewOllamaTranscriber creates a new OllamaTranscriber.
func NewOllamaTranscriber(ollamaHost, modelName string) (*OllamaTranscriber, error) {
	if ollamaHost == "" {
		return nil, fmt.Errorf("Ollama host is not set")
	}

	ollamaURL, err := url.Parse(ollamaHost)
	if err != nil {
		return nil, fmt.Errorf("invalid Ollama host URL: %w", err)
	}

	client, err := api.NewClient(ollamaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return &OllamaTranscriber{
		client:    client,
		modelName: modelName,
	}, nil
}

// Transcribe transcribes the audio file at the given path.
func (t *OllamaTranscriber) Transcribe(audioFilePath string) (string, error) {
	audioData, err := os.ReadFile(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read audio file: %w", err)
	}

	req := &api.GenerateRequest{
		Model:  t.modelName,
		Prompt: "transcribe the following audio file:",
		Images: []api.ImageData{api.ImageData(audioData)},
	}

	var transcription string
	respCh := make(chan *api.GenerateResponse)
	errCh := make(chan error)

	go func() {
		defer close(respCh)
		defer close(errCh)
		if err := t.client.Generate(context.Background(), req, func(resp api.GenerateResponse) error {
			respCh <- &resp
			return nil
		}); err != nil {
			errCh <- err
		}
	}()

	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				return transcription, nil
			}
			transcription += resp.Response
		case err, ok := <-errCh:
			if !ok {
				return transcription, nil
			}
			if err != nil {
				return "", fmt.Errorf("Ollama transcription failed: %w", err)
			}
		}
	}
}
