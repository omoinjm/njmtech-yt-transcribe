package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// PostgresMediaItemRepository implements MediaItemRepository against a Neon / Postgres database.
type PostgresMediaItemRepository struct {
	conn *pgx.Conn
}

// NewPostgresMediaItemRepository opens a single connection to the database at the given connString
// (i.e. the POSTGRES_URL env var) and returns a ready-to-use repository.
// The caller is responsible for calling Close() when finished.
func NewPostgresMediaItemRepository(ctx context.Context, connString string) (*PostgresMediaItemRepository, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &PostgresMediaItemRepository{conn: conn}, nil
}

// Close releases the underlying database connection.
func (r *PostgresMediaItemRepository) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

// FetchNextUnprocessed returns the oldest row in media_items where transcript_url IS NULL.
// Returns nil, nil when every item has already been processed.
func (r *PostgresMediaItemRepository) FetchNextUnprocessed(ctx context.Context) (*MediaItem, error) {
	const query = `
		SELECT id, url, platform, video_id
		FROM   media_items
		WHERE  transcript_url IS NULL
		ORDER  BY created_at ASC
		LIMIT  1`

	row := r.conn.QueryRow(ctx, query)

	var item MediaItem
	err := row.Scan(&item.ID, &item.URL, &item.Platform, &item.VideoID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next unprocessed item: %w", err)
	}
	return &item, nil
}

// UpdateTranscriptURL sets transcript_url for the row identified by id.
func (r *PostgresMediaItemRepository) UpdateTranscriptURL(ctx context.Context, id, transcriptURL string) error {
	const query = `UPDATE media_items SET transcript_url = $1 WHERE id = $2`

	tag, err := r.conn.Exec(ctx, query, transcriptURL, id)
	if err != nil {
		return fmt.Errorf("failed to update transcript_url for id %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no row found with id %s", id)
	}
	return nil
}
