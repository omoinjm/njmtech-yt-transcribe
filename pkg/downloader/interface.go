package downloader

// VideoDownloader defines the interface for downloading audio from videos.
// Applying the Interface Segregation Principle (ISP) and Dependency Inversion Principle (DIP).
type VideoDownloader interface {
	DownloadAudio(videoURL string, outputDir string) (string, error)
}
