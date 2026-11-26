package lyrics

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"

	"github.com/nerinejaponica/api/internal/constant"
)

var footnoteMarker = regexp.MustCompile(`\[\^(\d{1,9})\]`)

func Validate(doc Document, footnotes []Footnote) Result {
	var result Result

	// CHECK SCHEMA VERSION
	result.checkSchemaVersion(doc)
	// CHECK EVERY STANZA
	for i, stanza := range doc.Stanzas {
		result.checkSpan(i, stanza)
		result.checkRuby(i, stanza)
	}
	// CHECK FOOTNOTE MARKERS AGAINST FOOTNOTES
	result.checkFootnotes(doc, footnotes)

	return result
}

func (r *Result) checkSchemaVersion(doc Document) {
	// ZERO VALUE MEANS THE FIELD IS MISSING IN THE JSON
	if doc.SchemaVersion < 1 {
		r.addError("schemaVersion", constant.MsgSchemaVersionRequired)
	}
}

func (r *Result) checkSpan(index int, stanza Stanza) {
	path := fmt.Sprintf("stanzas[%d].translations", index)

	// SUM SPAN OF EVERY TRANSLATION
	total := 0
	for i, translation := range stanza.Translations {
		total += translation.Span
		// REJECT SPAN BELOW 1 BECAUSE IT CAN MAKE THE TOTAL MATCH BY ACCIDENT
		if translation.Span < 1 {
			r.addError(fmt.Sprintf("%s[%d].span", path, i), constant.MsgSpanMin)
		}
	}

	// TOTAL SPAN MUST EQUAL THE NUMBER OF LINES IN THE STANZA
	if total != len(stanza.Lines) {
		r.addError(path, fmt.Sprintf(constant.MsgSpanMismatchFmt, total, len(stanza.Lines)))
	}
}

func (r *Result) checkRuby(index int, stanza Stanza) {
	for i, line := range stanza.Lines {
		// RUBY NOTATION ONLY LIVES IN THE JAPANESE COLUMN
		if !rubyBalanced(line.JA) {
			r.addError(fmt.Sprintf("stanzas[%d].lines[%d].ja", index, i), constant.MsgRubyUnbalanced)
		}
	}
}

func rubyBalanced(text string) bool {
	depth := 0
	for _, char := range text {
		switch char {
		case '{':
			depth++
		case '}':
			depth--
		}
		// A CLOSING BRACE WITHOUT OPENING OR A NESTED OPENING IS UNBALANCED
		if depth < 0 || depth > 1 {
			return false
		}
	}
	return depth == 0
}

func (r *Result) checkFootnotes(doc Document, footnotes []Footnote) {
	// COLLECT EVERY MARKER USED IN LYRICS AND TRANSLATIONS
	used := map[int]struct{}{}
	for _, text := range doc.texts() {
		collectMarkers(text, used)
	}

	// COLLECT EVERY DEFINED FOOTNOTE
	defined := map[int]struct{}{}
	for _, footnote := range footnotes {
		defined[footnote.ID] = struct{}{}
	}

	// MARKER WITHOUT FOOTNOTE IS AN ERROR, SORTED SO THE OUTPUT IS STABLE
	for _, id := range slices.Sorted(maps.Keys(used)) {
		if _, ok := defined[id]; !ok {
			r.addError("footnotes", fmt.Sprintf(constant.MsgFootnoteMissingFmt, id))
		}
	}

	// FOOTNOTE WITHOUT MARKER IS ONLY A WARNING
	for i, footnote := range footnotes {
		if _, ok := used[footnote.ID]; !ok {
			r.addWarning(fmt.Sprintf("footnotes[%d]", i), fmt.Sprintf(constant.MsgFootnoteUnusedFmt, footnote.ID))
		}
	}
}

func collectMarkers(text string, into map[int]struct{}) {
	for _, match := range footnoteMarker.FindAllStringSubmatch(text, -1) {
		id, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		into[id] = struct{}{}
	}
}
