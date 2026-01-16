package application

import (
	"strings"

	"github.com/rafakmp18/gobirth/internal/gobirth/domain"
)

type EventParser struct{}

func (EventParser) Parse(e CalendarEvent) (domain.Contact, error) {
	name := strings.TrimSpace(e.Title)
	if name == "" {
		return domain.Contact{}, domain.ErrMissingName
	}

	phoneRaw, context := parseDescription(e.Description)

	phone, err := domain.NewPhone(phoneRaw)
	if err != nil {
		return domain.Contact{}, err
	}

	return domain.NewContact(name, phone, context)
}

func parseDescription(desc string) (phone string, context string) {
	lines := strings.Split(desc, "\n")

	var ctxLines []string
	var inContext bool

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		low := strings.ToLower(l)

		if strings.HasPrefix(low, "phone:") || strings.HasPrefix(low, "tel:") {
			inContext = false
			phone = strings.TrimSpace(l[strings.Index(l, ":")+1:])
			continue
		}

		if strings.HasPrefix(low, "context:") {
			inContext = true
			val := strings.TrimSpace(l[strings.Index(l, ":")+1:])
			if val != "" {
				ctxLines = append(ctxLines, val)
			}
			continue
		}

		if inContext {
			ctxLines = append(ctxLines, l)
		}
	}

	context = strings.TrimSpace(strings.Join(ctxLines, "\n"))
	return phone, context
}
