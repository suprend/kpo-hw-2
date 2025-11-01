package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	fileexport "kpo-hw-2/internal/application/files/export"
	fileimport "kpo-hw-2/internal/application/files/import"
	"kpo-hw-2/internal/infrastructure/di/bootstrap"
	infraexport "kpo-hw-2/internal/infrastructure/files/export"
	infraimport "kpo-hw-2/internal/infrastructure/files/import"
)

func main() {
	logFn, closeLog, err := openTimingLogger("logs/timings.log")
	if err != nil {
		log.Fatalf("не удалось открыть лог таймингов: %v", err)
	}
	defer func() {
		if err := closeLog(); err != nil {
			log.Printf("не удалось закрыть лог таймингов: %v", err)
		}
	}()

	app, err := bootstrap.Build(
		context.Background(),
		logFn,
		[]fileexport.Exporter{
			infraexport.NewJSONExporter(),
		},
		[]fileimport.Importer{
			infraimport.NewJSONImporter(),
		},
	)
	if err != nil {
		log.Fatalf("не удалось инициализировать приложение: %v", err)
	}

	if err := runProgram(app.Model); err != nil {
		log.Fatalf("не удалось запустить интерфейс: %v", err)
	}
}

func runProgram(model tea.Model) error {
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func openTimingLogger(path string) (func(name string, duration time.Duration, err error), func() error, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, nil, err
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(file, "", log.LstdFlags)
	logFn := func(name string, duration time.Duration, cmdErr error) {
		ms := float64(duration) / float64(time.Millisecond)
		logger.Printf("%s took %.3fms (err=%v)", name, ms, cmdErr)
	}

	closeFn := func() error {
		return file.Close()
	}

	return logFn, closeFn, nil
}
