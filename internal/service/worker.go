package service

import (
	"context"
	"sync"

	"github.com/whiterage/webserver_go/pkg/models"
)

type WorkerPool struct {
	service *Service
	workers int
	wg      sync.WaitGroup
	stop    chan struct{}
}

func NewWorkerPool(service *Service, workers int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	return &WorkerPool{
		service: service,
		workers: workers,
		stop:    make(chan struct{}),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.workerLoop(ctx)
	}
}

func (wp *WorkerPool) workerLoop(ctx context.Context) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-wp.stop:
			return
		case task, ok := <-wp.service.queue:
			if !ok {
				return
			}
			wp.processTask(ctx, task)
		}
	}
}

func (wp *WorkerPool) processTask(ctx context.Context, task *models.Task) {
	task.Status = "processing"
	wp.service.repo.Save(task)

	for i := range task.Results {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result := wp.service.checker.Check(ctx, task.Results[i].URL)
		task.Results[i] = result
	}

	task.Status = "done"
	wp.service.repo.Save(task)
}

func (wp *WorkerPool) Stop() {
	close(wp.stop)
	wp.wg.Wait()
}
