package events

import (
	"crypto/tls"
	"github.com/nats-io/nats.go"
	"gopkg.in/gomail.v2"
	"log"
)

type Event struct {
	ec *nats.EncodedConn
}

func New(ec *nats.EncodedConn) *Event {
	return &Event{ec: ec}
}

type message struct {
	Email   string
	Subject string
	Text    string
}

func (e *Event) CreateEmailCreatedListener() error {
	_, err := e.ec.Subscribe("email:created", func(msg *message) {
		log.Printf("email:created => %v", msg)
		m := gomail.NewMessage()
		m.SetHeader("From", "admin@pay9.co")
		m.SetHeader("To", msg.Email)
		m.SetHeader("Subject", msg.Subject)
		m.SetBody("text/html", msg.Text)
		d := gomail.NewDialer(
			"smtp.gmail.com",
			2525,
			"user",
			"pass",
		)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		err := d.DialAndSend(m)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
