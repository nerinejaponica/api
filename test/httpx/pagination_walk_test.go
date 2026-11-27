package httpx_test

import (
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/nerinejaponica/api/internal/httpx"
)

// before TELLS IF a COMES BEFORE b IN THE LIST ORDER: published_at DESC NULLS LAST, id DESC
func before(a, b item) bool {
	switch {
	case a.PublishedAt != nil && b.PublishedAt == nil:
		return true
	case a.PublishedAt == nil && b.PublishedAt != nil:
		return false
	case a.PublishedAt != nil && !a.PublishedAt.Equal(*b.PublishedAt):
		return a.PublishedAt.After(*b.PublishedAt)
	}
	return a.ID > b.ID
}

// testDataset HAS TIED published_at (IDS 7 AND 5) AND TWO DRAFTS WITHOUT published_at (IDS 4 AND 2).
// IDS ARE NOT IN THE SAME ORDER AS DATES ON PURPOSE. EXPECTED ORDER: 3, 7, 5, 6, 1, 4, 2.
func testDataset() []item {
	day := func(d int) *time.Time {
		t := time.Date(2026, 11, d, 0, 0, 0, 0, time.UTC)
		return &t
	}
	return []item{
		{ID: 1, PublishedAt: day(7)},
		{ID: 2},
		{ID: 3, PublishedAt: day(10)},
		{ID: 4},
		{ID: 5, PublishedAt: day(9)},
		{ID: 6, PublishedAt: day(8)},
		{ID: 7, PublishedAt: day(9)},
	}
}

// fetch ACTS LIKE THE FUTURE QUERY: SORT, SKIP EVERYTHING UP TO THE CURSOR, TAKE limit+1 ROWS
func fetch(all []item, cursor *httpx.Cursor, limit int) []item {
	sorted := slices.Clone(all)
	sort.Slice(sorted, func(i, j int) bool { return before(sorted[i], sorted[j]) })

	// KEEP ONLY ROWS THAT COME AFTER THE CURSOR
	rows := sorted
	if cursor != nil {
		last := item{ID: cursor.ID, PublishedAt: cursor.PublishedAt}
		rows = nil
		for _, row := range sorted {
			if before(last, row) {
				rows = append(rows, row)
			}
		}
	}
	// TAKE ONE ROW MORE THAN LIMIT, THE ROW THAT TELLS IF THERE IS A NEXT PAGE
	if len(rows) > limit+1 {
		rows = rows[:limit+1]
	}
	return rows
}

// walk FOLLOWS next_cursor FROM THE FIRST PAGE UNTIL IT IS NULL AND RETURNS THE IDS OF EVERY PAGE
func walk(t *testing.T, all []item, limit int) [][]int64 {
	t.Helper()
	var pages [][]int64
	var cursor *httpx.Cursor

	for range 100 {
		resp := httpx.NewListResponse(fetch(all, cursor, limit), limit, cursorOfItem)

		var ids []int64
		for _, row := range resp.Data {
			ids = append(ids, row.ID)
		}
		pages = append(pages, ids)

		if resp.NextCursor == nil {
			return pages
		}
		decoded, err := httpx.DecodeCursor(*resp.NextCursor)
		if err != nil {
			t.Fatalf("next_cursor %q cannot be decoded: %v", *resp.NextCursor, err)
		}
		cursor = &decoded
	}
	t.Fatal("still paging after 100 pages, next_cursor never became null")
	return nil
}

func TestPaging_LimitTwoSplitsSevenRowsIntoFourPagesInOrder(t *testing.T) {
	got := walk(t, testDataset(), 2)

	want := [][]int64{{3, 7}, {5, 6}, {1, 4}, {2}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("pages = %v, want %v", got, want)
	}
}

func TestPaging_TiedPublishedAtAreNeitherSkippedNorRepeatedAtPageBoundary(t *testing.T) {
	// LIMIT 2 ENDS THE FIRST PAGE AT ID 7 AND STARTS THE SECOND AT ID 5, WHICH SHARE published_at
	got := walk(t, testDataset(), 2)

	if got[0][1] != 7 || got[1][0] != 5 {
		t.Errorf("pages = %v, want the tied rows 7 and 5 split across pages 1 and 2", got)
	}
}

func TestPaging_EveryRowAppearsExactlyOnceForEveryLimit(t *testing.T) {
	for limit := 1; limit <= 8; limit++ {
		var all []int64
		for _, page := range walk(t, testDataset(), limit) {
			all = append(all, page...)
		}

		want := []int64{3, 7, 5, 6, 1, 4, 2}
		if !slices.Equal(all, want) {
			t.Errorf("limit %d: rows = %v, want %v", limit, all, want)
		}
	}
}

func TestPaging_DraftsComeLastAndCanBeReachedWithCursor(t *testing.T) {
	got := walk(t, testDataset(), 1)

	last := got[len(got)-2:]
	if last[0][0] != 4 || last[1][0] != 2 {
		t.Errorf("last pages = %v, want the drafts 4 and 2 at the very end", last)
	}
}

func TestPaging_LimitLargerThanDataGivesOnePageWithoutCursor(t *testing.T) {
	got := walk(t, testDataset(), 50)

	if len(got) != 1 || len(got[0]) != 7 {
		t.Errorf("pages = %v, want one page with all 7 rows", got)
	}
}

func TestPaging_LimitEqualToRowCountGivesOnePageWithoutCursor(t *testing.T) {
	got := walk(t, testDataset(), 7)

	if len(got) != 1 {
		t.Errorf("pages = %v, want exactly one page", got)
	}
}

func TestPaging_NoRowsGivesOneEmptyPage(t *testing.T) {
	got := walk(t, nil, 20)

	if len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("pages = %v, want one empty page", got)
	}
}
