package template

import (
	"context"
	"fmt"

	"github.com/rafakmp18/gobirth/internal/gobirth/application"
	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

type Generator struct {
	Emoji string
}

func (g Generator) Generate(ctx context.Context, in application.MessageInput) (domain.GreetingMessage, error) {
	_ = ctx

	emoji := g.Emoji
	if emoji == "" {
		emoji = "🎉"
	}

	var text string
	if in.Context != "" {
		text = fmt.Sprintf("¡Feliz cumpleaños, %s! %s\n%s", in.Name, emoji, in.Context)
	} else {
		text = fmt.Sprintf("¡Feliz cumpleaños, %s! %s", in.Name, emoji)
	}

	return domain.NewGreetingMessage(text), nil
}
