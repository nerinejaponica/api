package song_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/mock/gomock"

	"github.com/nerinejaponica/api/internal/db"
	"github.com/nerinejaponica/api/internal/song"
	"github.com/nerinejaponica/api/test/mock"
)

// newService BUILDS A SERVICE ON TOP OF A MOCKED DATABASE
func newService(t *testing.T) (*song.Service, *mock.MockQuerier) {
	t.Helper()
	querier := mock.NewMockQuerier(gomock.NewController(t))
	return song.NewService(querier), querier
}

func TestList_MapsRowsToResponses(t *testing.T) {
	svc, querier := newService(t)
	publishedAt := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	querier.EXPECT().ListSongs(gomock.Any()).Return([]db.ListSongsRow{
		{ID: 1, Slug: "izayoi", TitleJa: "十六夜", Status: "published", PublishedAt: &publishedAt},
	}, nil)

	got, err := svc.List(context.Background())

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(got) != 1 || got[0].ID != 1 || got[0].Slug != "izayoi" || got[0].TitleJA != "十六夜" {
		t.Errorf("got %+v, want one song izayoi", got)
	}
	if got[0].PublishedAt == nil || !got[0].PublishedAt.Equal(publishedAt) {
		t.Errorf("published_at = %v, want %v", got[0].PublishedAt, publishedAt)
	}
}

func TestList_NoRowsReturnsEmptyListNotNil(t *testing.T) {
	svc, querier := newService(t)
	querier.EXPECT().ListSongs(gomock.Any()).Return(nil, nil)

	got, err := svc.List(context.Background())

	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("got %#v, err %v, want empty non-nil list", got, err)
	}
}

func TestList_DatabaseErrorIsWrapped(t *testing.T) {
	svc, querier := newService(t)
	dbErr := errors.New("connection lost")
	querier.EXPECT().ListSongs(gomock.Any()).Return(nil, dbErr)

	_, err := svc.List(context.Background())

	if !errors.Is(err, dbErr) {
		t.Errorf("err = %v, want it to wrap %v", err, dbErr)
	}
}

func TestCreate_DraftHasNoPublishedAt(t *testing.T) {
	svc, querier := newService(t)
	var sent db.CreateSongParams
	querier.EXPECT().CreateSong(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, arg db.CreateSongParams) (db.CreateSongRow, error) {
			sent = arg
			return db.CreateSongRow{ID: 1, Slug: arg.Slug, TitleJa: arg.TitleJa, Status: arg.Status}, nil
		})

	got, err := svc.Create(context.Background(), song.CreateRequest{Slug: "izayoi", TitleJA: "十六夜", Status: song.StatusDraft})

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if sent.PublishedAt != nil || got.PublishedAt != nil {
		t.Errorf("published_at sent %v, returned %v, want nil for a draft", sent.PublishedAt, got.PublishedAt)
	}
}

func TestCreate_PublishedSetsPublishedAtToNow(t *testing.T) {
	svc, querier := newService(t)
	var sent db.CreateSongParams
	querier.EXPECT().CreateSong(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, arg db.CreateSongParams) (db.CreateSongRow, error) {
			sent = arg
			return db.CreateSongRow{ID: 1, Slug: arg.Slug, Status: arg.Status, PublishedAt: arg.PublishedAt}, nil
		})
	before := time.Now().UTC().Add(-time.Second)

	_, err := svc.Create(context.Background(), song.CreateRequest{Slug: "izayoi", TitleJA: "十六夜", Status: song.StatusPublished})

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if sent.PublishedAt == nil || sent.PublishedAt.Before(before) {
		t.Errorf("published_at = %v, want a time close to now", sent.PublishedAt)
	}
}

func TestCreate_UniqueViolationBecomesSlugTaken(t *testing.T) {
	svc, querier := newService(t)
	querier.EXPECT().CreateSong(gomock.Any(), gomock.Any()).Return(db.CreateSongRow{}, &pgconn.PgError{Code: "23505"})

	_, err := svc.Create(context.Background(), song.CreateRequest{Slug: "izayoi", TitleJA: "十六夜", Status: song.StatusDraft})

	if !errors.Is(err, song.ErrSlugTaken) {
		t.Errorf("err = %v, want ErrSlugTaken", err)
	}
}

func TestCreate_OtherDatabaseErrorIsWrappedNotSlugTaken(t *testing.T) {
	svc, querier := newService(t)
	dbErr := errors.New("connection lost")
	querier.EXPECT().CreateSong(gomock.Any(), gomock.Any()).Return(db.CreateSongRow{}, dbErr)

	_, err := svc.Create(context.Background(), song.CreateRequest{Slug: "izayoi", TitleJA: "十六夜", Status: song.StatusDraft})

	if !errors.Is(err, dbErr) || errors.Is(err, song.ErrSlugTaken) {
		t.Errorf("err = %v, want it to wrap %v and not be ErrSlugTaken", err, dbErr)
	}
}
