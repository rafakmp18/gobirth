package application

import "time"

type CalendarEvent struct {
	ID          string
	Title       string
	Description string
	StartDate   time.Time
}
