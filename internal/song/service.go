package song

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nerinejaponica/api/internal/db"
)

var ErrSlugTaken = errors.New("slug sudah dipakai")

const pgUniqueViolation = "23505"

type Service struct {
	q *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) List(ctx context.Context) ([]Response, error) {
	rows, err := s.q.ListSongs(ctx)
	if err != nil {
		return nil, fmt.Errorf("list songs: %w", err)
	}
	out := make([]Response, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Response, error) {
	var publishedAt *time.Time
	if req.Status == StatusPublished {
		now := time.Now().UTC()
		publishedAt = &now
	}

	row, err := s.q.CreateSong(ctx, db.CreateSongParams{
		Slug:        req.Slug,
		TitleJa:     req.TitleJA,
		Status:      req.Status,
		PublishedAt: publishedAt,
	})
	if isUniqueViolation(err) {
		return Response{}, ErrSlugTaken
	}
	if err != nil {
		return Response{}, fmt.Errorf("create song: %w", err)
	}
	return toResponse(row), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}
