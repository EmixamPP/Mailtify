package message

import (
	"net/smtp"
)

// Messenger is a wrapper allowing to send messages.
type Messenger struct {
	send func(subject string, body string, to []string) error
}

// Create creates a Messenger.
func Create(from, username, password, host, port string) *Messenger {
	return &Messenger{
		send: func(subject string, body string, to []string) error {
			var msg string
			msg += "From:" + from + "\n"
			msg += "Subject:" + subject + "\n"
			msg += "\n" + body + "\n"
			auth := smtp.PlainAuth("", username, password, host)
			return smtp.SendMail(host+":"+port, auth, from, to, []byte(msg))
		},
	}
}

// SendMessage sends a message to the recipients.
// If an error occurs, it is returned.
func (m Messenger) SendMessage(title string, message string, recipients []string) error {
	return m.send(title, message, recipients)
}
