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

// TestDownloadAudio_YTDLPNotFound tests the scenario where yt-dlp is not found.
func TestDownloadAudio_YTDLPNotFound(t *testing.T) {
	testMux.Lock()
	oldOsLookPath := osLookPath
	oldCmdCombinedOutput := cmdCombinedOutput
	oldOsStat := osStat

	osLookPath = func(file string) (string, error) {
		return "", errors.New("not found in PATH")
	}
	osStat = func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist } // Default mock for osStat
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) { // Default mock for CombinedOutput
		return nil, errors.New("command not found")
	}
	testMux.Unlock()
	defer restoreGlobals(commandExecutor, oldOsLookPath, oldOsStat, oldCmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	_, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", os.TempDir())
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
		return "/usr/local/bin/yt-dlp", nil
	}

	tempDir := t.TempDir()
	expectedFilePath := filepath.Join(tempDir, "Test_Video_Title.mp3")
	dummyFileContent := []byte("dummy audio content")

	// Create the dummy file *before* mocking osStat to ensure os.Stat finds it
	err := os.WriteFile(expectedFilePath, dummyFileContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy audio file: %v", err)
	}

	// Mock os.Stat to report the dummy file exists
	osStat = func(name string) (os.FileInfo, error) {
		if name == expectedFilePath {
			return os.Stat(name) // Use real os.Stat for the dummy file we created
		}
		return nil, os.ErrNotExist
	}

	commandExecutor = func(name string, args ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: name,
			Args: append([]string{name}, args...),
		}
	}
	// Mock CombinedOutput for the dummy exec.Cmd
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte(fmt.Sprintf("Some yt-dlp output\n[ExtractAudio] Destination: %s\n", expectedFilePath)), nil
	}
	testMux.Unlock()
	defer restoreGlobals(oldCommandExecutor, oldOsLookPath, oldOsStat, oldCmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	downloadedPath, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", tempDir)

	if err != nil {
		t.Fatalf("DownloadAudio failed unexpectedly: %v", err)
	}
	if downloadedPath != expectedFilePath {
		t.Errorf("Downloaded path mismatch. Expected: %s, Got: %s", expectedFilePath, downloadedPath)
	}

	// Verify the file was "downloaded" and exists (via mocked osStat)
	_, err = osStat(downloadedPath)
	if os.IsNotExist(err) {
		t.Errorf("Downloaded file %s does not exist according to mocked osStat", downloadedPath)
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
		return "/usr/local/bin/yt-dlp", nil
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
	_, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", os.TempDir())

	if err == nil {
		t.Error("Expected an error from failed yt-dlp command, but got none")
	}
	if !strings.Contains(err.Error(), "yt-dlp command failed") || !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected 'yt-dlp command failed' error with specific message, got: %v", err)
	}
}

// TestDownloadAudio_NoDestinationInOutput tests when yt-dlp doesn't print destination.
func TestDownloadAudio_NoDestinationInOutput(t *testing.T) {
	testMux.Lock()
	oldCommandExecutor := commandExecutor
	oldOsLookPath := osLookPath
	oldOsStat := osStat
	oldCmdCombinedOutput := cmdCombinedOutput

	osLookPath = func(file string) (string, error) {
		return "/usr/local/bin/yt-dlp", nil
	}
	osStat = func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	commandExecutor = func(name string, args ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: name,
			Args: append([]string{name}, args...),
		}
	}
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte("Some other yt-dlp output without destination line."), nil
	}
	testMux.Unlock()
	defer restoreGlobals(oldCommandExecutor, oldOsLookPath, oldOsStat, oldCmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	_, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", os.TempDir())

	if err == nil {
		t.Error("Expected an error for no destination in output, but got none")
	}
	if !strings.Contains(err.Error(), "could not determine downloaded file path") {
		t.Errorf("Expected 'could not determine downloaded file path' error, got: %v", err)
	}
}

// TestDownloadAudio_FileNotExistAfterSuccess tests when yt-dlp reports success but file is missing.
func TestDownloadAudio_FileNotExistAfterSuccess(t *testing.T) {
	testMux.Lock()
	oldCommandExecutor := commandExecutor
	oldOsLookPath := osLookPath
	oldOsStat := osStat
	oldCmdCombinedOutput := cmdCombinedOutput

	osLookPath = func(file string) (string, error) {
		return "/usr/local/bin/yt-dlp", nil
	}

	tempDir := t.TempDir()
	expectedFilePath := filepath.Join(tempDir, "NonExistent_File.mp3") // This file will *not* be created

	// Mock os.Stat to report the file does NOT exist
	osStat = func(name string) (os.FileInfo, error) {
		if name == expectedFilePath {
			return nil, os.ErrNotExist
		}
		return os.Stat(name) // Use real os.Stat for other files if needed
	}

	commandExecutor = func(name string, args ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: name,
			Args: append([]string{name}, args...),
		}
	}
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte(fmt.Sprintf("[ExtractAudio] Destination: %s\n", expectedFilePath)), nil
	}
	testMux.Unlock()
	defer restoreGlobals(oldCommandExecutor, oldOsLookPath, oldOsStat, oldCmdCombinedOutput)

	downloader := NewYTDLPAudioDownloader()
	_, err := downloader.DownloadAudio("https://youtube.com/watch?v=test", tempDir)

	if err == nil {
		t.Error("Expected an error for file not existing after reported success, but got none")
	}
	if !strings.Contains(err.Error(), "file not found at expected path") {
		t.Errorf("Expected 'file not found' error, got: %v", err)
	}
}