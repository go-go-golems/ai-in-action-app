package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wesen/ai-in-action-app/internal/repository/mock"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initialize repositories
	eventRepo := mock.NewMockEventRepository()
	timerRepo := mock.NewMockTimerRepository()
	noteRepo := mock.NewMockNoteRepository()
	questionRepo := mock.NewMockQuestionRepository()

	// Verify repositories are working
	fmt.Println("=== AI in Action App ===")

	// Test event repository
	fmt.Println("\nUpcoming Events:")
	upcomingEvents, err := eventRepo.GetUpcomingEvents(ctx)
	if err != nil {
		log.Fatalf("Error getting upcoming events: %v", err)
	}
	for _, event := range upcomingEvents {
		fmt.Printf("- %s by %s on %s\n", event.Title, event.Speaker, event.Date.Format("Jan 2, 2006"))
	}

	fmt.Println("\nPast Events:")
	pastEvents, err := eventRepo.GetPastEvents(ctx)
	if err != nil {
		log.Fatalf("Error getting past events: %v", err)
	}
	for _, event := range pastEvents {
		fmt.Printf("- %s by %s on %s\n", event.Title, event.Speaker, event.Date.Format("Jan 2, 2006"))
	}

	// Test timer repository
	timer, err := timerRepo.GetTimer(ctx)
	if err != nil {
		log.Fatalf("Error getting timer: %v", err)
	}
	fmt.Printf("\nTimer: %v remaining, running: %v\n", timer.RemainingTime, timer.IsRunning)

	// Test note repository
	note, err := noteRepo.GetNote(ctx, 1)
	if err != nil {
		log.Fatalf("Error getting note: %v", err)
	}
	fmt.Printf("\nNote (Page %d/%d): %s\n", note.PageNumber, note.TotalPages, note.Content)

	// Test question repository
	fmt.Println("\nQuestions:")
	questions, err := questionRepo.GetQuestions(ctx)
	if err != nil {
		log.Fatalf("Error getting questions: %v", err)
	}
	if len(questions) == 0 {
		fmt.Println("No questions submitted yet.")
	}

	log.Println("Repository setup complete!")
}
