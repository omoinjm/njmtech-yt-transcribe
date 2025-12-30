package transcriber

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// MockAPIKeyProvider is a mock implementation of the APIKeyProvider interface for testing.
type MockAPIKeyProvider struct {
	Key string
}

// GetAPIKey returns the predefined mock key.
func (m *MockAPIKeyProvider) GetAPIKey() string {
	return m.Key
}

// TestNewOpenAITranscriber ensures that the constructor returns a non-nil instance
// and correctly sets the APIKeyProvider.
func TestNewOpenAITranscriber(t *testing.T) {
	mockProvider := &MockAPIKeyProvider{Key: "test-api-key"}
	transcriber := NewOpenAITranscriber(mockProvider)
	if transcriber == nil {
		t.Errorf("NewOpenAITranscriber returned nil, expected an instance")
	}
	if transcriber.apiKeyProvider != mockProvider {
		t.Errorf("NewOpenAITranscriber did not set the correct APIKeyProvider")
	}

	// Test with nil provider
	nilTranscriber := NewOpenAITranscriber(nil)
	if nilTranscriber == nil {
		t.Errorf("NewOpenAITranscriber returned nil with nil provider, expected an instance (with warning)")
	}
	if nilTranscriber.apiKeyProvider != nil {
		t.Errorf("NewOpenAITranscriber with nil provider should have nil apiKeyProvider")
	}
}

// TestOpenAITranscriber_Transcribe tests the mock transcription functionality
// with various API key provider scenarios.
func TestOpenAITranscriber_Transcribe(t *testing.T) {
	audioFilePath := "/tmp/test_audio.mp3" // Use a dummy path for the mock
	expectedFileName := filepath.Base(audioFilePath)

	tests := []struct {
		name         string
		apiKey       string
		expectError  bool
		expectedErrMsg string
		expectedSubstring string
	}{
		{
			name:         "Valid API Key",
			apiKey:       "sk-test1234567890",
			expectError:  false,
			expectedSubstring: fmt.Sprintf("Mock Transcriber: Simulating transcription for %s using API key (first 5 chars): sk-te...", audioFilePath),
		},
		{
			name:         "Empty API Key",
			apiKey:       "",
			expectError:  true,
			expectedErrMsg: "API key not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &MockAPIKeyProvider{Key: tt.apiKey}
			transcriber := NewOpenAITranscriber(mockProvider)

			transcription, err := transcriber.Transcribe(audioFilePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
				if !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error message to contain '%s', got: %v", tt.expectedErrMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if !strings.Contains(transcription, expectedFileName) {
					t.Errorf("Transcription did not contain expected filename.\nGot: %s", transcription)
				}
				if !strings.Contains(transcription, "The quick brown fox jumps over the lazy dog.") {
					t.Errorf("Transcription did not contain the expected filler text.")
				}
				// Verify API key snippet in mock output
				if !strings.Contains(transcription, tt.expectedSubstring) {
					t.Errorf("Mock output did not contain expected API key snippet.\nExpected to contain: %s\nGot: %s", tt.expectedSubstring, transcription)
				}
			}
		})
	}

	t.Run("Nil APIKeyProvider during Transcribe", func(t *testing.T) {
		transcriber := NewOpenAITranscriber(nil) // Initialize with a nil provider
		_, err := transcriber.Transcribe(audioFilePath)
		if err == nil {
			t.Errorf("Expected an error when APIKeyProvider is nil during Transcribe, but got none")
		}
		if !strings.Contains(err.Error(), "transcriber is not properly initialized: missing APIKeyProvider") {
			t.Errorf("Expected specific error message for nil APIKeyProvider, got: %v", err)
		}
	})
}