package fixture

import (
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/nerinejaponica/api/internal/lyrics"
)

// LyricsFile MATCHES THE JSON SHAPE OF A FIXTURE FILE
type LyricsFile struct {
	Lyrics    lyrics.Document   `json:"lyrics"`
	Footnotes []lyrics.Footnote `json:"footnotes"`
}

// Stanza BUILDS A STANZA WITH N PLAIN LINES AND ONE TRANSLATION PER SPAN
func Stanza(lineCount int, spans ...int) lyrics.Stanza {
	stanza := lyrics.Stanza{}
	for i := 0; i < lineCount; i++ {
		stanza.Lines = append(stanza.Lines, lyrics.Line{JA: "あ", Romaji: "a"})
	}
	for _, span := range spans {
		stanza.Translations = append(stanza.Translations, lyrics.Translation{Span: span, Text: "uji"})
	}
	return stanza
}

// Doc WRAPS STANZAS IN A DOCUMENT WITH A VALID SCHEMA VERSION
func Doc(stanzas ...lyrics.Stanza) lyrics.Document {
	return lyrics.Document{SchemaVersion: 1, Stanzas: stanzas}
}

// LoadLyrics READS A JSON FIXTURE, THE PATH IS RELATIVE TO THE TEST PACKAGE
func LoadLyrics(t *testing.T, path string) LyricsFile {
	t.Helper()

	// READ FILE
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	// DECODE JSON
	var file LyricsFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatalf("decode fixture %s: %v", path, err)
	}
	return file
}

// TextStanza BUILDS A ONE LINE STANZA WITH CUSTOM TEXT, FOR RUBY AND FOOTNOTE CASES
func TextStanza(ja, romaji, translation string) lyrics.Stanza {
	return lyrics.Stanza{
		Lines:        []lyrics.Line{{JA: ja, Romaji: romaji}},
		Translations: []lyrics.Translation{{Span: 1, Text: translation}},
	}
}

// AssertIssues CHECKS ERROR PATHS, WARNING PATHS, AND OK() OF A VALIDATION RESULT
func AssertIssues(t *testing.T, got lyrics.Result, wantErrors, wantWarnings []string) {
	t.Helper()

	// COMPARE ERROR PATHS
	if gotErrors := issuePaths(got.Errors); !slices.Equal(gotErrors, wantErrors) {
		t.Errorf("errors = %v, want %v", gotErrors, wantErrors)
	}
	// COMPARE WARNING PATHS
	if gotWarnings := issuePaths(got.Warnings); !slices.Equal(gotWarnings, wantWarnings) {
		t.Errorf("warnings = %v, want %v", gotWarnings, wantWarnings)
	}
	// OK() MUST FOLLOW THE ERRORS ONLY
	if wantOK := len(wantErrors) == 0; got.OK() != wantOK {
		t.Errorf("OK() = %v, want %v", got.OK(), wantOK)
	}
}

func issuePaths(issues []lyrics.Issue) []string {
	var out []string
	for _, issue := range issues {
		out = append(out, issue.Path)
	}
	return out
}
