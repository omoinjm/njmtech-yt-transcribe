package repository

import (
	"context"
	"errors"
	"testing"
)

// mockRepo is a test double for MediaItemRepository.
type mockRepo struct {
	fetchResult *MediaItem
	fetchErr    error
	updateErr   error

	lastUpdateID  string
	lastUpdateURL string
}

func (m *mockRepo) FetchNextUnprocessed(_ context.Context) (*MediaItem, error) {
	return m.fetchResult, m.fetchErr
}

func (m *mockRepo) UpdateTranscriptURL(_ context.Context, id, transcriptURL string) error {
	m.lastUpdateID = id
	m.lastUpdateURL = transcriptURL
	return m.updateErr
}

// Verify mockRepo satisfies the interface at compile time.
var _ MediaItemRepository = (*mockRepo)(nil)

func TestFetchNextUnprocessed_ReturnsItem(t *testing.T) {
	expected := &MediaItem{
		ID:       "abc-123",
		URL:      "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		Platform: "youtube",
		VideoID:  "dQw4w9WgXcQ",
	}
	repo := &mockRepo{fetchResult: expected}

	item, err := repo.FetchNextUnprocessed(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item == nil {
		t.Fatal("expected a MediaItem, got nil")
	}
	if item.ID != expected.ID {
		t.Errorf("ID: want %q, got %q", expected.ID, item.ID)
	}
	if item.URL != expected.URL {
		t.Errorf("URL: want %q, got %q", expected.URL, item.URL)
	}
	if item.Platform != expected.Platform {
		t.Errorf("Platform: want %q, got %q", expected.Platform, item.Platform)
	}
	if item.VideoID != expected.VideoID {
		t.Errorf("VideoID: want %q, got %q", expected.VideoID, item.VideoID)
	}
}

func TestFetchNextUnprocessed_ReturnsNilWhenEmpty(t *testing.T) {
	repo := &mockRepo{fetchResult: nil, fetchErr: nil}

	item, err := repo.FetchNextUnprocessed(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item != nil {
		t.Errorf("expected nil item when queue is empty, got %+v", item)
	}
}

func TestFetchNextUnprocessed_PropagatesError(t *testing.T) {
	expectedErr := errors.New("connection refused")
	repo := &mockRepo{fetchErr: expectedErr}

	_, err := repo.FetchNextUnprocessed(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("want %v, got %v", expectedErr, err)
	}
}

func TestUpdateTranscriptURL_StoresValues(t *testing.T) {
	repo := &mockRepo{}
	id := "abc-123"
	blobURL := "https://blob.vercel-storage.com/yt-transcribe/youtube/dQw4w9WgXcQ"

	if err := repo.UpdateTranscriptURL(context.Background(), id, blobURL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastUpdateID != id {
		t.Errorf("ID: want %q, got %q", id, repo.lastUpdateID)
	}
	if repo.lastUpdateURL != blobURL {
		t.Errorf("URL: want %q, got %q", blobURL, repo.lastUpdateURL)
	}
}

func TestUpdateTranscriptURL_PropagatesError(t *testing.T) {
	expectedErr := errors.New("update failed")
	repo := &mockRepo{updateErr: expectedErr}

	err := repo.UpdateTranscriptURL(context.Background(), "id", "url")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("want %v, got %v", expectedErr, err)
	}
}
