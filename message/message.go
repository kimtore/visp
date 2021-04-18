// Package message provides simple text message communication.
package message

import (
	"fmt"
)

// Message is a message passed from anywhere inside PMS, relayed to the user
// through the statusbar.
type Message struct {
	Text     string
	Severity int
	Type     int
}

// Message severities. INFO messages and above will end up in the statusbar.
const (
	Debug = iota
	Info
	Error
)

// Message types.
const (
	Normal = iota
	SequenceText
)

// format formats using Sprintf, and returns a new Message.
func format(severity int, t int, format string, a ...interface{}) Message {
	return Message{
		Text:     fmt.Sprintf(format, a...),
		Severity: severity,
		Type:     t,
	}
}

// Format returns a normal info message.
func Format(fmt string, a ...interface{}) Message {
	return format(Info, Normal, fmt, a...)
}

// Errorf returns a normal error message.
func Errorf(fmt string, a ...interface{}) Message {
	return format(Error, Normal, fmt, a...)
}

// Sequencef returns a sequence text message.
func Sequencef(fmt string, a ...interface{}) Message {
	return format(Info, SequenceText, fmt, a...)
}
