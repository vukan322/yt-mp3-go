package downloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/vukan322/yt-mp3-go/internal/jobs"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Metadata struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
}

type Downloader struct{}

func (d *Downloader) GetMetadata(url string) (*Metadata, error) {
	cmd := exec.Command("yt-dlp", "--cookies", "cookies.txt", "--no-playlist", "--dump-single-json", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("metadata command failed: %s", string(output))
	}
	var meta Metadata
	if err := json.Unmarshal(output, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata json: %w", err)
	}
	return &meta, nil
}

func (d *Downloader) Download(store *jobs.JobStore, jobID, url string) {
	log.Printf("[JOB %s] Starting download...", jobID)
	store.UpdateStatus(jobID, jobs.StatusProcessing)

	outputDir := filepath.Join("downloads", jobID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		errMsg := fmt.Sprintf("could not create output dir: %v", err)
		log.Printf("[JOB %s] ERROR: %s", jobID, errMsg)
		store.SetResult(jobID, "", 0, errMsg)
		return
	}

	cmd := exec.Command("yt-dlp",
		"--cookies", "cookies.txt",
		"--no-playlist", "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0",
		"-o", "%(title)s.%(ext)s",
		"-P", outputDir,
		url,
	)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		errMsg := fmt.Sprintf("failed to start command: %v", err)
		log.Printf("[JOB %s] ERROR: %s", jobID, errMsg)
		store.SetResult(jobID, "", 0, errMsg)
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go logPipe(stdout, jobID, &wg)
	go logPipe(stderr, jobID, &wg)
	err := cmd.Wait()
	wg.Wait()

	if err != nil {
		errMsg := "yt-dlp command finished with an error."
		log.Printf("[JOB %s] ERROR: %s", jobID, errMsg)
		store.SetResult(jobID, "", 0, errMsg)
		return
	}

	log.Printf("[JOB %s] Command finished. Searching for MP3 file.", jobID)
	files, err := os.ReadDir(outputDir)
	if err != nil {
		errMsg := fmt.Sprintf("could not read output dir: %v", err)
		log.Printf("[JOB %s] ERROR: %s", jobID, errMsg)
		store.SetResult(jobID, "", 0, errMsg)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp3") {
			mp3Path := filepath.Join(outputDir, file.Name())

			fileInfo, err := os.Stat(mp3Path)
			if err != nil {
				errMsg := fmt.Sprintf("could not get file info for %s: %v", mp3Path, err)
				log.Printf("[JOB %s] ERROR: %s", jobID, errMsg)
				store.SetResult(jobID, "", 0, errMsg)
				return
			}

			fileSize := fileInfo.Size()
			log.Printf("[JOB %s] Found MP3 file: %s (Size: %d bytes)", jobID, mp3Path, fileSize)
			store.SetResult(jobID, mp3Path, fileSize, "")
			return
		}
	}

	errMsg := "no mp3 file found after download"
	log.Printf("[JOB %s] ERROR: %s", jobID, errMsg)
	store.SetResult(jobID, "", 0, errMsg)
}

func logPipe(pipe io.ReadCloser, jobID string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.Printf("[JOB %s | yt-dlp] %s", jobID, scanner.Text())
	}
}
