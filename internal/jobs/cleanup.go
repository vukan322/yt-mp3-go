package jobs

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func StartCleanupWorker(interval, maxAge time.Duration) {
	slog.Info("Starting background file cleanup worker", "interval", interval.String(), "maxAge", maxAge.String())

	ticker := time.NewTicker(interval)

	go func() {
		runCleanup(maxAge)

		for range ticker.C {
			runCleanup(maxAge)
		}
	}()
}

func runCleanup(maxAge time.Duration) {
	downloadsDir := "downloads"
	slog.Info("Running cleanup task for old files", "maxAge", maxAge.String())

	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		slog.Error("Cleanup: failed to read downloads directory", "path", downloadsDir, "error", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		jobDir := filepath.Join(downloadsDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			slog.Error("Cleanup: failed to get info for item", "path", jobDir, "error", err)
			continue
		}

		if time.Since(info.ModTime()) > maxAge {
			slog.Info("Cleanup: deleting expired directory", "path", jobDir)
			if err := os.RemoveAll(jobDir); err != nil {
				slog.Error("Cleanup: failed to delete directory", "path", jobDir, "error", err)
			}
		}
	}
}
