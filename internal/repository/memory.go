package repository

import (
	"sync"

	"github.com/whiterage/webserver_go/pkg/models"
)

type MemoryRepo struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		tasks: make(map[string]*models.Task),
	}
}

func (r *MemoryRepo) Save(task *models.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
}

func (r *MemoryRepo) Get(id string) (*models.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	return task, ok
}

func (r *MemoryRepo) List(ids []string) []*models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]*models.Task, 0, len(ids))
	for _, id := range ids {
		if task, ok := r.tasks[id]; ok {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
