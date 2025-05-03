package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"ra2-ini-editor/internal/logger"
)

const (
	MinWindowWidth  = 1200
	MinWindowHeight = 900
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	logFile, err := getLogFile()
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "ra2-ini-editor",
		Width:     MinWindowWidth,
		Height:    MinWindowHeight,
		MinWidth:  MinWindowWidth,
		MinHeight: MinWindowHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Logger:           &logger.Logger{},
		BackgroundColour: options.NewRGBA(255, 255, 255, 0),
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}

func getLogFilename() (string, error) {
	name := fmt.Sprintf("%d.log", time.Now().Unix())
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "ra2-ini-editor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

func getLogFile() (*os.File, error) {
	filename, err := getLogFilename()
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return f, nil
}
