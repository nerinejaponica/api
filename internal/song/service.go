package song

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/db"
)

var ErrSlugTaken = errors.New(constant.MsgSlugTaken)

type Service struct {
	q *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) List(ctx context.Context) ([]Response, error) {
	// GET ALL SONGS FROM DATABASE
	rows, err := s.q.ListSongs(ctx)
	if err != nil {
		return nil, fmt.Errorf("list songs: %w", err)
	}
	// MAP ROWS INTO RESPONSES, NEVER RETURN NULL LIST
	out := make([]Response, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Response, error) {
	// SET PUBLISHED TIME ONLY WHEN STATUS IS PUBLISHED
	var publishedAt *time.Time
	if req.Status == StatusPublished {
		now := time.Now().UTC()
		publishedAt = &now
	}

	// INSERT SONG INTO DATABASE
	row, err := s.q.CreateSong(ctx, db.CreateSongParams{
		Slug:        req.Slug,
		TitleJa:     req.TitleJA,
		Status:      req.Status,
		PublishedAt: publishedAt,
	})
	// TREAT UNIQUE VIOLATION AS SLUG TAKEN
	if isUniqueViolation(err) {
		return Response{}, ErrSlugTaken
	}
	if err != nil {
		return Response{}, fmt.Errorf("create song: %w", err)
	}
	// CREATE ROW AND LIST ROW HAVE THE SAME COLUMNS, SO CONVERT IT
	return toResponse(db.ListSongsRow(row)), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}
