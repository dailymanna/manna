package bible

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dailymanna/manna/internal/utils"
	_ "modernc.org/sqlite"
)

func TestBibleService_getCrossReferences(t *testing.T) {
	utils.Load()
	appDir := utils.GetAppConfigDir()
	dbPath := filepath.Join(appDir, "data", "persistence", "manna.db")
	dbPathWithOptions := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-64000)", dbPath)
	db, err := sql.Open("sqlite", dbPathWithOptions)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	// cfg := BibleServiceConfig{
	// 	DB:     db,
	// 	App:    nil,
	// 	DataFS: nil,
	// }
	// svc := NewBibleService(&cfg)
	svc := BibleService{
		db: db,
	}
	inp := GetCrossReferencesInput{
		Book:        "Genesis",
		Chapter:     1,
		VerseNumber: 2,
	}
	out, err := svc.GetCrossReferences(&inp)
	if err != nil {
		t.Fatalf("expected err to be nil but was %v", err)
	}
	t.Logf("Out: %v", out)
	t.FailNow()
}
