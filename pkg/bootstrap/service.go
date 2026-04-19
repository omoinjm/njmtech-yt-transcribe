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

func loadDotEnv() {
	loadDotEnvOnce.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
		}
	})
}


func NewTranscriptionServiceFromEnv() (src.TranscriptionService, error) {
	loadDotEnv()
	ctx := context.Background()

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

	videoDownloader := downloader.NewYTDLPAudioDownloader()
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(whisperModelPath)
	blobUploader := uploader.NewVercelBlobUploader(vercelBlobAPIURL, vercelBlobAPIToken, &http.Client{})

	return src.NewTranscriptionService(videoDownloader, audioTranscriber, blobUploader), nil
}
