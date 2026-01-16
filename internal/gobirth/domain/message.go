package domain

import "strings"

type GreetingMessage struct {
	text string
}

func NewGreetingMessage(text string) GreetingMessage {
	return GreetingMessage{text: strings.TrimSpace(text)}
}

func (message GreetingMessage) Text() string {
	return message.text
}
