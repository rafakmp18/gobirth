package application

import "testing"

func TestEventParser_Parse_OK_WithContext(t *testing.T) {
	p := EventParser{}

	ev := CalendarEvent{
		ID:    "1",
		Title: "Pepe",
		Description: `
phone: +34600111222
context: le encanta el fitness
y el k1
`,
	}

	c, err := p.Parse(ev)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if c.Name() != "Pepe" {
		t.Fatalf("expected name Pepe, got %q", c.Name())
	}

	if c.Phone().String() != "+34600111222" {
		t.Fatalf("expected phone +34600111222, got %q", c.Phone().String())
	}

	wantCtx := "le encanta el fitness\ny el k1"
	if c.Context() != wantCtx {
		t.Fatalf("expected context %q, got %q", wantCtx, c.Context())
	}
}

func TestEventParser_Parse_OK_NoContext(t *testing.T) {
	p := EventParser{}

	ev := CalendarEvent{
		ID:          "1",
		Title:       "Mamá",
		Description: "phone: +34611112222",
	}

	c, err := p.Parse(ev)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if c.Context() != "" {
		t.Fatalf("expected empty context, got %q", c.Context())
	}
}

func TestEventParser_Parse_Fails_WhenMissingPhone(t *testing.T) {
	p := EventParser{}

	ev := CalendarEvent{
		ID:          "1",
		Title:       "Pepe",
		Description: "context: hola",
	}

	_, err := p.Parse(ev)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
