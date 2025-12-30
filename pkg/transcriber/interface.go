package transcriber

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
