package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/domain"
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"github.com/go-go-golems/ai-in-action-app/internal/templates/pages"
	"github.com/labstack/echo/v4"
)

// EventHandler handles event-related requests
type EventHandler struct {
	eventRepo repository.EventRepository
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventRepo repository.EventRepository) *EventHandler {
	return &EventHandler{
		eventRepo: eventRepo,
	}
}

// RegisterRoutes registers the event routes
func (h *EventHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.HandleTimelinePage)
	e.GET("/events/add-form", h.HandleAddEventForm)
	e.POST("/events/add", h.HandleAddEvent)
}

// HandleTimelinePage renders the timeline page with upcoming and past events
func (h *EventHandler) HandleTimelinePage(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Get upcoming events
	upcomingEvents, err := h.eventRepo.GetUpcomingEvents(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get upcoming events: "+err.Error())
	}

	// Get past events
	pastEvents, err := h.eventRepo.GetPastEvents(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get past events: "+err.Error())
	}

	// Render the timeline page
	return pages.Timeline(upcomingEvents, pastEvents).Render(ctx, c.Response().Writer)
}

// HandleAddEventForm renders the form for adding a new event
func (h *EventHandler) HandleAddEventForm(c echo.Context) error {
	ctx := c.Request().Context()
	return pages.AddEventForm().Render(ctx, c.Response().Writer)
}

// HandleAddEvent handles the submission of a new event
func (h *EventHandler) HandleAddEvent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Parse form data
	title := c.FormValue("title")
	speaker := c.FormValue("speaker")
	description := c.FormValue("description")
	dateStr := c.FormValue("date")
	timeStr := c.FormValue("time")

	// Validate required fields
	if title == "" || speaker == "" || description == "" || dateStr == "" || timeStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "All fields are required")
	}

	// Parse date and time
	dateTimeStr := dateStr + "T" + timeStr + ":00"
	eventDate, err := time.Parse("2006-01-02T15:04:05", dateTimeStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid date or time format")
	}

	// Create new event
	event := domain.Event{
		Title:       title,
		Speaker:     speaker,
		Description: description,
		Date:        eventDate,
		IsUpcoming:  eventDate.After(time.Now()),
	}

	// Add event to repository
	_, err = h.eventRepo.AddEvent(ctx, event)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add event: "+err.Error())
	}

	// Get updated events for the response
	upcomingEvents, err := h.eventRepo.GetUpcomingEvents(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get upcoming events: "+err.Error())
	}

	pastEvents, err := h.eventRepo.GetPastEvents(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get past events: "+err.Error())
	}

	// Check if this is an HTMX request
	if c.Request().Header.Get("HX-Request") == "true" {
		// Return only the timeline content for HTMX requests
		return pages.TimelineContent(upcomingEvents, pastEvents).Render(ctx, c.Response().Writer)
	}

	// Return the full timeline page for regular requests
	return pages.Timeline(upcomingEvents, pastEvents).Render(ctx, c.Response().Writer)
}
