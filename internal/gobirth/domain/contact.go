package domain

import "strings"

type Contact struct {
	name    string
	phone   Phone
	context string
}

func NewContact(name string, phone Phone, context string) (Contact, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Contact{}, ErrMissingName
	}

	return Contact{
		name:    name,
		phone:   phone,
		context: strings.TrimSpace(context),
	}, nil
}

func (contact Contact) Name() string {
	return contact.name
}

func (contact Contact) Phone() Phone {
	return contact.phone
}

func (contact Contact) Context() string {
	return contact.context
}
