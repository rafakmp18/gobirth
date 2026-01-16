package stdout

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

func TestSender_SendText_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	s := New(&buf)

	phone, err := domain.NewPhone("+34600111222")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.SendText(context.Background(), phone, "hola")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "TO: +34600111222") {
		t.Fatalf("expected output to contain phone, got:\n%s", out)
	}
	if !strings.Contains(out, "MSG:\nhola") {
		t.Fatalf("expected output to contain message, got:\n%s", out)
	}
}
