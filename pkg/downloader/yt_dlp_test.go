package downloader

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// A mutex to protect global variables during concurrent tests
var testMux sync.Mutex

// Helper function to restore original globals after each test
func restoreGlobals(
	oldCommandExecutor execCommandFunc,
	oldOsLookPath func(string) (string, error),
	oldOsStat func(string) (os.FileInfo, error),
	oldCmdCombinedOutput func(cmd *exec.Cmd) ([]byte, error),
) {
	testMux.Lock()
	defer testMux.Unlock()
	commandExecutor = oldCommandExecutor
	osLookPath = oldOsLookPath
	osStat = oldOsStat
	cmdCombinedOutput = oldCmdCombinedOutput
}

// TestNewYTDLPAudioDownloader ensures the constructor works correctly.
func TestNewYTDLPAudioDownloader(t *testing.T) {
	downloader := NewYTDLPAudioDownloader()
	if downloader == nil {
		t.Errorf("NewYTDLPAudioDownloader returned nil, expected an instance")
	}
}

// TestDownloadAudio_FFMPEGNotFound tests the scenario where ffmpeg is not found.
func TestDownloadAudio_FFMPEGNotFound(t *testing.T) {
	testMux.Lock()
	oldOsLookPath := osLookPath
	osLookPath = func(file string) (string, error) {
		if file == "ffmpeg" {
			return "", errors.New("not found in PATH")
		}
		return "/usr/local/bin/yt-dlp", nil // Assume yt-dlp is found
	}
	testMux.Unlock()
	defer restoreGlobals(commandExecutor, oldOsLookPath, osStat, cmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	_, _, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", os.TempDir())
	if err == nil {
		t.Error("Expected an error when ffmpeg is not found, but got none")
	}
	if !strings.Contains(err.Error(), "ffmpeg not found") {
		t.Errorf("Expected 'ffmpeg not found' error, got: %v", err)
	}
}

// TestDownloadAudio_YTDLPNotFound tests the scenario where yt-dlp is not found.
func TestDownloadAudio_YTDLPNotFound(t *testing.T) {
	testMux.Lock()
	oldOsLookPath := osLookPath
	osLookPath = func(file string) (string, error) {
		if file == "ffmpeg" {
			return "/usr/local/bin/ffmpeg", nil
		}
		if file == "yt-dlp" {
			return "", errors.New("not found in PATH")
		}
		return "", nil
	}
	testMux.Unlock()
	defer restoreGlobals(commandExecutor, oldOsLookPath, osStat, cmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	_, _, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", os.TempDir())
	if err == nil {
		t.Error("Expected an error when yt-dlp is not found, but got none")
	}
	if !strings.Contains(err.Error(), "yt-dlp not found") {
		t.Errorf("Expected 'yt-dlp not found' error, got: %v", err)
	}
}

// TestDownloadAudio_Success tests a successful download scenario.
func TestDownloadAudio_Success(t *testing.T) {
	testMux.Lock()
	oldCommandExecutor := commandExecutor
	oldOsLookPath := osLookPath
	oldOsStat := osStat
	oldCmdCombinedOutput := cmdCombinedOutput

	// Mock os.LookPath to always succeed
	osLookPath = func(file string) (string, error) {
		return "/usr/local/bin/" + file, nil
	}

	tempDir := t.TempDir()
	expectedVideoID := "test-video-id"
	expectedFilename := fmt.Sprintf("%s.wav", expectedVideoID)
	expectedFilePath := filepath.Join(tempDir, expectedFilename)
	dummyFileContent := []byte("dummy audio content")

	// Mock exec.Command to do nothing
	commandExecutor = func(name string, args ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: name,
			Args: append([]string{name}, args...),
		}
	}
	// Mock CombinedOutput to simulate success and create the file
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		if strings.Contains(cmd.String(), "--get-id") {
			return []byte(expectedVideoID), nil
		}
		// Create the dummy file to simulate yt-dlp downloading it
		err := os.WriteFile(expectedFilePath, dummyFileContent, 0644)
		if err != nil {
			return nil, fmt.Errorf("mock failed to create dummy file: %w", err)
		}
		return []byte("yt-dlp success"), nil
	}

	// Mock os.Stat to report the dummy file exists
	osStat = func(name string) (os.FileInfo, error) {
		return os.Stat(name) // Use real os.Stat
	}

	testMux.Unlock()
	defer restoreGlobals(oldCommandExecutor, oldOsLookPath, oldOsStat, oldCmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	downloadedPath, videoID, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", tempDir)

	if err != nil {
		t.Fatalf("DownloadAudio failed unexpectedly: %v", err)
	}
	if downloadedPath != expectedFilePath {
		t.Errorf("Downloaded path mismatch. Expected: %s, Got: %s", expectedFilePath, downloadedPath)
	}
	if videoID != expectedVideoID {
		t.Errorf("Video ID mismatch. Expected: %s, Got: %s", expectedVideoID, videoID)
	}

	// Verify the file was "downloaded" and exists
	_, err = os.Stat(downloadedPath)
	if os.IsNotExist(err) {
		t.Errorf("Downloaded file %s does not exist", downloadedPath)
	}
}

// TestDownloadAudio_CommandFailed tests when yt-dlp command returns an error.
func TestDownloadAudio_CommandFailed(t *testing.T) {
	testMux.Lock()
	oldCommandExecutor := commandExecutor
	oldOsLookPath := osLookPath
	oldOsStat := osStat
	oldCmdCombinedOutput := cmdCombinedOutput

	osLookPath = func(file string) (string, error) {
		if file == "ffmpeg" {
			return "/usr/local/bin/ffmpeg", nil
		}
		if file == "yt-dlp" {
			return "/usr/local/bin/yt-dlp", nil
		}
		return "", nil
	}
	osStat = func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	expectedErrorMsg := "yt-dlp error output"
	commandExecutor = func(name string, args ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: name,
			Args: append([]string{name}, args...),
		}
	}
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte("Error output from yt-dlp: " + expectedErrorMsg), errors.New("exit status 1")
	}
	testMux.Unlock()
	defer restoreGlobals(oldCommandExecutor, oldOsLookPath, oldOsStat, oldCmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	_, _, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", os.TempDir())

	if err == nil {
		t.Error("Expected an error from failed yt-dlp command, but got none")
	}
	if !strings.Contains(err.Error(), "failed") {
		t.Errorf("Expected 'failed' error with specific message, got: %v", err)
	}
}
