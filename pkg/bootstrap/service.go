package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/transcriber"
	"yt-transcribe/pkg/uploader"
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

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s environment variable not set", name)
	}

	return value, nil
}

func NewTranscriptionServiceFromEnv() (src.TranscriptionService, error) {
	loadDotEnv()

	whisperModelPath, err := requiredEnv("WHISPER_MODEL_PATH")
	if err != nil {
		return nil, err
	}

	vercelBlobAPIURL, err := requiredEnv("VERCEL_BLOB_API_URL")
	if err != nil {
		return nil, err
	}

	vercelBlobAPIToken, err := requiredEnv("VERCEL_BLOB_API_TOKEN")
	if err != nil {
		return nil, err
	}

	videoDownloader := downloader.NewYTDLPAudioDownloader()
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(whisperModelPath)
	blobUploader := uploader.NewVercelBlobUploader(vercelBlobAPIURL, vercelBlobAPIToken, &http.Client{})

	return src.NewTranscriptionService(videoDownloader, audioTranscriber, blobUploader), nil
}
