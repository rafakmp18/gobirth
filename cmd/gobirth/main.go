package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/calendar/file"
	clocksys "github.com/rafakmp18/gobirth/internal/gobirth/adapters/clock/system"
	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/message/template"
	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/whatsapp/stdout"
	"github.com/rafakmp18/gobirth/internal/gobirth/application"
)

func main() {
	var (
		calendarFile = flag.String("calendar-file", "", "Path to a JSON file with calendar events")
		tag          = flag.String("tag", "gobirth", "Calendar tag/label to filter events")
		dryRun       = flag.Bool("dry-run", true, "If true, messages will not be sent (only printed)")
		maxPerRun    = flag.Int("max", 10, "Maximum number of greetings to process per run")
		dateStr      = flag.String("date", "", "Run for a specific date (YYYY-MM-DD). Defaults to today.")
		emoji        = flag.String("emoji", "🎉", "Emoji to include in template messages")
	)
	flag.Parse()

	if *calendarFile == "" {
		fmt.Fprintln(os.Stderr, "error: --calendar-file is required for now")
		os.Exit(2)
	}

	//TODO: Make it configurable
	loc := time.Local
	runDate, err := resolveRunDate(*dateStr, loc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}


	cal := file.Provider{Path: *calendarFile}
	gen := template.Generator{Emoji: *emoji}
	sender := stdout.New(os.Stdout)
	clk := clocksys.Clock{}


	uc := application.RunDailyGreetings{
		Calendar:  cal,
		Parser:    application.EventParser{},
		Generator: gen,
		Sender:    sender,
		Clock:     fixedOrSystemClock{fixed: runDate, system: clk, useFixed: *dateStr != ""},
		Tag:       *tag,
		MaxPerRun:  *maxPerRun,
		DryRun:     *dryRun,
	}

	res := uc.Run(context.Background())

	fmt.Printf("Run finished. total=%d sent=%d skipped=%d failed=%d\n",
		res.Total, res.Sent, res.Skipped, res.Failed)

	if len(res.Errors) > 0 {
		fmt.Println("Errors:")
		for _, e := range res.Errors {
			fmt.Printf("- %v\n", e)
		}
		if res.Failed > 0 {
			os.Exit(1)
		}
	}
}


func resolveRunDate(dateStr string, loc *time.Location) (time.Time, error) {
	now := time.Now().In(loc)

	if dateStr == "" {
		y, m, d := now.Date()
		return time.Date(y, m, d, 9, 0, 0, 0, loc), nil
	}

	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --date %q, expected YYYY-MM-DD", dateStr)
	}


	y, m, d := t.Date()
	return time.Date(y, m, d, 9, 0, 0, 0, loc), nil
}


type fixedOrSystemClock struct {
	fixed    time.Time
	system   clocksys.Clock
	useFixed bool
}

func (c fixedOrSystemClock) Now() time.Time {
	if c.useFixed {
		return c.fixed
	}
	return c.system.Now()
}
