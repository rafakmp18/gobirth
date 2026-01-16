package application

import (
	"context"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

type RunDailyGreetings struct {
	Calendar  CalendarProvider
	Parser    EventParser
	Generator MessageGenerator
	Sender    WhatsAppSender
	Clock     Clock

	Tag       string
	MaxPerRun int
	DryRun    bool
}

type RunResult struct {
	Total   int
	Sent    int
	Skipped int
	Failed  int
	Errors  []error
}

func (uc RunDailyGreetings) Run(ctx context.Context) RunResult {
	now := uc.Clock.Now()
	date := startOfDay(now)

	events, err := uc.Calendar.EventsForDate(ctx, date, uc.Tag)
	if err != nil {
		return RunResult{Failed: 1, Errors: []error{err}}
	}

	res := RunResult{Total: len(events)}

	limit := uc.MaxPerRun
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}

	for i := 0; i < limit; i++ {
		ev := events[i]

		contact, err := uc.Parser.Parse(ev)
		if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, err)
			continue
		}

		msg, err := uc.Generator.Generate(ctx, MessageInput{
			Name:    contact.Name(),
			Context: contact.Context(),
			Date:    date,
		})
		if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, err)
			continue
		}

		if uc.DryRun {
			res.Skipped++
			continue
		}

		if err := uc.Sender.SendText(ctx, contact.Phone(), msg.Text()); err != nil {
			res.Failed++
			res.Errors = append(res.Errors, err)
			continue
		}

		res.Sent++
	}

	if len(events) > limit {
		res.Skipped += len(events) - limit
	}

	return res
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}


var _ = domain.ErrMissingName
