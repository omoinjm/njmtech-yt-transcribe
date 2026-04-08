package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"yt-transcribe/pkg/downloader"
	"yt-transcribe/pkg/repository"
	"yt-transcribe/pkg/transcriber"
	"yt-transcribe/pkg/uploader"
	"yt-transcribe/src"
)

const (
	DEFAULT_VIDEO_URL = "https://www.youtube.com/watch?v=rdWZo5PD9Ek"
	URL_FLAG          = "url"
	OUTPUT_FLAG       = "output"
	DB_FLAG           = "db"
	REPROCESS_ALL_FLAG = "reprocess-all"
)

// handleFatalError logs a fatal error and exits the program.
func handleFatalError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
	log.Fatal(message)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found or error loading .env file. Proceeding without .env variables.")
	}

	// Define command-line flags
	videoURL := flag.String(URL_FLAG, "", "Video URL to download audio from. Can also be provided as a positional argument.")
	outputDir := flag.String(OUTPUT_FLAG, os.TempDir(), "Directory to save downloaded audio")
	useDB        := flag.Bool(DB_FLAG, false, "Fetch the next unprocessed video URL from the database instead of using -url")
	reprocessAll := flag.Bool(REPROCESS_ALL_FLAG, false, "Re-transcribe every record in the database, overwriting existing transcript URLs")
	flag.Parse()

	// --- Shared dependency setup ---
	whisperModelPath := os.Getenv("WHISPER_MODEL_PATH")
	if whisperModelPath == "" {
		handleFatalError("WHISPER_MODEL_PATH environment variable not set", nil)
	}

	vercelBlobAPIURL := os.Getenv("VERCEL_BLOB_API_URL")
	if vercelBlobAPIURL == "" {
		handleFatalError("VERCEL_BLOB_API_URL environment variable not set", nil)
	}
	vercelBlobAPIToken := os.Getenv("VERCEL_BLOB_API_TOKEN")
	if vercelBlobAPIToken == "" {
		handleFatalError("VERCEL_BLOB_API_TOKEN environment variable not set", nil)
	}

	videoDownloader := downloader.NewYTDLPAudioDownloader()
	audioTranscriber := transcriber.NewWhisperCPPTranscriber(whisperModelPath)
	blobUploader := uploader.NewVercelBlobUploader(vercelBlobAPIURL, vercelBlobAPIToken, &http.Client{})
	transcriptionService := src.NewTranscriptionService(videoDownloader, audioTranscriber, blobUploader)

	ctx := context.Background()

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		handleFatalError(fmt.Sprintf("Error creating output directory %s", *outputDir), err)
	}

	if *reprocessAll {
		runReprocessAll(ctx, transcriptionService, *outputDir)
	} else if *useDB {
		runFromDB(ctx, transcriptionService, *outputDir)
	} else {
		runFromCLI(ctx, transcriptionService, *videoURL, *outputDir)
	}
}

// runFromCLI processes a single URL provided via flags or positional args.
func runFromCLI(ctx context.Context, svc src.TranscriptionService, videoURL, outputDir string) {
	if videoURL == "" {
		if len(flag.Args()) > 0 {
			videoURL = flag.Args()[0]
		} else {
			log.Printf("No URL provided. Using default URL: %s", DEFAULT_VIDEO_URL)
			videoURL = DEFAULT_VIDEO_URL
		}
	}

	if _, err := url.ParseRequestURI(videoURL); err != nil {
		handleFatalError(fmt.Sprintf("Error: Invalid video URL provided: %s", videoURL), err)
	}

	fmt.Printf("Transcribing video from URL: %s\n", videoURL)
	fmt.Printf("Output directory: %s\n", outputDir)

	if _, err := svc.Execute(ctx, videoURL, outputDir); err != nil {
		handleFatalError("Error executing transcription service", err)
	}
}

// runFromDB fetches the next unprocessed media_items row, transcribes it,
// and writes the resulting Vercel Blob URL back to transcript_url.
func runFromDB(ctx context.Context, svc src.TranscriptionService, outputDir string) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		handleFatalError("POSTGRES_URL environment variable not set (required for -db mode)", nil)
	}

	repo, err := repository.NewPostgresMediaItemRepository(ctx, postgresURL)
	if err != nil {
		handleFatalError("Failed to connect to database", err)
	}
	defer repo.Close(ctx)

	item, err := repo.FetchNextUnprocessed(ctx)
	if err != nil {
		handleFatalError("Failed to fetch next unprocessed item", err)
	}
	if item == nil {
		fmt.Println("No unprocessed items found in the database. Nothing to do.")
		return
	}

	fmt.Printf("Fetched item from DB — id: %s  platform: %s  url: %s\n", item.ID, item.Platform, item.URL)
	fmt.Printf("Output directory: %s\n", outputDir)

	blobURL, err := svc.Execute(ctx, item.URL, outputDir)
	if err != nil {
		handleFatalError("Error executing transcription service", err)
	}

	if err := repo.UpdateTranscriptURL(ctx, item.ID, blobURL); err != nil {
		handleFatalError("Transcription succeeded but failed to update transcript_url in database", err)
	}

	fmt.Printf("transcript_url updated in database for id %s\n", item.ID)
}

// runReprocessAll fetches every record in media_items and re-transcribes each one,
// overwriting the existing transcript_url. Failures on individual items are logged
// and skipped so the rest of the batch can continue.
func runReprocessAll(ctx context.Context, svc src.TranscriptionService, outputDir string) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		handleFatalError("POSTGRES_URL environment variable not set (required for -reprocess-all mode)", nil)
	}

	repo, err := repository.NewPostgresMediaItemRepository(ctx, postgresURL)
	if err != nil {
		handleFatalError("Failed to connect to database", err)
	}
	defer repo.Close(ctx)

	items, err := repo.FetchAll(ctx)
	if err != nil {
		handleFatalError("Failed to fetch all items from database", err)
	}
	if len(items) == 0 {
		fmt.Println("No records found in the database. Nothing to do.")
		return
	}

	total := len(items)
	succeeded, failed := 0, 0

	fmt.Printf("Reprocessing %d record(s)...\n\n", total)

	for i, item := range items {
		fmt.Printf("[%d/%d] id: %s  platform: %s  url: %s\n", i+1, total, item.ID, item.Platform, item.URL)

		blobURL, err := svc.Execute(ctx, item.URL, outputDir)
		if err != nil {
			log.Printf("  ✗ transcription failed: %v — skipping\n", err)
			failed++
			continue
		}

		if err := repo.UpdateTranscriptURL(ctx, item.ID, blobURL); err != nil {
			log.Printf("  ✗ db update failed: %v — skipping\n", err)
			failed++
			continue
		}

		fmt.Printf("  ✓ transcript_url updated\n")
		succeeded++
	}

	fmt.Printf("\nDone. %d succeeded, %d failed out of %d total.\n", succeeded, failed, total)
}
