package handler

import (
	"encoding/json"
	"fmt"
	"github.com/vukan322/yt-mp3-go/internal/downloader"
	"github.com/vukan322/yt-mp3-go/internal/jobs"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type AppHandler struct {
	I18nBundle *i18n.Bundle
	Downloader *downloader.Downloader
	JobStore   *jobs.JobStore
	BasePath   string
	Templates  *template.Template
}

func (h *AppHandler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	staticFs := http.FileServer(http.Dir("./web/static"))
	mux.Handle(h.BasePath+"/static/", http.StripPrefix(h.BasePath+"/static/", staticFs))

	mux.HandleFunc(h.BasePath+"/", h.HandleIndex)
	mux.HandleFunc(h.BasePath+"/download", h.HandleDownload)
	mux.HandleFunc(h.BasePath+"/downloads/", h.HandleServeDownload)
	mux.HandleFunc(h.BasePath+"/events", h.HandleStatusEvents)

	return mux
}

func (h *AppHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	meta, err := h.Downloader.GetMetadata(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get video metadata: %v", err), http.StatusBadRequest)
		return
	}

	job := h.JobStore.Create(url)
	go h.Downloader.Download(h.JobStore, job.ID, url)
	log.Printf("Created job %s for URL: %s", job.ID, url)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"jobID":     job.ID,
		"title":     meta.Title,
		"thumbnail": meta.Thumbnail,
	})
}

func (h *AppHandler) HandleServeDownload(w http.ResponseWriter, r *http.Request) {
	prefix := h.BasePath + "/downloads/"
	filePath := strings.TrimPrefix(r.URL.Path, prefix)

	diskPath := filepath.Join("downloads", filePath)

	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(diskPath)+"\"")
	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, diskPath)
}

func (h *AppHandler) HandleStatusEvents(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		job, found := h.JobStore.Get(jobID)
		if !found {
			return
		}
		jobJSON, _ := json.Marshal(job)
		fmt.Fprintf(w, "data: %s\n\n", jobJSON)
		flusher.Flush()
		if job.Status == jobs.StatusComplete || job.Status == jobs.StatusFailed {
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func (h *AppHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		accept := r.Header.Get("Accept-Language")
		tags, _, _ := language.ParseAcceptLanguage(accept)
		if len(tags) > 0 {
			base, _ := tags[0].Base()
			lang = base.String()
		} else {
			lang = "en"
		}
	}
	localizer := i18n.NewLocalizer(h.I18nBundle, lang)

	tmpl, err := h.Templates.Clone()
	if err != nil {
		http.Error(w, "Failed to clone templates", http.StatusInternalServerError)
		log.Println("Template cloning error:", err)
		return
	}

	tmpl.Funcs(template.FuncMap{
		"Localize": func(messageID string) template.HTML {
			msg, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
			return template.HTML(msg)
		},
	})

	data := map[string]any{
		"Lang":     lang,
		"BasePath": h.BasePath,
	}

	err = tmpl.ExecuteTemplate(w, "layout.gohtml", data)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}
