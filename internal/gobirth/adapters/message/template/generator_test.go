package template

import (
	"context"
	"testing"
	"time"

	"github.com/rafakmp18/gobirth/internal/gobirth/application"
)

func TestTemplateGenerator_Generate(t *testing.T) {
	g := Generator{Emoji: "🎂"}

	msg, err := g.Generate(context.Background(), application.MessageInput{
		Name:    "Pepe",
		Context: "Pásalo genial!",
		Date:    time.Now(),
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if msg.Text() == "" {
		t.Fatalf("expected non-empty message")
	}
}
