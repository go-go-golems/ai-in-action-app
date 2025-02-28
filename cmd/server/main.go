package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/handlers"
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"github.com/go-go-golems/ai-in-action-app/internal/repository/mock"
	"github.com/go-go-golems/ai-in-action-app/internal/repository/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// Flags
	useSQLite  bool
	dbPath     string
	serverPort int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "AI in Action App Server",
		Long:  `Web server for the AI in Action application with timeline, timer, notes, and question queue features.`,
		RunE:  runServer,
	}

	// Add flags
	rootCmd.Flags().BoolVar(&useSQLite, "sqlite", false, "Use SQLite repositories instead of mock repositories")
	rootCmd.Flags().StringVar(&dbPath, "db-path", "ai-in-action.db", "Path to SQLite database file (only used with --sqlite)")
	rootCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	var (
		eventRepo     repository.EventRepository
		timerRepo     repository.TimerRepository
		noteRepo      repository.NoteRepository
		questionRepo  repository.QuestionRepository
		sqliteFactory *sqlite.RepositoryFactory
		err           error
	)

	// Initialize repositories based on flag
	if useSQLite {
		log.Printf("Using SQLite repositories with database: %s\n", dbPath)
		sqliteFactory, err = sqlite.NewRepositoryFactory(dbPath)
		if err != nil {
			return errors.Wrap(err, "failed to initialize SQLite repositories")
		}
		defer sqliteFactory.Close()

		eventRepo = sqliteFactory.GetEventRepository()
		timerRepo = sqliteFactory.GetTimerRepository()
		noteRepo = sqliteFactory.GetNoteRepository()
		questionRepo = sqliteFactory.GetQuestionRepository()
	} else {
		log.Println("Using mock repositories")
		eventRepo = mock.NewMockEventRepository()
		timerRepo = mock.NewMockTimerRepository()
		noteRepo = mock.NewMockNoteRepository()
		questionRepo = mock.NewMockQuestionRepository()
	}

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	// Register handlers
	handlers.RegisterHandlers(e, eventRepo, timerRepo, noteRepo, questionRepo)

	// Start server in a goroutine
	go func() {
		serverAddr := fmt.Sprintf(":%d", serverPort)
		if err := e.Start(serverAddr); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	log.Printf("Server started at http://localhost:%d\n", serverPort)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "error during server shutdown")
	}

	return nil
}
