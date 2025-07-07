package jobs

import (
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusComplete   JobStatus = "complete"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID       string    `json:"jobID"`
	URL      string    `json:"-"`
	Status   JobStatus `json:"status"`
	FilePath string    `json:"filePath"`
	FileSize int64     `json:"fileSize"`
	Error    string    `json:"error"`
}

type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*Job),
	}
}

func (s *JobStore) Create(url string) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	job := &Job{
		ID:     uuid.New().String(),
		URL:    url,
		Status: StatusPending,
	}
	s.jobs[job.ID] = job
	return job
}

func (s *JobStore) Get(id string) (*Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, found := s.jobs[id]
	return job, found
}

func (s *JobStore) UpdateStatus(id string, status JobStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, found := s.jobs[id]; found {
		job.Status = status
	}
}

func (s *JobStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
}

func (s *JobStore) SetResult(id, filePath string, fileSize int64, errStr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, found := s.jobs[id]; found {
		if errStr != "" {
			job.Status = StatusFailed
			job.Error = errStr
		} else {
			job.Status = StatusComplete
			job.FilePath = filePath
			job.FileSize = fileSize
		}

		go func() {
			time.Sleep(15 * time.Minute)
			slog.Info("deleting old job from memory", "jobID", id)
			s.Delete(id)
		}()
	}
}
