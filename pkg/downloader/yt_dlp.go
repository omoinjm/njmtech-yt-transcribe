package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// execCommandFunc is a type that allows us to mock os/exec.Command in tests.
type execCommandFunc func(name string, arg ...string) *exec.Cmd

var commandExecutor execCommandFunc = exec.Command

// osLookPath is a variable that can be overridden for testing purposes.
var osLookPath = os.LookPath

// osStat is a variable that can be overridden for testing purposes.
var osStat = os.Stat

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
// yt-dlp -x --audio-format wav --output "/path/to/output/%(title)s.%(ext)s" <video-url>
// cmdCombinedOutput is a variable that can be overridden for testing purposes.
var cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

func (d *YTDLPAudioDownloader) DownloadAudio(videoURL string, outputDir string) (string, error) {
	// Check if yt-dlp is installed
	if _, err := osLookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp not found in PATH. Please install it to use this feature: %w", err)
	}

	cmd := commandExecutor(
		"yt-dlp",
		"-x", // Extract audio
		"--audio-format", "mp3", // Convert audio to mp3 format
		"--output", filepath.Join(outputDir, "%(title)s.%(ext)s"), // Output template
		"--print-to-stdout", // Print metadata to stdout (though we're parsing stderr for "Destination")
		"--restrict-filenames", // Keep filenames simple
		videoURL,
	)

	fmt.Printf("Executing command: %s\n", cmd.String())

	output, err := cmdCombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("yt-dlp command failed: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	var downloadedFilePath string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "[ExtractAudio] Destination: ") {
			downloadedFilePath = strings.TrimSpace(strings.TrimPrefix(line, "[ExtractAudio] Destination: "))
			break
		}
	}

	if downloadedFilePath == "" {
		return "", fmt.Errorf("could not determine downloaded file path from yt-dlp output:\n%s", outputStr)
	}

	if !filepath.IsAbs(downloadedFilePath) {
		downloadedFilePath = filepath.Join(outputDir, filepath.Base(downloadedFilePath))
	}

	// Verify the file exists
	if _, err := osStat(downloadedFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("yt-dlp reported successful download, but file not found at expected path: %s", downloadedFilePath)
	}

	return downloadedFilePath, nil
}
