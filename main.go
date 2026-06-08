package main

import (
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dailymanna/manna/internal/app"
	"github.com/dailymanna/manna/internal/utils"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist migrations/*.go all:data
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	utils.Load()
	appDir := utils.GetAppConfigDir()
	err := os.MkdirAll(appDir, 0755)
	if err != nil {
		log.Fatalf("unable to create the directory %s and failing with error: %v", appDir, err)
	}

	dbPath := filepath.Join(appDir, "data", "persistence", "manna.db")
	err = os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		log.Fatalf("failed to create db directory: %v", err)
	}
	dbPathWithOptions := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-64000)", dbPath)
	db, err := sql.Open("sqlite", dbPathWithOptions)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	cfg := new(app.Config)
	cfg.FS = assets
	cfg.DB = db
	cfg.ConfigPath = appDir

	app := app.NewMannaApp(cfg)

	// Create context menu
	contextMenu := app.ContextMenu.New()
	contextMenu.Add("Add Note").OnClick(func(ctx *application.Context) {})
	contextMenu.Add("Highlight").OnClick(func(ctx *application.Context) {})
	contextMenu.Add("Bookmark this").OnClick(func(ctx *application.Context) {})

	// Register with ID
	app.ContextMenu.Add("verse-ctx-menu", contextMenu)

	// // Auto updater
	// gh, _ := github.New(github.Config{
	// 	Repository: "dailymanna/manna",
	// 	Prerelease: true,
	// })
	// if err := app.Updater.Init(updater.Config{
	// 	CurrentVersion: "0.0.0-alpha-1",
	// 	Providers:      []updater.Provider{gh},
	// }); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
	// 	log.Printf("error during update: %v", err)
	// }

	// Run the application. This blocks until the application has been exited.
	err = app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
