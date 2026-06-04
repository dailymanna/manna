package utils

import (
	"log"
	"os"
	"path/filepath"
)

var appConfigLoc string

func init() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("unable to get the config dir: %v", err)
	}
	appDir := filepath.Join(cfgDir, "manna")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		log.Fatalf("unable to create the directory %s and failing with error: %v", appDir, err)
	}
	appConfigLoc = appDir
}

func Load() {

}

func GetAppConfigDir() string {
	return appConfigLoc
}
