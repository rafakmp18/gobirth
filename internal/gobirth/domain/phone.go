package domain

import (
	"strings"
	"unicode"
)

type Phone struct {
	value string
}

func NewPhone(phone string) (Phone, error) {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return Phone{}, ErrMissingPhone
	}

	if !isValidE164(phone) {
		return Phone{}, ErrInvalidPhone
	}

	return Phone{value: phone}, nil
}

func (phone Phone) String() string {
	return phone.value
}

func isValidE164(phone string) bool {
	if !strings.HasPrefix(phone, "+") {
		return false
	}

	if len(phone) < 8 || len(phone) > 16 {
		return false
	}

	for _, r := range phone[1:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
