package bible

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"strconv"

	biblepkg "github.com/dailymanna/manna/pkg/bible"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type BibleServiceConfig struct {
	App    *application.App
	DataFS fs.FS
	DB     *sql.DB
}

func NewBibleService(cfg *BibleServiceConfig) *BibleService {
	bs := new(BibleService)
	bs.app = cfg.App
	bs.bibleBooks = bibleBooks
	bs.bibleChapters = bibleChapters
	bs.bibleVerses = bibleVerses
	bs.translationsAvailable = []string{"KJV", "AMP", "CEV", "ESV", "MSG", "NASB", "NIV", "NKJV"}
	bs.dataFS = cfg.DataFS
	bs.db = cfg.DB
	for _, t := range bs.translationsAvailable {
		path := fmt.Sprintf("data/bibles/%s.xmm", t)
		err := bs.importOpenSongBible(path, t)
		if err != nil {
			panic(err)
		}
	}
	for _, bb := range bs.bibleBooks {
		fmt.Println(bb)
	}
	return bs
}

type BibleService struct {
	app                   *application.App
	bibleBooks            []string
	bibleChapters         map[string]int
	bibleVerses           map[string][]int
	bibles                []*Bible
	notes                 []Note
	translationsAvailable []string
	dataFS                fs.FS
	db                    *sql.DB
}

type GetCrossReferencesInput struct {
	Book        string `json:"book"`
	Chapter     int    `json:"chapter"`
	VerseNumber int    `json:"verse_number"`
}

type GetCrossReferencesOutput struct {
	CrossReferences []*CrossReference `json:"cross_references"`
}

type CrossReference struct {
	ID           int    `json:"id"`
	FromBook     string `json:"from_book"`
	FromChapter  int    `json:"from_chapter"`
	FromVerse    int    `json:"from_verse"`
	ToBook       string `json:"to_book"`
	ToChapter    int    `json:"to_chapter"`
	ToVerseStart int    `json:"to_verse_start"`
	ToVerseEnd   int    `json:"to_verse_end"`
	Votes        int    `json:"votes"`
}

func (bs *BibleService) GetCrossReferences(input *GetCrossReferencesInput) (*GetCrossReferencesOutput, error) {
	tableName := "cross_references"
	fmt.Println("Book:", input.Book)
	query := fmt.Sprintf("select * from %s where from_book = ? AND from_chapter = ? AND from_verse = ? ORDER BY ID ASC;", tableName)
	rows, err := bs.db.Query(query, input.Book, input.Chapter, input.VerseNumber)
	if err != nil {
		log.Printf("Error when querying: %v", err)
		return nil, err
	}
	var refs []*CrossReference
	for rows.Next() {
		ref := new(CrossReference)
		err := rows.Scan(&ref.ID, &ref.FromBook, &ref.FromChapter, &ref.FromVerse, &ref.ToBook, &ref.ToChapter, &ref.ToVerseStart, &ref.ToVerseEnd, &ref.Votes)
		if err != nil {
			log.Printf("error when scanning the results from table %s: %v", tableName, err)
			continue
		}
		refs = append(refs, ref)
	}
	result := &GetCrossReferencesOutput{
		CrossReferences: refs,
	}
	fmt.Printf("Count of References: %d\n", len(refs))

	return result, nil
}

type Reference struct {
	ID           int    `json:"id"`
	FromBook     string `json:"from_book"`
	FromChapter  int    `json:"from_chapter"`
	FromVerse    int    `json:"from_verse"`
	ToBook       string `json:"to_book"`
	ToChapter    int    `json:"to_chapter"`
	ToVerseStart int    `json:"to_verse_start"`
	ToVerseEnd   int    `json:"to_verse_end"`
	Votes        int    `json:"votes"`
}

type GetCrossReferencesResult struct {
	References []*Reference `json:"references"`
}

func (bs *BibleService) getCrossReferences(book string, chapterNum, verseNum int) (*GetCrossReferencesResult, error) {
	tableName := "cross_references"
	fmt.Println("Book:", book)
	query := fmt.Sprintf("select * from %s where from_book = ? AND from_chapter = ? AND from_verse = ?;", tableName)
	rows, err := bs.db.Query(query, book, chapterNum, verseNum)
	if err != nil {
		log.Printf("Error when querying: %v", err)
		return nil, err
	}
	var refs []*Reference
	for rows.Next() {
		ref := new(Reference)
		err := rows.Scan(&ref.ID, &ref.FromBook, &ref.FromChapter, &ref.FromVerse, &ref.ToBook, &ref.ToChapter, &ref.ToVerseStart, &ref.ToVerseEnd, &ref.Votes)
		if err != nil {
			log.Printf("error when scanning the results from table %s: %v", tableName, err)
			continue
		}
		refs = append(refs, ref)
	}
	result := &GetCrossReferencesResult{
		References: refs,
	}
	fmt.Printf("References: %v\n", refs)
	return result, nil
}

func (bs *BibleService) GetBooksOfTheBible() []string {
	return bs.bibleBooks
}

func (bs *BibleService) GetCountOfChaptersInTheBook(book string) int {
	c := bibleChapters[book]
	return c
}

func (bs *BibleService) GetCountOfVersesInTheChapter(book string, chapter int) int {
	b := bibleVerses[book]
	v := b[chapter-1]
	return v
}

type ImportResult struct {
}

func (bs *BibleService) importOpenSongBible(path string, version string) error {
	f, err := bs.dataFS.Open(path)
	if err != nil {
		return err
	}
	osb, err := biblepkg.ParseOpenSongBible(f)
	if err != nil {
		return err
	}
	bbl := new(Bible)
	bbl.Version = version
	bbl.Code = version
	for bkIndex, book := range osb.B {
		bk := new(Book)
		bk.Name = book.N
		for _, chapter := range book.C {
			ch := new(Chapter)
			chNum, _ := strconv.Atoi(chapter.N)
			ch.Number = chNum
			for _, verse := range chapter.V {
				v := new(Verse)
				bkID := bkIndex + 1
				v.BookID = bkID
				v.Chapter = chNum
				v.Verse, _ = strconv.Atoi(verse.N)
				v.Text = verse.Text
				ch.Verses = append(ch.Verses, v)
			}
			bk.Chapters = append(bk.Chapters, ch)
		}
		bk.ChapterCount = len(bk.Chapters)
		bbl.Books = append(bbl.Books, bk)
	}
	bs.bibles = append(bs.bibles, bbl)

	return nil
}

func (bs *BibleService) updateDBWithImportedBibles() {
	tx, _ := bs.db.Begin()
	q := `
	INSERT INTO bible_translations
	(code, name, language, is_installed, installed_at) VALUES
	(?, ?, ?, ?, ?)
	`
	insertBibleTranslationStmt, _ := tx.Prepare(q)
	for _, bbl := range bs.bibles {
		insertBibleTranslationStmt.Exec(bbl.Code)
	}
}

func (bs *BibleService) ShowWarning(title, message string) {
	warningDialog := bs.app.Dialog.Warning()
	warningDialog.Title = title
	warningDialog.Message = message
	warningDialog.Show()

}

func (bs *BibleService) ShowError(title, message string) {
	errorDialog := bs.app.Dialog.Error()
	errorDialog.Title = title
	errorDialog.Message = message
	errorDialog.Show()

}

type GetVersesResponse struct {
	Num  string `json:"num"`
	Text string `json:"text"`
	// Verses []*Verse1 `json:"verses"`
}

func (bs *BibleService) isTranslationAvailable(v string) bool {
	versionAvailable := false
	for _, ver := range bs.translationsAvailable {
		if ver == v {
			versionAvailable = true
			return versionAvailable
		}
	}
	return versionAvailable
}

func (bs *BibleService) GetListOfTranslationsAvailable() []string {
	return bs.translationsAvailable
}

func (bs *BibleService) GetVerses(version string, book string, chapter int) ([]*GetVersesResponse, error) {
	if version == "" || !bs.isTranslationAvailable(version) {
		return nil, fmt.Errorf("version of bible not found")
	}
	var verses []*GetVersesResponse
	for _, bbl := range bs.bibles {
		if bbl.Version == version {
			for _, bk := range bbl.Books {
				if bk.Name == book {
					for _, chap := range bk.Chapters {
						if chap.Number == chapter {
							// verses = chap.Verses
							for i, v := range chap.Verses {
								gvr := &GetVersesResponse{
									Num:  strconv.Itoa(i + 1),
									Text: v.Text,
								}
								verses = append(verses, gvr)
							}
						}
					}
				}
			}
		}
	}
	return verses, nil
}

func (bs *BibleService) AddNote(text string, book, chapter, verse int) error {
	return nil
}

func (bs *BibleService) UpdateNote(id int, text string, book, chapter, verse int) error {
	return nil
}

func (bs *BibleService) DeleteNote(id int, text string, book, chapter, verse int) error {
	return nil
}

func (bs *BibleService) GetNoteByID(id int) ([]*Note, error) {
	return nil, nil
}

func (bs *BibleService) GetNoteByVerse(book, chapter, verse int) ([]*Note, error) {
	return nil, nil
}

func (bs *BibleService) GetNotesByBook(book int) ([]*Note, error) {
	return nil, nil
}

func (bs *BibleService) GetNotesByChapter(book, chapter int) ([]*Note, error) {
	return nil, nil
}

func (bs *BibleService) GetAllNotes() ([]*Note, error) {
	return nil, nil
}
