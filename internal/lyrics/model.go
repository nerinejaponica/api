package lyrics

type Document struct {
	SchemaVersion int      `json:"schemaVersion"`
	Stanzas       []Stanza `json:"stanzas"`
}

type Stanza struct {
	Lines        []Line        `json:"lines"`
	Translations []Translation `json:"translations"`
}

type Line struct {
	JA     string `json:"ja"`
	Romaji string `json:"romaji"`
}

type Translation struct {
	Span int    `json:"span"`
	Text string `json:"text"`
}

type Footnote struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Issue struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type Result struct {
	Errors   []Issue `json:"errors"`
	Warnings []Issue `json:"warnings"`
}

func (r Result) OK() bool {
	return len(r.Errors) == 0
}

func (r *Result) addError(path, message string) {
	r.Errors = append(r.Errors, Issue{Path: path, Message: message})
}

func (r *Result) addWarning(path, message string) {
	r.Warnings = append(r.Warnings, Issue{Path: path, Message: message})
}

func (s Stanza) texts() []string {
	var texts []string
	for _, line := range s.Lines {
		texts = append(texts, line.JA, line.Romaji)
	}
	for _, translation := range s.Translations {
		texts = append(texts, translation.Text)
	}
	return texts
}

func (d Document) texts() []string {
	var texts []string
	for _, stanza := range d.Stanzas {
		texts = append(texts, stanza.texts()...)
	}
	return texts
}
