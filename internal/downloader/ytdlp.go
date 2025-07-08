package downloader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/vukan322/yt-mp3-go/internal/jobs"
)

const cookiesFile = "cookies.txt"

type Metadata struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
}

type Downloader struct{}

func commandArgs(baseArgs []string) []string {
	if _, err := os.Stat(cookiesFile); err == nil {
		slog.Debug("cookies.txt found, using cookies for request")
		return append([]string{"--cookies", cookiesFile}, baseArgs...)
	}
	slog.Debug("cookies.txt not found, proceeding without cookies")
	return baseArgs
}

func (d *Downloader) GetMetadata(url string) (*Metadata, error) {
	args := commandArgs([]string{"--no-playlist", "--dump-single-json", url})
	cmd := exec.Command("yt-dlp", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("metadata command failed: %s", string(output))
	}

	jsonStartIndex := bytes.Index(output, []byte("{"))
	if jsonStartIndex == -1 {
		return nil, fmt.Errorf("yt-dlp returned no JSON output: %s", string(output))
	}

	jsonOutput := output[jsonStartIndex:]

	var meta Metadata
	if err := json.Unmarshal(jsonOutput, &meta); err != nil {
		return nil, fmt.Errorf("yt-dlp returned a non-JSON response: %s", string(output))
	}
	return &meta, nil
}

func (d *Downloader) Download(store *jobs.JobStore, jobID, videoID string, quality AudioQuality) {
	slog.Info("starting download", "jobID", jobID, "videoID", videoID, "quality", quality)
	store.UpdateStatus(jobID, jobs.StatusProcessing)

	outputDir := filepath.Join("downloads", jobID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		errMsg := fmt.Sprintf("could not create output dir: %v", err)
		slog.Error("download error", "jobID", jobID, "error", errMsg)
		store.SetResult(jobID, "", errMsg)
		return
	}

	audioQuality := quality.ToYtDlp()

	baseArgs := []string{
		"--no-playlist", "--extract-audio", "--audio-format", "mp3",
		"--audio-quality", audioQuality, "-o", "%(title)s.%(ext)s",
		"-P", outputDir,
		videoID,
	}
	cmd := exec.Command("yt-dlp", commandArgs(baseArgs)...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		errMsg := fmt.Sprintf("failed to start command: %v", err)
		slog.Error("download error", "jobID", jobID, "error", errMsg)
		store.SetResult(jobID, "", errMsg)
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
		slog.Error("download error", "jobID", jobID, "error", errMsg)
		store.SetResult(jobID, "", errMsg)
		return
	}

	slog.Info("command finished, searching for MP3 file", "jobID", jobID)
	files, err := os.ReadDir(outputDir)
	if err != nil {
		errMsg := fmt.Sprintf("could not read output dir: %v", err)
		slog.Error("download error", "jobID", jobID, "error", errMsg)
		store.SetResult(jobID, "", errMsg)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp3") {
			mp3Path := filepath.Join(outputDir, file.Name())

			fileInfo, err := os.Stat(mp3Path)
			if err != nil {
				errMsg := fmt.Sprintf("could not get file info: %v", err)
				slog.Error("download error", "jobID", jobID, "path", mp3Path, "error", errMsg)
				store.SetResult(jobID, "", errMsg)
				return
			}

			fileSize := fileInfo.Size()
			slog.Info("found MP3 file", "jobID", jobID, "path", mp3Path, "size", fileSize)
			store.SetResult(jobID, mp3Path, "")
			return
		}
	}

	errMsg := "no mp3 file found after download"
	slog.Error("download error", "jobID", jobID, "error", errMsg)
	store.SetResult(jobID, "", errMsg)
}

func logPipe(pipe io.ReadCloser, jobID string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[download]") {
			continue
		}
		slog.Debug("yt-dlp output", "jobID", jobID, "output", line)
	}
}
