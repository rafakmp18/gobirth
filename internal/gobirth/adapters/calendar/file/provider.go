package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/application"
)

type Provider struct {
	Path string
}

type eventDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
}

func (p Provider) EventsForDate(ctx context.Context, date time.Time, tag string) ([]application.CalendarEvent, error) {
	_ = ctx
	_ = tag

	f, err := os.Open(p.Path)
	if err != nil {
		return nil, fmt.Errorf("file calendar: open %s: %w", p.Path, err)
	}
	defer f.Close()

	dtos, err := decodeEvents(f)
	if err != nil {
		return nil, err
	}

	want := ymd(date)
	out := make([]application.CalendarEvent, 0, len(dtos))

	for _, d := range dtos {
		evDate, err := time.ParseInLocation("2006-01-02", d.StartDate, date.Location())
		if err != nil {
			return nil, fmt.Errorf("file calendar: invalid start_date for id=%s: %w", d.ID, err)
		}

		if ymd(evDate) != want {
			continue
		}

		out = append(out, application.CalendarEvent{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			StartDate:   evDate,
		})
	}

	return out, nil
}

func decodeEvents(r io.Reader) ([]eventDTO, error) {
	var dtos []eventDTO
	if err := json.NewDecoder(r).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("file calendar: decode json: %w", err)
	}
	return dtos, nil
}

func ymd(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%04d-%02d-%02d", y, int(m), d)
}
