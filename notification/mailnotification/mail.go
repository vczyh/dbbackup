package mailnotification

import (
	"bytes"
	"fmt"
	"github.com/vczyh/dbbackup/notification"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

type Mail struct {
	username string
	password string
	host     string
	port     int
	to       []string
}

func New(username, password, host string, to []string, opts ...Option) (*Mail, error) {
	m := &Mail{
		username: username,
		password: password,
		host:     host,
		port:     25,
		to:       to,
	}
	for _, opt := range opts {
		opt.apply(m)
	}

	return m, nil
}

func (m *Mail) BackupNotify(notification *notification.BackupNotification) error {
	from := m.username
	password := m.password
	toAddress := strings.Join(m.to, ";")
	smtpHost := m.host
	smtpPort := m.port

	var body string
	if notification.IsSucceed {
		body = "Backup successfully"
	} else {
		body = "Backup failed: " + notification.Message
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", toAddress))
	buf.WriteString(fmt.Sprintf("Subject: Backup Result Notification\r\n"))
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	buf.WriteString(body)
	message := buf.Bytes()

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err := smtp.SendMail(net.JoinHostPort(smtpHost, strconv.Itoa(smtpPort)), auth, from, m.to, message); err != nil {
		return err
	}
	return nil
}

type Option interface {
	apply(mail *Mail)
}

type optionFunc func(mail *Mail)

func (f optionFunc) apply(b *Mail) {
	f(b)
}

func WithPort(port int) Option {
	return Option(optionFunc(func(m *Mail) {
		m.port = port
	}))
}
