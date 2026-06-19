package bible

import "time"

type Bible struct {
	Books       []*Book
	Code        string // 'NIV', 'ESV', 'KJV', 'NLT'
	Version     string
	Language    string
	IsInstalled bool
	InstalledAt *time.Time
}

// Book represents a book of the Bible with its metadata.
type Book struct {
	ID           int          `json:"id"`         // 1 for Genesis, 66 for Revelation
	Name         string       `json:"name"`       // e.g., "Genesis"
	ShortName    string       `json:"short_name"` // e.g., "Gen" or "Ge"
	Testament    string       `json:"testament"`  // "Old" or "New"
	Genre        string       `json:"genre"`      // e.g., "Pentateuch", "Gospels", "Epistle"
	Metadata     BookMetadata `json:"metadata"`
	ChapterCount int          `json:"chapter_count"`
	Chapters     []*Chapter
}

type BookMetadata struct {
	Authorship string `json:"authorship"` // e.g., "Moses"
	Dating     string `json:"dating"`     // e.g., "c. 1445–1405 BC"
	Summary    string `json:"summary"`    // Short introduction for study
}

type Chapter struct {
	Number int
	Verses []*Verse
}

// VerseKey acts as a unique, lightweight identifier for any verse in the Bible.
// This is incredibly useful for map lookups and frontend routing.
type VerseKey struct {
	BookID  int `json:"book_id"`
	Chapter int `json:"chapter"`
	Verse   int `json:"verse"`
}

// Verse represents the actual text and its direct properties.
type Verse struct {
	VerseKey
	Text       string           `json:"text"`
	References []CrossReference `json:"references,omitempty"`
	Notes      []Note           `json:"notes,omitempty"`
}

// CrossReference connects a source verse to a target verse.

// Note represents user-generated study content or commentary.
type Note struct {
	ID        int        `json:"id"`
	VerseRef  VerseKey   `json:"verse_ref"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`    // Supports Markdown or HTML string from Vue
	CreatedAt *time.Time `json:"created_at"` // Unix timestamp for easy frontend formatting
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
