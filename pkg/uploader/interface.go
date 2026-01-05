package uploader

// Uploader defines the interface for uploading content.
type Uploader interface {
	// Upload takes content as a string and uploads it, returning the response from the upload service.
	Upload(content string, filename string) (string, error)
}
