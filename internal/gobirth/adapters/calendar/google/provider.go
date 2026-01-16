package google

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/application"
	"google.golang.org/api/calendar/v3"
)

type Provider struct {
	Svc          *calendar.Service
	CalendarName string

	mu          sync.Mutex
	cachedCalID string
}

func (p *Provider) EventsForDate(ctx context.Context, date time.Time) ([]application.CalendarEvent, error) {
	if p.Svc == nil {
		return nil, fmt.Errorf("google calendar: nil service")
	}
	if strings.TrimSpace(p.CalendarName) == "" {
		return nil, fmt.Errorf("google calendar: CalendarName is required")
	}

	calID, err := p.calendarID(ctx)
	if err != nil {
		return nil, err
	}

	start := startOfDay(date)
	end := start.Add(24 * time.Hour)

	call := p.Svc.Events.List(calID).
		Context(ctx).
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime")

	events, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("google calendar: events.list: %w", err)
	}

	out := make([]application.CalendarEvent, 0, len(events.Items))
	for _, ev := range events.Items {
		out = append(out, application.CalendarEvent{
			ID:          ev.Id,
			Title:       ev.Summary,
			Description: ev.Description,
			StartDate:   parseEventStart(date.Location(), ev),
		})
	}

	return out, nil
}

func (p *Provider) calendarID(ctx context.Context) (string, error) {
	p.mu.Lock()
	if p.cachedCalID != "" {
		id := p.cachedCalID
		p.mu.Unlock()
		return id, nil
	}
	p.mu.Unlock()

	name := strings.TrimSpace(p.CalendarName)

	var pageToken string
	for {
		call := p.Svc.CalendarList.List().Context(ctx)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		cl, err := call.Do()
		if err != nil {
			return "", fmt.Errorf("google calendar: calendarList.list: %w", err)
		}

		for _, item := range cl.Items {

			if strings.EqualFold(strings.TrimSpace(item.Summary), name) {
				p.mu.Lock()
				p.cachedCalID = item.Id
				p.mu.Unlock()
				return item.Id, nil
			}
		}

		if cl.NextPageToken == "" {
			break
		}
		pageToken = cl.NextPageToken
	}

	return "", fmt.Errorf("google calendar: calendar %q not found (create it in Google Calendar)", name)
}

func parseEventStart(loc *time.Location, ev *calendar.Event) time.Time {
	if ev.Start == nil {
		return time.Time{}
	}


	if ev.Start.Date != "" {
		t, err := time.ParseInLocation("2006-01-02", ev.Start.Date, loc)
		if err == nil {
			return t
		}
	}

	if ev.Start.DateTime != "" {
		t, err := time.Parse(time.RFC3339, ev.Start.DateTime)
		if err == nil {
			return t.In(loc)
		}
	}

	return time.Time{}
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
