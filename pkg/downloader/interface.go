package downloader

// YouTubeDownloader defines the interface for downloading audio from YouTube videos.
// Applying the Interface Segregation Principle (ISP) and Dependency Inversion Principle (DIP).
type YouTubeDownloader interface {
	DownloadAudio(videoURL string, outputDir string) (string, error)
}
