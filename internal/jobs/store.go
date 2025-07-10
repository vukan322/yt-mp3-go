package jobs

import (
	"context"
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
	ID       string             `json:"jobID"`
	VideoID  string             `json:"-"`
	Status   JobStatus          `json:"status"`
	FilePath string             `json:"filePath"`
	Error    string             `json:"error"`
	cancel   context.CancelFunc `json:"-"`
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

func (s *JobStore) Add(videoID string, cancel context.CancelFunc) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	job := &Job{
		ID:      uuid.New().String(),
		VideoID: videoID,
		Status:  StatusPending,
		cancel:  cancel,
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

func (s *JobStore) Cancel(id string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if job, found := s.jobs[id]; found {
		if job.cancel != nil {
			job.cancel()
		}
	}
}

func (s *JobStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
}

func (s *JobStore) SetResult(id, filePath string, errStr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, found := s.jobs[id]; found {
		if errStr != "" {
			job.Status = StatusFailed
			job.Error = errStr
		} else {
			job.Status = StatusComplete
			job.FilePath = filePath
		}

		job.cancel = nil

		go func() {
			time.Sleep(15 * time.Minute)
			slog.Info("deleting old job from memory", "jobID", id)
			s.Delete(id)
		}()
	}
}
