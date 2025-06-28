package main

import (
	"encoding/json"
	"github.com/vukan322/yt-mp3-go/internal/downloader"
	"github.com/vukan322/yt-mp3-go/internal/handler"
	"github.com/vukan322/yt-mp3-go/internal/jobs"
	"github.com/vukan322/yt-mp3-go/internal/view"
	"log"
	"net/http"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const basePath = "/yt-downloader"

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("./locales/en.json")
	bundle.LoadMessageFile("./locales/sr.json")

	jobStore := jobs.NewStore()
	downloader := &downloader.Downloader{}
	templates := view.ParseTemplates()

	appHandler := &handler.AppHandler{
		I18nBundle: bundle,
		Downloader: downloader,
		JobStore:   jobStore,
		BasePath:   basePath,
		Templates:  templates,
	}

	go jobs.StartCleanupWorker(30*time.Minute, 2*time.Hour)

	mux := appHandler.Routes()

	log.Printf("Server starting on :8080, available at http://localhost:8080%s", basePath)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
