# AI in Action App Implementation Plan

## HTML Structure as YAML DSL

### Timeline Page (01-timeline.png)
```yaml
page:
  header:
    title: "AI in Action Group"
    navigation:
      - label: "Timeline"
        active: true
      - label: "Timer & Notes"
        link: "/timer"
      - label: "Questions"
        link: "/questions"
  content:
    sections:
      - title: "Upcoming Talks"
        actions:
          - label: "Add Event"
            type: "button"
            class: "primary"
        items:
          - type: "event"
            title: "Generative AI for Scientific Discovery"
            speaker: "Dr. Alex Chen"
            description: "Exploring how generative AI models can accelerate scientific discovery in various domains."
            date: "Thursday, Mar 6, 2025"
          - type: "event"
            title: "Ethical Considerations in AI Development"
            speaker: "Prof. Maya Johnson"
            description: "Discussing the ethical frameworks necessary for responsible AI development."
            date: "Thursday, Mar 13, 2025"
          - type: "event"
            title: "New AI Workshop"
            speaker: "Add speaker name"
            description: "Add description here"
            date: "Thursday, Mar 20, 2025"
          # More placeholder events...
      - title: "Past Talks"
        items:
          - type: "event"
            title: "Multimodal Learning in AI"
            speaker: "Sam Rodriguez"
            description: "How combining different data modalities can enhance AI model capabilities."
            date: "Thursday, Feb 20, 2025"
          - type: "event"
            title: "Reinforcement Learning from Human Feedback"
            speaker: "Dr. Jamie Park"
            description: "Deep dive into how RLHF is transforming the alignment of AI systems."
            date: "Thursday, Feb 13, 2025"
```

### Timer & Notes Page (03-timer.png)
```yaml
page:
  header:
    title: "AI in Action Group"
    navigation:
      - label: "Timeline"
        link: "/"
      - label: "Timer & Notes"
        active: true
      - label: "Questions"
        link: "/questions"
  content:
    sections:
      - title: "Countdown Timer"
        components:
          - type: "timer"
            display: "13:38"
            actions:
              - label: "Pause"
                type: "button"
                class: "danger"
              - label: "Reset"
                type: "button"
                class: "secondary"
          - type: "preset-times"
            label: "Set Custom Time (minutes)"
            options:
              - value: 5
              - value: 10
              - value: 15
              - value: 20
              - value: 30
      - title: "Speaker Notes"
        components:
          - type: "textarea"
            placeholder: "Questions and discussion time"
            rows: 10
          - type: "pagination"
            current: 9
            total: 10
            actions:
              - label: "Previous"
                type: "button"
                class: "primary"
              - label: "Next"
                type: "button"
                class: "primary"
```

### Question Queue Page (04-question-queue.png)
```yaml
page:
  header:
    title: "AI in Action"
    navigation:
      - label: "Timeline"
        link: "/"
      - label: "Countdown & Notes"
        link: "/timer"
      - label: "Question Queue"
        active: true
  content:
    sections:
      - title: "Question Queue"
        components:
          - type: "form"
            fields:
              - label: "Your Name"
                type: "text"
                placeholder: "Enter your name"
              - label: "Your Question"
                type: "text"
                placeholder: "Enter your question"
            actions:
              - label: "Submit Question"
                type: "submit"
                class: "primary"
      - title: "Current Questions"
        components:
          - type: "question-list"
            empty_message: "No questions submitted yet."
            items: []  # Will be populated dynamically
```

## Object Types (Domain Models)

```yaml
models:
  - name: Event
    fields:
      - name: ID
        type: uint
      - name: Title
        type: string
      - name: Speaker
        type: string
      - name: Description
        type: string
      - name: Date
        type: time.Time
      - name: IsUpcoming
        type: bool
  
  - name: Timer
    fields:
      - name: ID
        type: uint
      - name: Duration
        type: time.Duration
      - name: RemainingTime
        type: time.Duration
      - name: IsRunning
        type: bool
      - name: LastStartedAt
        type: time.Time
  
  - name: Note
    fields:
      - name: ID
        type: uint
      - name: Content
        type: string
      - name: PageNumber
        type: int
      - name: TotalPages
        type: int
  
  - name: Question
    fields:
      - name: ID
        type: uint
      - name: Name
        type: string
      - name: Content
        type: string
      - name: SubmittedAt
        type: time.Time
      - name: Answered
        type: bool
```

## Repository Interfaces

```yaml
repositories:
  - name: EventRepository
    methods:
      - name: GetUpcomingEvents
        return: []Event
      - name: GetPastEvents
        return: []Event
      - name: AddEvent
        params:
          - name: event
            type: Event
        return: Event
      - name: UpdateEvent
        params:
          - name: event
            type: Event
        return: bool
  
  - name: TimerRepository
    methods:
      - name: GetTimer
        return: Timer
      - name: UpdateTimer
        params:
          - name: timer
            type: Timer
        return: bool
      - name: ResetTimer
        params:
          - name: duration
            type: time.Duration
        return: Timer
  
  - name: NoteRepository
    methods:
      - name: GetNote
        params:
          - name: pageNumber
            type: int
        return: Note
      - name: SaveNote
        params:
          - name: note
            type: Note
        return: bool
  
  - name: QuestionRepository
    methods:
      - name: GetQuestions
        return: []Question
      - name: AddQuestion
        params:
          - name: question
            type: Question
        return: Question
      - name: MarkAsAnswered
        params:
          - name: id
            type: uint
        return: bool
```

## API Handlers (with HTMX support)

```yaml
handlers:
  # Timeline Handlers
  - path: "/"
    method: GET
    handler: HandleTimelinePage
    description: "Renders the main timeline page"
  
  - path: "/events/upcoming"
    method: GET
    handler: HandleUpcomingEvents
    description: "Returns HTML fragment of upcoming events (for HTMX)"
  
  - path: "/events/past"
    method: GET
    handler: HandlePastEvents
    description: "Returns HTML fragment of past events (for HTMX)"
  
  - path: "/events/add"
    method: POST
    handler: HandleAddEvent
    description: "Adds a new event and returns updated event list"
  
  # Timer Handlers
  - path: "/timer"
    method: GET
    handler: HandleTimerPage
    description: "Renders the timer and notes page"
  
  - path: "/timer/display"
    method: GET
    handler: HandleTimerDisplay
    description: "Returns HTML fragment with current timer (for HTMX)"
  
  - path: "/timer/toggle"
    method: POST
    handler: HandleTimerToggle
    description: "Toggles timer between running/paused and returns updated timer"
  
  - path: "/timer/reset"
    method: POST
    handler: HandleTimerReset
    description: "Resets timer and returns updated timer"
  
  - path: "/timer/set"
    method: POST
    handler: HandleTimerSet
    description: "Sets timer to specific duration and returns updated timer"
  
  # Notes Handlers
  - path: "/notes"
    method: GET
    handler: HandleGetNote
    description: "Returns HTML fragment with current note (for HTMX)"
  
  - path: "/notes/save"
    method: POST
    handler: HandleSaveNote
    description: "Saves current note content"
  
  - path: "/notes/navigate"
    method: POST
    handler: HandleNavigateNotes
    description: "Navigates to previous/next note and returns updated note"
  
  # Question Queue Handlers
  - path: "/questions"
    method: GET
    handler: HandleQuestionsPage
    description: "Renders the question queue page"
  
  - path: "/questions/list"
    method: GET
    handler: HandleQuestionsList
    description: "Returns HTML fragment with current questions (for HTMX)"
  
  - path: "/questions/add"
    method: POST
    handler: HandleAddQuestion
    description: "Adds a new question and returns updated question list"
  
  - path: "/questions/answer"
    method: POST
    handler: HandleAnswerQuestion
    description: "Marks a question as answered and returns updated question list"
```

## Directory Structure

  ```
  ├── cmd
  │   └── server
  │       └── main.go
  ├── internal
  │   ├── domain
  │   │   └── models.go
  │   ├── repository
  │   │   ├── interfaces.go
  │   │   └── mock
  │   │       └── repositories.go
  │   ├── handlers
  │   │   ├── events.go
  │   │   ├── timer.go
  │   │   ├── notes.go
  │   │   └── questions.go
  │   └── templates
  │       ├── components
  │       │   ├── event.templ
  │       │   ├── timer.templ
  │       │   ├── notes.templ
  │       │   └── questions.templ
  │       ├── layouts
  │       │   └── base.templ
  │       └── pages
  │           ├── timeline.templ
  │           ├── timer.templ
  │           └── questions.templ
  ├── static
  │   ├── css
  │   │   └── custom.css
  │   └── js
  │       └── app.js
  ├── go.mod
  ├── go.sum
  └── README.md
  ```
## Implementation Plan

### 1. Project Setup

- [x] Initialize Go module
  ```bash
  mkdir -p cmd/server
  go mod init github.com/go-go-golems/ai-in-action-app
  ```

- [x] Install dependencies
  ```bash
  go get -u github.com/labstack/echo/v4
  go get -u github.com/a-h/templ
  go get -u github.com/spf13/cobra
  ```

### 2. Domain Models Implementation

- [x] Define domain models in `internal/domain/models.go`
  ```go
  package domain

  import (
      "time"
  )

  type Event struct {
      ID          uint
      Title       string
      Speaker     string
      Description string
      Date        time.Time
      IsUpcoming  bool
  }

  // Define other models (Timer, Note, Question)
  ```

### 3. Repository Interfaces

- [x] Define repository interfaces in `internal/repository/interfaces.go`
  ```go
  package repository

  import (
      "github.com/go-go-golems/ai-in-action-app/internal/domain"
      "time"
  )

  type EventRepository interface {
      GetUpcomingEvents() []domain.Event
      GetPastEvents() []domain.Event
      AddEvent(event domain.Event) domain.Event
      UpdateEvent(event domain.Event) bool
  }

  // Define other repository interfaces
  ```

### 4. Mock Repository Implementation

- [x] Implement mock repositories in `internal/repository/mock/repositories.go`
  ```go
  package mock

  import (
      "github.com/go-go-golems/ai-in-action-app/internal/domain"
      "github.com/go-go-golems/ai-in-action-app/internal/repository"
      "time"
  )

  type MockEventRepository struct {
      events []domain.Event
  }

  var _ repository.EventRepository = &MockEventRepository{}

  // Implement methods with sample data
  func (m *MockEventRepository) GetUpcomingEvents() []domain.Event {
      // Return mock upcoming events
  }

  // Implement other repository methods and types
  ```

### 5. Templ Templates

- [ ] Create base layout in `internal/templates/layouts/base.templ`
  ```go
  package layouts

  templ Base(title string, activeNav string) {
      <!DOCTYPE html>
      <html lang="en">
          <head>
              <meta charset="UTF-8"/>
              <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
              <title>{ title }</title>
              <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet"/>
              <script src="https://unpkg.com/htmx.org@1.9.2"></script>
              <link href="/static/css/custom.css" rel="stylesheet"/>
          </head>
          <body>
              <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
                  <div class="container">
                      <a class="navbar-brand" href="/">AI in Action Group</a>
                      <div class="navbar-nav">
                          <!-- Navigation items -->
                      </div>
                  </div>
              </nav>
              <div class="container mt-4">
                  { children... }
              </div>
              <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
              <script src="/static/js/app.js"></script>
          </body>
      </html>
  }
  ```

- [ ] Create page templates for timeline, timer, and questions

### 6. HTTP Handlers

- [ ] Implement handlers for each endpoint defined in the API handlers section

### 7. Main Application

- [ ] Implement main.go to wire everything together
  ```go
  package main

  import (
      "github.com/labstack/echo/v4"
      "github.com/labstack/echo/v4/middleware"
      "github.com/go-go-golems/ai-in-action-app/internal/handlers"
      "github.com/go-go-golems/ai-in-action-app/internal/repository/mock"
  )

  func main() {
      // Initialize repositories
      eventRepo := &mock.MockEventRepository{}
      timerRepo := &mock.MockTimerRepository{}
      noteRepo := &mock.MockNoteRepository{}
      questionRepo := &mock.MockQuestionRepository{}

      // Initialize Echo
      e := echo.New()
      e.Use(middleware.Logger())
      e.Use(middleware.Recover())
      e.Static("/static", "static")

      // Register handlers
      handlers.RegisterHandlers(e, eventRepo, timerRepo, noteRepo, questionRepo)

      // Start server
      e.Logger.Fatal(e.Start(":8080"))
  }
  ```

### 8. HTMX Integration

- [ ] Create JavaScript for HTMX interactions in `static/js/app.js`
  ```javascript
  // Timer auto-refresh
  document.addEventListener('DOMContentLoaded', function() {
      // Set up timer polling if on timer page
      if (document.getElementById('timer-display')) {
          setInterval(() => {
              htmx.trigger('#timer-display', 'refresh');
          }, 1000);
      }
  });
  ```

### 9. Testing

- [ ] Write unit tests for repositories and handlers
- [ ] Write integration tests for API endpoints

### 10. Documentation

- [ ] Update README.md with setup and usage instructions
- [ ] Add API documentation

## Implementation Order

1. Set up project structure and dependencies
2. Implement domain models and repository interfaces
3. Create mock repositories with sample data
4. Implement templ templates for all pages
5. Create HTTP handlers for each endpoint
6. Wire everything together in main.go
7. Add HTMX interactions
8. Test and refine
9. Add documentation 