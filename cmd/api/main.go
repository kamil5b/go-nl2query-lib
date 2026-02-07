package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	_ = godotenv.Load()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis: %v", err)
	}

	// wsRepo := repository.NewRedisWorkspaceRepository(redisClient)
	// statusRepo := repository.NewRedisStatusRepository(redisClient)

	// hashGen := hash.New()
	// metadataExt := postgres.NewMetadataExtractor()
	// checksumGen := postgres.NewChecksumGenerator()
	// taskQueue := queue.NewAsynqTaskQueue(redisClient)

	// // ucContainer := usecase.NewContainer(wsRepo, statusRepo, hashGen, metadataExt, checksumGen, taskQueue)

	// // workspaceHandler := handler.NewWorkspaceHandler(
	// // 	ucContainer.AddWorkspace(),
	// // 	ucContainer.ListWorkspaces(),
	// // 	ucContainer.GetStatus(),
	// // )

	// mux := http.NewServeMux()
	// mux.HandleFunc("POST /workspace", workspaceHandler.AddWorkspace)
	// mux.HandleFunc("GET /workspace", workspaceHandler.ListWorkspaces)
	// mux.HandleFunc("GET /status/{workspace_id}", workspaceHandler.GetWorkspaceStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: ":" + port,
		// Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}

	_ = redisClient.Close()
	log.Println("done")
}
