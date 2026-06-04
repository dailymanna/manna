package settings

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dailymanna/manna/internal/utils"
	_ "modernc.org/sqlite"
)

func TestSetupOfflineCrossReferences(t *testing.T) {
	utils.Load()
	appDir := utils.GetAppConfigDir()
	dbPath := filepath.Join(appDir, "data", "persistence", "manna.db")
	dbPathWithOptions := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-64000)", dbPath)
	db, err := sql.Open("sqlite", dbPathWithOptions)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	cfg := SettingsServiceConfig{
		DB: db,
	}
	ss := NewSettingsService(&cfg)
	err = ss.SetupOfflineCrossReferences()
	if err != nil {
		t.Fatalf("expected err to be nil but was %v", err)
	}

}
