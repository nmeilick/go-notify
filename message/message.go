// Package message provides a simple way to create and manage messages.
package message

// Message represents a message with various attributes.
type Message struct {
	Text     string   `json:"text,omitempty"`
	Format   string   `json:"format,omitempty"`
	Subject  string   `json:"subject,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
	File     string   `json:"file,omitempty"`
}

// Option is a function that applies a specific attribute to a Message.
type Option func(*Message)

// Text sets the text content of a Message.
func Text(text string) Option {
	return func(m *Message) {
		m.Text = text
	}
}

// Format sets the text format of a Message.
func Format(format string) Option {
	return func(m *Message) {
		m.Format = format
	}
}

// Subject sets the subject of a Message.
func Subject(subject string) Option {
	return func(m *Message) {
		m.Subject = subject
	}
}

// Summary sets the summary of a Message.
func Summary(summary string) Option {
	return func(m *Message) {
		m.Summary = summary
	}
}

// Keywords sets the keywords of a Message.
func Keywords(keywords ...string) Option {
	return func(m *Message) {
		m.Keywords = keywords
	}
}

// File sets the files associated with a Message.
func File(file string) Option {
	return func(m *Message) {
		m.File = file
	}
}

// NewMessage creates a new Message with the given text and applies the provided options.
func NewMessage(text string, opts ...Option) *Message {
	msg := &Message{Text: text}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}
