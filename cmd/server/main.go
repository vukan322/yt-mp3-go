package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vukan322/yt-mp3-go/internal/config"
	"github.com/vukan322/yt-mp3-go/internal/downloader"
	"github.com/vukan322/yt-mp3-go/internal/handler"
	"github.com/vukan322/yt-mp3-go/internal/jobs"
	"github.com/vukan322/yt-mp3-go/internal/localization"
	"github.com/vukan322/yt-mp3-go/internal/logger"
	"github.com/vukan322/yt-mp3-go/internal/view"
)

var version = "development"

func main() {
	_ = godotenv.Load()

	conf := config.New()
	logger.Setup(conf.Environment)

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
		BasePath:   conf.BasePath,
		Version:    version,
	}

	go jobs.StartCleanupWorker(30*time.Minute, 2*time.Hour)

	mux := appHandler.Routes()

	startServer(mux, conf)
}

func startServer(mux *http.ServeMux, conf *config.Config) {
	var serverURL string

	if conf.Environment == "production" {
		serverURL = fmt.Sprintf("https://%s%s", conf.Domain, conf.BasePath)
	} else {
		serverURL = fmt.Sprintf("http://%s:%s%s", conf.Domain, conf.Port, conf.BasePath)
	}

	slog.Info("server starting", "address", serverURL, "env", conf.Environment)

	listenAddr := fmt.Sprintf(":%s", conf.Port)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
