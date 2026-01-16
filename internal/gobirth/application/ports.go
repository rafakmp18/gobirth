package application

import (
	"context"
	"time"
	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

type CalendarProvider interface {
	EventsForDate(ctx context.Context, date time.Time, tag string) ([]CalendarEvent, error)
}

type MessageInput struct {
	Name    string
	Context string
	Date    time.Time
}

type MessageGenerator interface {
	Generate(ctx context.Context, in MessageInput) (domain.GreetingMessage, error)
}

type WhatsAppSender interface {
	SendText(ctx context.Context, to domain.Phone, text string) error
}

type Clock interface {
	Now() time.Time
}
