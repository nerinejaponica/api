package lyrics

import (
	"slices"
	"testing"
)

func plainStanza(lineCount int, spans ...int) Stanza {
	stanza := Stanza{}
	for i := 0; i < lineCount; i++ {
		stanza.Lines = append(stanza.Lines, Line{JA: "あ", Romaji: "a"})
	}
	for _, span := range spans {
		stanza.Translations = append(stanza.Translations, Translation{Span: span, Text: "uji"})
	}
	return stanza
}

func docOf(stanzas ...Stanza) Document {
	return Document{SchemaVersion: 1, Stanzas: stanzas}
}

func paths(issues []Issue) []string {
	var out []string
	for _, issue := range issues {
		out = append(out, issue.Path)
	}
	return out
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name         string
		doc          Document
		footnotes    []Footnote
		wantErrors   []string
		wantWarnings []string
	}{
		{
			name: "valid document passes",
			doc:  docOf(plainStanza(2, 1, 1), plainStanza(3, 3)),
		},
		{
			name:       "schemaVersion missing",
			doc:        Document{Stanzas: []Stanza{plainStanza(1, 1)}},
			wantErrors: []string{"schemaVersion"},
		},
		{
			name: "span total equals lines passes with merged translation",
			doc:  docOf(plainStanza(3, 2, 1)),
		},
		{
			name:       "span total smaller than lines",
			doc:        docOf(plainStanza(3, 1, 1)),
			wantErrors: []string{"stanzas[0].translations"},
		},
		{
			name:       "span total bigger than lines",
			doc:        docOf(plainStanza(1, 1), plainStanza(2, 2, 1)),
			wantErrors: []string{"stanzas[1].translations"},
		},
		{
			name:       "span below one",
			doc:        docOf(plainStanza(2, 3, -1)),
			wantErrors: []string{"stanzas[0].translations[1].span"},
		},
		{
			name: "ruby balanced passes",
			doc: docOf(Stanza{
				Lines:        []Line{{JA: "{十六夜|いざよい}の{月|つき}", Romaji: "x"}},
				Translations: []Translation{{Span: 1, Text: "uji"}},
			}),
		},
		{
			name: "ruby not closed",
			doc: docOf(Stanza{
				Lines:        []Line{{JA: "{十六夜|いざよい", Romaji: "x"}},
				Translations: []Translation{{Span: 1, Text: "uji"}},
			}),
			wantErrors: []string{"stanzas[0].lines[0].ja"},
		},
		{
			name: "ruby closed without opening",
			doc: docOf(Stanza{
				Lines:        []Line{{JA: "十六夜|いざよい}", Romaji: "x"}},
				Translations: []Translation{{Span: 1, Text: "uji"}},
			}),
			wantErrors: []string{"stanzas[0].lines[0].ja"},
		},
		{
			name: "ruby nested",
			doc: docOf(Stanza{
				Lines:        []Line{{JA: "{十{六|ろく}|じゅう}", Romaji: "x"}},
				Translations: []Translation{{Span: 1, Text: "uji"}},
			}),
			wantErrors: []string{"stanzas[0].lines[0].ja"},
		},
		{
			name: "footnote used in romaji and translation passes",
			doc: docOf(Stanza{
				Lines:        []Line{{JA: "あ", Romaji: "toki[^1] wa"}},
				Translations: []Translation{{Span: 1, Text: "waktu[^2]"}},
			}),
			footnotes: []Footnote{{ID: 1, Text: "uji"}, {ID: 2, Text: "uji"}},
		},
		{
			name: "marker without footnote",
			doc: docOf(Stanza{
				Lines:        []Line{{JA: "あ", Romaji: "toki[^3]"}},
				Translations: []Translation{{Span: 1, Text: "uji"}},
			}),
			wantErrors: []string{"footnotes"},
		},
		{
			name:         "footnote without marker is only a warning",
			doc:          docOf(plainStanza(1, 1)),
			footnotes:    []Footnote{{ID: 7, Text: "uji"}},
			wantWarnings: []string{"footnotes[0]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.doc, tt.footnotes)

			if !slices.Equal(paths(got.Errors), tt.wantErrors) {
				t.Errorf("errors = %v, want %v", paths(got.Errors), tt.wantErrors)
			}
			if !slices.Equal(paths(got.Warnings), tt.wantWarnings) {
				t.Errorf("warnings = %v, want %v", paths(got.Warnings), tt.wantWarnings)
			}
			if got.OK() != (len(tt.wantErrors) == 0) {
				t.Errorf("OK() = %v, want %v", got.OK(), len(tt.wantErrors) == 0)
			}
		})
	}
}

func TestValidateSpanMismatchMessage(t *testing.T) {
	got := Validate(docOf(plainStanza(3, 1, 1)), nil)

	want := "total span 2 tidak sama dengan jumlah baris 3"
	if len(got.Errors) != 1 || got.Errors[0].Message != want {
		t.Errorf("errors = %+v, want one error with message %q", got.Errors, want)
	}
}
