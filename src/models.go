
package src

// VideoDownloader defines the interface for downloading audio from videos.
// Applying the Interface Segregation Principle (ISP) and Dependency Inversion Principle (DIP).
type VideoDownloader interface {
	DownloadAudio(videoURL string, outputDir string) (filePath string, videoID string, err error)
}

// Transcriber defines the interface for transcribing audio files into text.
// This adheres to the Interface Segregation Principle (ISP) and Dependency Inversion Principle (DIP).
type Transcriber interface {
	Transcribe(audioFilePath string) (string, error)
}

// APIKeyProvider defines the interface for retrieving an API key.
// This allows the Transcriber to depend on an abstraction for fetching credentials,
// rather than a concrete implementation like os.Getenv, adhering to DIP and SRP.
type APIKeyProvider interface {
	GetAPIKey() string
}

// Uploader defines the interface for uploading content.
type Uploader interface {
	// Upload takes content as a string and uploads it, returning the response from the upload service.
	Upload(content string, filename string) (string, error)
}

// TranscriptionService defines the interface for the main transcription service.
type TranscriptionService interface {
	Execute(videoURL, outputDir string) error
}
