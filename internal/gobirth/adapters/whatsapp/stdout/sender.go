package stdout

import (
	"context"
	"fmt"
	"io"

	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

type Sender struct {
	Out io.Writer
}

func New(out io.Writer) Sender {
	return Sender{Out: out}
}

func (s Sender) SendText(ctx context.Context, to domain.Phone, text string) error {
	_ = ctx

	out := s.Out
	if out == nil {
		return fmt.Errorf("stdout sender: Out writer is nil")
	}

	_, err := fmt.Fprintf(out, "---- GOBIRTH (DRY WHATSAPP) ----\nTO: %s\nMSG:\n%s\n-------------------------------\n\n",
		to.String(),
		text,
	)
	return err
}
