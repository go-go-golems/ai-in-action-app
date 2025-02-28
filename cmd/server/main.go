package main

import (
	"fmt"
	"log"

	"github.com/wesen/ai-in-action-app/internal/repository/mock"
)

func main() {
	// Initialize repositories
	eventRepo := mock.NewMockEventRepository()
	timerRepo := mock.NewMockTimerRepository()
	noteRepo := mock.NewMockNoteRepository()
	questionRepo := mock.NewMockQuestionRepository()

	// Verify repositories are working
	fmt.Println("=== AI in Action App ===")

	// Test event repository
	fmt.Println("\nUpcoming Events:")
	for _, event := range eventRepo.GetUpcomingEvents() {
		fmt.Printf("- %s by %s on %s\n", event.Title, event.Speaker, event.Date.Format("Jan 2, 2006"))
	}

	fmt.Println("\nPast Events:")
	for _, event := range eventRepo.GetPastEvents() {
		fmt.Printf("- %s by %s on %s\n", event.Title, event.Speaker, event.Date.Format("Jan 2, 2006"))
	}

	// Test timer repository
	timer := timerRepo.GetTimer()
	fmt.Printf("\nTimer: %v remaining, running: %v\n", timer.RemainingTime, timer.IsRunning)

	// Test note repository
	note := noteRepo.GetNote(1)
	fmt.Printf("\nNote (Page %d/%d): %s\n", note.PageNumber, note.TotalPages, note.Content)

	// Test question repository
	fmt.Println("\nQuestions:")
	if len(questionRepo.GetQuestions()) == 0 {
		fmt.Println("No questions submitted yet.")
	}

	log.Println("Repository setup complete!")
}
