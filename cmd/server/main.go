package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/whiterage/webserver_go/internal/repository"
	"github.com/whiterage/webserver_go/internal/service"
	"github.com/whiterage/webserver_go/internal/worker"
)

func main() {
	repo := repository.NewMemoryRepo()
	checker := worker.NewHTTPChecker(5 * time.Second)
	svc := service.NewService(repo, checker, 20)
	pool := service.NewWorkerPool(svc, 4)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer pool.Stop()

	pool.Start(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
