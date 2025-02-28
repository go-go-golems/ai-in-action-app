package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/handlers"
	"github.com/go-go-golems/ai-in-action-app/internal/repository/mock"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize repositories
	eventRepo := mock.NewMockEventRepository()
	timerRepo := mock.NewMockTimerRepository()
	noteRepo := mock.NewMockNoteRepository()
	questionRepo := mock.NewMockQuestionRepository()

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	// Register handlers
	handlers.RegisterHandlers(e, eventRepo, timerRepo, noteRepo, questionRepo)

	// Start server in a goroutine
	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	log.Println("Server started at http://localhost:8080")

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
