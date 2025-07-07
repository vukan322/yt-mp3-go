package main

import (
	"github.com/vukan322/yt-mp3-go/internal/downloader"
	"github.com/vukan322/yt-mp3-go/internal/handler"
	"github.com/vukan322/yt-mp3-go/internal/jobs"
	"github.com/vukan322/yt-mp3-go/internal/localization"
	"github.com/vukan322/yt-mp3-go/internal/logger"
	"github.com/vukan322/yt-mp3-go/internal/view"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const basePath = "/yt-downloader"

func main() {
	logger.Setup()

	bundle, err := localization.NewBundle("./locales")
	if err != nil {
		slog.Error("failed to create i18n bundle", "error", err)
		os.Exit(1)
	}

	jobStore := jobs.NewStore()
	templates := view.ParseTemplates()

	appHandler := &handler.AppHandler{
		Downloader: &downloader.Downloader{},
		I18nBundle: bundle,
		JobStore:   jobStore,
		Templates:  templates,
		BasePath:   basePath,
	}

	go jobs.StartCleanupWorker(30*time.Minute, 2*time.Hour)

	mux := appHandler.Routes()

	slog.Info("server starting", "address", "http://localhost:8080"+basePath)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
