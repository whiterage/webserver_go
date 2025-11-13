package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/whiterage/webserver_go/internal/repository"
	"github.com/whiterage/webserver_go/pkg/models"
)

var ErrTaskNotFound = errors.New("task not found")

type Checker interface {
	Check(ctx context.Context, url string) models.LinkStatus
}

type Service struct {
	repo    *repository.MemoryRepo
	queue   chan *models.Task
	checker Checker
}

func NewService(repo *repository.MemoryRepo, checker Checker, queueSize int) *Service {
	if queueSize <= 0 {
		queueSize = 10
	}

	return &Service{
		repo:    repo,
		queue:   make(chan *models.Task, queueSize),
		checker: checker,
	}
}

func (s *Service) CreateTask(ctx context.Context, links []string) (string, error) {
	if len(links) == 0 {
		return "", errors.New("empty links payload")
	}

	task := &models.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		Status:    "pending",
		Results:   make([]models.LinkStatus, len(links)),
	}

	for i, url := range links {
		task.Results[i] = models.LinkStatus{
			URL:    url,
			Status: "pending",
		}
	}

	s.repo.Save(task)

	select {
	case s.queue <- task:
	case <-ctx.Done():
		return "", ctx.Err()
	}

	return task.ID, nil
}

func (s *Service) GetTask(id string) (*models.Task, error) {
	task, ok := s.repo.Get(id)
	if !ok {
		return nil, ErrTaskNotFound
	}
	return task, nil
}
