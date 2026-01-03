package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	// "time"
)

// execCommandFunc is a type that allows us to mock os/exec.Command in tests.
type execCommandFunc func(name string, arg ...string) *exec.Cmd

var commandExecutor execCommandFunc = exec.Command

// osLookPath is a variable that can be overridden for testing purposes.
var osLookPath = exec.LookPath

// osStat is a variable that can be overridden for testing purposes.
var osStat = os.Stat

// cmdCombinedOutput is a variable that can be overridden for testing purposes.
var cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

// YTDLPAudioDownloader is an implementation of YouTubeDownloader that uses the `yt-dlp` external tool.
// It downloads the audio stream of a given YouTube video.
type YTDLPAudioDownloader struct{}

// NewYTDLPAudioDownloader creates and returns a new instance of YTDLPAudioDownloader.
// This acts as a constructor, promoting consistency in object creation.
func NewYTDLPAudioDownloader() *YTDLPAudioDownloader {
	return &YTDLPAudioDownloader{}
}

// DownloadAudio downloads the audio stream from the specified YouTube video URL
// and saves it to the given output directory.
// It returns the full path to the downloaded audio file or an error if the download fails.
//
// Dependencies: This function relies on the `yt-dlp` command-line tool being installed
// and accessible in the system's PATH.
//
// Example yt-dlp command:
// yt-dlp -x --audio-format wav --output "/path/to/output/video_2025-12-31.wav" <video-url>
func (d *YTDLPAudioDownloader) DownloadAudio(videoURL string, outputDir string) (string, error) {
	// Check if ffmpeg is installed
	if _, err := osLookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH. ffmpeg is required by yt-dlp to process audio. Please install it to use this feature: %w", err)
	}

	// Check if yt-dlp is installed
	if _, err := osLookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp not found in PATH. Please install it to use this feature: %w", err)
	}

	// Generate filename based on current date
	// dateStr := time.Now().Format("2006-01-02")
	// outputFilename := fmt.Sprintf("video_%s.wav", dateStr)
	outputFilename := "audio.wav"

	cmd := commandExecutor(
		"yt-dlp",
		"-x",                    // Extract audio
		"--audio-format", "wav", // Convert audio to wav format
		"--output", filepath.Join(outputDir, outputFilename), // Output template
		"--restrict-filenames", // Keep filenames simple
		videoURL,
	)

	fmt.Printf("Executing command: %s\n", cmd.String())

	output, err := cmdCombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("yt-dlp command failed: %v\nOutput: %s", err, string(output))
	}

	// The output from yt-dlp now directly gives us the filename
	downloadedFilePath := filepath.Join(outputDir, outputFilename)

	// Verify the file exists
	if _, err := osStat(downloadedFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("yt-dlp reported successful download, but file not found at expected path: %s", downloadedFilePath)
	}

	return downloadedFilePath, nil
}
