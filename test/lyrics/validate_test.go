package lyrics_test

import (
	"testing"

	"github.com/nerinejaponica/api/internal/lyrics"
	"github.com/nerinejaponica/api/test/fixture"
)

func TestValidate_ValidDocumentPasses(t *testing.T) {
	got := lyrics.Validate(fixture.Doc(fixture.Stanza(2, 1, 1), fixture.Stanza(3, 3)), nil)

	fixture.AssertIssues(t, got, nil, nil)
}

func TestValidate_SchemaVersionMissing(t *testing.T) {
	doc := lyrics.Document{Stanzas: []lyrics.Stanza{fixture.Stanza(1, 1)}}

	got := lyrics.Validate(doc, nil)

	fixture.AssertIssues(t, got, []string{"schemaVersion"}, nil)
}

func TestValidate_SpanTotalEqualsLinesWithMergedTranslation(t *testing.T) {
	got := lyrics.Validate(fixture.Doc(fixture.Stanza(3, 2, 1)), nil)

	fixture.AssertIssues(t, got, nil, nil)
}

func TestValidate_SpanTotalSmallerThanLines(t *testing.T) {
	got := lyrics.Validate(fixture.Doc(fixture.Stanza(3, 1, 1)), nil)

	fixture.AssertIssues(t, got, []string{"stanzas[0].translations"}, nil)
}

func TestValidate_SpanTotalBiggerThanLines(t *testing.T) {
	got := lyrics.Validate(fixture.Doc(fixture.Stanza(1, 1), fixture.Stanza(2, 2, 1)), nil)

	fixture.AssertIssues(t, got, []string{"stanzas[1].translations"}, nil)
}

func TestValidate_SpanBelowOne(t *testing.T) {
	got := lyrics.Validate(fixture.Doc(fixture.Stanza(2, 3, -1)), nil)

	fixture.AssertIssues(t, got, []string{"stanzas[0].translations[1].span"}, nil)
}

func TestValidate_SpanMismatchMessage(t *testing.T) {
	got := lyrics.Validate(fixture.Doc(fixture.Stanza(3, 1, 1)), nil)

	want := "total span 2 tidak sama dengan jumlah baris 3"
	if len(got.Errors) != 1 || got.Errors[0].Message != want {
		t.Errorf("errors = %+v, want one error with message %q", got.Errors, want)
	}
}

func TestValidate_RubyBalancedPasses(t *testing.T) {
	doc := fixture.Doc(fixture.TextStanza("{十六夜|いざよい}の{月|つき}", "x", "uji"))

	got := lyrics.Validate(doc, nil)

	fixture.AssertIssues(t, got, nil, nil)
}

func TestValidate_RubyNotClosed(t *testing.T) {
	doc := fixture.Doc(fixture.TextStanza("{十六夜|いざよい", "x", "uji"))

	got := lyrics.Validate(doc, nil)

	fixture.AssertIssues(t, got, []string{"stanzas[0].lines[0].ja"}, nil)
}

func TestValidate_RubyClosedWithoutOpening(t *testing.T) {
	doc := fixture.Doc(fixture.TextStanza("十六夜|いざよい}", "x", "uji"))

	got := lyrics.Validate(doc, nil)

	fixture.AssertIssues(t, got, []string{"stanzas[0].lines[0].ja"}, nil)
}

func TestValidate_RubyNested(t *testing.T) {
	doc := fixture.Doc(fixture.TextStanza("{十{六|ろく}|じゅう}", "x", "uji"))

	got := lyrics.Validate(doc, nil)

	fixture.AssertIssues(t, got, []string{"stanzas[0].lines[0].ja"}, nil)
}

func TestValidate_FootnoteUsedInRomajiAndTranslationPasses(t *testing.T) {
	doc := fixture.Doc(fixture.TextStanza("あ", "toki[^1] wa", "waktu[^2]"))
	footnotes := []lyrics.Footnote{{ID: 1, Text: "uji"}, {ID: 2, Text: "uji"}}

	got := lyrics.Validate(doc, footnotes)

	fixture.AssertIssues(t, got, nil, nil)
}

func TestValidate_MarkerWithoutFootnote(t *testing.T) {
	doc := fixture.Doc(fixture.TextStanza("あ", "toki[^3]", "uji"))

	got := lyrics.Validate(doc, nil)

	fixture.AssertIssues(t, got, []string{"footnotes"}, nil)
}

func TestValidate_FootnoteWithoutMarkerIsOnlyAWarning(t *testing.T) {
	doc := fixture.Doc(fixture.Stanza(1, 1))
	footnotes := []lyrics.Footnote{{ID: 7, Text: "uji"}}

	got := lyrics.Validate(doc, footnotes)

	fixture.AssertIssues(t, got, nil, []string{"footnotes[0]"})
}

func TestValidate_FixtureValidSongPasses(t *testing.T) {
	data := fixture.LoadLyrics(t, "testdata/valid_song.json")

	got := lyrics.Validate(data.Lyrics, data.Footnotes)

	fixture.AssertIssues(t, got, nil, nil)
}

func TestValidate_FixtureSpanMismatch(t *testing.T) {
	data := fixture.LoadLyrics(t, "testdata/span_mismatch.json")

	got := lyrics.Validate(data.Lyrics, data.Footnotes)

	fixture.AssertIssues(t, got, []string{"stanzas[0].translations"}, nil)
}

func TestValidate_FixtureUnbalancedRuby(t *testing.T) {
	data := fixture.LoadLyrics(t, "testdata/unbalanced_ruby.json")

	got := lyrics.Validate(data.Lyrics, data.Footnotes)

	fixture.AssertIssues(t, got, []string{"stanzas[0].lines[0].ja"}, nil)
}
