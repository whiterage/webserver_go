package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/whiterage/webserver_go/internal/http"
	"github.com/whiterage/webserver_go/internal/repository"
	"github.com/whiterage/webserver_go/internal/service"
	"github.com/whiterage/webserver_go/internal/worker"
)

func main() {
	repo := repository.NewMemoryRepo()
	checker := worker.NewHTTPChecker(5 * time.Second)
	svc := service.NewService(repo, checker, 20)
	pool := service.NewWorkerPool(svc, 4)
	handlers := api.NewHandlers(svc)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool.Start(ctx)

	handlers.Register(mux)

	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutdown: stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	pool.Stop()

	log.Println("shutdown: complete")
}
