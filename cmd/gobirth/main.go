package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/calendar/file"
	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/calendar/google"
	clocksys "github.com/rafakmp18/gobirth/internal/gobirth/adapters/clock/system"
	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/message/template"
	"github.com/rafakmp18/gobirth/internal/gobirth/adapters/whatsapp/stdout"
	"github.com/rafakmp18/gobirth/internal/gobirth/application"
)

func main() {
	var (
		calendarFile       = flag.String("calendar-file", "", "Path to a JSON file with calendar events")
		dryRun             = flag.Bool("dry-run", true, "If true, messages will not be sent (only printed)")
		maxPerRun          = flag.Int("max", 10, "Maximum number of greetings to process per run")
		dateStr            = flag.String("date", "", "Run for a specific date (YYYY-MM-DD). Defaults to today.")
		emoji              = flag.String("emoji", "🎉", "Emoji to include in template messages")
		calendarProvider   = flag.String("calendar-provider", "file", "Calendar provider: file|google")
		googleCalendarName = flag.String("google-calendar", "gobirth", "Google Calendar name to use (e.g. gobirth)")
		googleCredentials  = flag.String("google-credentials", "", "Path to Google OAuth credentials.json (default ~/.config/gobirth/credentials.json)")
		googleToken        = flag.String("google-token", "", "Path to Google OAuth token.json (default ~/.config/gobirth/token.json)")
	)
	flag.Parse()

	cfgDir, _ := os.UserConfigDir()
	defaultCreds := cfgDir + "/gobirth/credentials.json"
	defaultToken := cfgDir + "/gobirth/token.json"

	credPath := *googleCredentials
	if credPath == "" {
		credPath = defaultCreds
	}

	tokenPath := *googleToken
	if tokenPath == "" {
		tokenPath = defaultToken
	}

	//TODO: Make it configurable
	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot load Europe/Madrid timezone:", err)
		os.Exit(1)
	}
	runDate, err := resolveRunDate(*dateStr, loc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	var cal application.CalendarProvider

	switch *calendarProvider {
	case "file":
		if *calendarFile == "" {
			fmt.Fprintln(os.Stderr, "error: --calendar-file is required when --calendar-provider=file")
			os.Exit(2)
		}
		cal = file.Provider{Path: *calendarFile}

	case "google":
		svc, err := google.NewCalendarService(context.Background(), google.AuthConfig{
			CredentialsPath: credPath,
			TokenPath:       tokenPath,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(2)
		}
		cal = &google.Provider{
			Svc:          svc,
			CalendarName: *googleCalendarName,
		}

	default:
		fmt.Fprintln(os.Stderr, "error: invalid --calendar-provider (use file|google)")
		os.Exit(2)
	}

	gen := template.Generator{Emoji: *emoji}
	sender := stdout.New(os.Stdout)
	clk := clocksys.Clock{}

	uc := application.RunDailyGreetings{
		Calendar:  cal,
		Parser:    application.EventParser{},
		Generator: gen,
		Sender:    sender,
		Clock:     fixedOrSystemClock{fixed: runDate, system: clk, useFixed: *dateStr != ""},
		MaxPerRun: *maxPerRun,
		DryRun:    *dryRun,
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
