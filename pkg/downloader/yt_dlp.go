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

// cmdCombinedOutput is a variable that can be overridden for testing purposes.
var cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

// YTDLPAudioDownloader is an implementation of YouTubeDownloader that uses the `yt-dlp` external tool.
// It downloads the audio stream of a given YouTube video.
