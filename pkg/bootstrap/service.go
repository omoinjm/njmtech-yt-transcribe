package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/transcriber"
	"yt-transcribe/pkg/uploader"
	"yt-transcribe/pkg/secrets"
	"yt-transcribe/src"
)

var loadDotEnvOnce sync.Once

// Config holds application configuration loaded from environment or Infisical.
type Config struct {
	WhisperModelPath  string
	VercelBlobAPIURL  string
	VercelBlobAPIToken string
	PostgresURL       string
}

func loadDotEnv() {
	loadDotEnvOnce.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
		}
	})
}

// LoadConfigFromEnv loads all configuration from environment or Infisical.
func LoadConfigFromEnv(ctx context.Context) (*Config, error) {
	loadDotEnv()

	// Optional Infisical configuration: set INFISICAL_ENABLED=true, INFISICAL_PROJECT_ID and INFISICAL_ENVIRONMENT
	infisicalProjectID := os.Getenv("INFISICAL_PROJECT_ID")
	infisicalEnvironment := os.Getenv("INFISICAL_ENVIRONMENT")

	whisperModelPath, err := secrets.GetSecret(ctx, "WHISPER_MODEL_PATH", "WHISPER_MODEL_PATH", infisicalProjectID, infisicalEnvironment)
	if err != nil {
		return nil, err
	}
	if whisperModelPath == "" {
		return nil, fmt.Errorf("WHISPER_MODEL_PATH not set")
	}

	vercelBlobAPIURL, err := secrets.GetSecret(ctx, "VERCEL_BLOB_API_URL", "VERCEL_BLOB_API_URL", infisicalProjectID, infisicalEnvironment)
	if err != nil {
		return nil, err
	}
	if vercelBlobAPIURL == "" {
		return nil, fmt.Errorf("VERCEL_BLOB_API_URL not set")
	}

	vercelBlobAPIToken, err := secrets.GetSecret(ctx, "VERCEL_BLOB_API_TOKEN", "VERCEL_BLOB_API_TOKEN", infisicalProjectID, infisicalEnvironment)
	if err != nil {
		return nil, err
	}
	if vercelBlobAPIToken == "" {
		return nil, fmt.Errorf("VERCEL_BLOB_API_TOKEN not set")
	}

	// POSTGRES_URL is optional (only needed for -db mode)
	postgresURL, err := secrets.GetSecret(ctx, "POSTGRES_URL", "POSTGRES_URL", infisicalProjectID, infisicalEnvironment)
	if err != nil {
		// If POSTGRES_URL is not found and Infisical is not enabled, return empty string (optional)
		postgresURL = ""
	}

	return &Config{
		WhisperModelPath:   whisperModelPath,
		VercelBlobAPIURL:   vercelBlobAPIURL,
		VercelBlobAPIToken: vercelBlobAPIToken,
		PostgresURL:        postgresURL,
	}, nil
}

func NewTranscriptionServiceFromEnv() (src.TranscriptionService, error) {
	ctx := context.Background()
	cfg, err := LoadConfigFromEnv(ctx)
	if err != nil {
		return nil, err
	}

	videoDownloader := downloader.NewYTDLPAudioDownloader()
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(cfg.WhisperModelPath)
	blobUploader := uploader.NewVercelBlobUploader(cfg.VercelBlobAPIURL, cfg.VercelBlobAPIToken, &http.Client{})

	return src.NewTranscriptionService(videoDownloader, audioTranscriber, blobUploader), nil
}
