package application

import (
	"context"
	"testing"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

type fakeClock struct{ t time.Time }

func (f fakeClock) Now() time.Time { return f.t }

type fakeCalendar struct {
	events []CalendarEvent
	err    error
}

func (f fakeCalendar) EventsForDate(ctx context.Context, date time.Time) ([]CalendarEvent, error) {
	return f.events, f.err
}

type fakeGenerator struct {
	text string
	err  error
}

func (f fakeGenerator) Generate(ctx context.Context, in MessageInput) (domain.GreetingMessage, error) {
	if f.err != nil {
		return domain.GreetingMessage{}, f.err
	}

	return domain.NewGreetingMessage("Feliz cumple " + in.Name + "! " + f.text), nil
}

type fakeSender struct {
	sent []struct {
		to   string
		text string
	}
	err error
}

func (f *fakeSender) SendText(ctx context.Context, to domain.Phone, text string) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, struct {
		to   string
		text string
	}{to: to.String(), text: text})
	return nil
}

func TestRunDailyGreetings_SendsMessages(t *testing.T) {
	now := time.Date(2026, 1, 16, 9, 0, 0, 0, time.UTC)

	cal := fakeCalendar{
		events: []CalendarEvent{
			{ID: "1", Title: "Pepe", Description: "phone: +34600111222", StartDate: now},
			{ID: "2", Title: "Ana", Description: "phone: +34600333444\ncontext: amiga del curro", StartDate: now},
		},
	}

	sender := &fakeSender{}
	uc := RunDailyGreetings{
		Calendar:  cal,
		Parser:    EventParser{},
		Generator: fakeGenerator{text: "🎉"},
		Sender:    sender,
		Clock:     fakeClock{t: now},

		MaxPerRun: 10,
		DryRun:    false,
	}

	res := uc.Run(context.Background())

	if res.Total != 2 {
		t.Fatalf("expected total 2, got %d", res.Total)
	}
	if res.Sent != 2 {
		t.Fatalf("expected sent 2, got %d", res.Sent)
	}
	if len(sender.sent) != 2 {
		t.Fatalf("expected sender to send 2 messages, got %d", len(sender.sent))
	}
}

func TestRunDailyGreetings_DryRun_DoesNotSend(t *testing.T) {
	now := time.Date(2026, 1, 16, 9, 0, 0, 0, time.UTC)

	cal := fakeCalendar{
		events: []CalendarEvent{
			{ID: "1", Title: "Pepe", Description: "phone: +34600111222", StartDate: now},
		},
	}

	sender := &fakeSender{}
	uc := RunDailyGreetings{
		Calendar:  cal,
		Parser:    EventParser{},
		Generator: fakeGenerator{text: "🎉"},
		Sender:    sender,
		Clock:     fakeClock{t: now},

		MaxPerRun: 10,
		DryRun:    true,
	}

	res := uc.Run(context.Background())

	if res.Sent != 0 {
		t.Fatalf("expected sent 0, got %d", res.Sent)
	}
	if res.Skipped != 1 {
		t.Fatalf("expected skipped 1, got %d", res.Skipped)
	}
	if len(sender.sent) != 0 {
		t.Fatalf("expected sender to send 0 messages, got %d", len(sender.sent))
	}
}

func TestRunDailyGreetings_RespectsMaxPerRun(t *testing.T) {
	now := time.Date(2026, 1, 16, 9, 0, 0, 0, time.UTC)

	cal := fakeCalendar{
		events: []CalendarEvent{
			{ID: "1", Title: "A", Description: "phone: +34600000001", StartDate: now},
			{ID: "2", Title: "B", Description: "phone: +34600000002", StartDate: now},
			{ID: "3", Title: "C", Description: "phone: +34600000003", StartDate: now},
		},
	}

	sender := &fakeSender{}
	uc := RunDailyGreetings{
		Calendar:  cal,
		Parser:    EventParser{},
		Generator: fakeGenerator{text: "ok"},
		Sender:    sender,
		Clock:     fakeClock{t: now},
		MaxPerRun: 2,
		DryRun:    false,
	}

	res := uc.Run(context.Background())

	if res.Sent != 2 {
		t.Fatalf("expected sent 2, got %d", res.Sent)
	}
	if res.Skipped != 1 {
		t.Fatalf("expected skipped 1, got %d", res.Skipped)
	}
	if len(sender.sent) != 2 {
		t.Fatalf("expected sender to send 2 messages, got %d", len(sender.sent))
	}
}
