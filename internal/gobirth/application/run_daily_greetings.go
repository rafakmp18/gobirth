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

func (useCase RunDailyGreetings) Run(ctx context.Context) RunResult {
	now := useCase.Clock.Now()
	date := startOfDay(now)

	events, err := useCase.Calendar.EventsForDate(ctx, date, useCase.Tag)
	if err != nil {
		return RunResult{Failed: 1, Errors: []error{err}}
	}

	res := RunResult{Total: len(events)}

	limit := useCase.MaxPerRun
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}

	for i := 0; i < limit; i++ {
		ev := events[i]

		contact, err := useCase.Parser.Parse(ev)
		if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, err)
			continue
		}

		msg, err := useCase.Generator.Generate(ctx, MessageInput{
			Name:    contact.Name(),
			Context: contact.Context(),
			Date:    date,
		})
		if err != nil {
			res.Failed++
			res.Errors = append(res.Errors, err)
			continue
		}

		if useCase.DryRun {
			res.Skipped++
			continue
		}

		if err := useCase.Sender.SendText(ctx, contact.Phone(), msg.Text()); err != nil {
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
