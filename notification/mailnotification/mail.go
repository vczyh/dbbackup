package mailnotification

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/log/zaplog"
	"github.com/vczyh/dbbackup/notification"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

type Mail struct {
	logger   log.Logger
	username string
	password string
	host     string
	port     int
	to       []string
}

func New(username, password, host string, to []string, opts ...Option) (*Mail, error) {
	m := &Mail{
		logger:   zaplog.Default,
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

func (m *Mail) Test() error {
	var buf bytes.Buffer
	buf.WriteString("This email ensures that you can receive the mail.")
	buf.WriteString("<br>")
	buf.WriteString("https://github.com/vczyh/dbbackup")
	return m.send("DBBackup", buf.String())
}

func (m *Mail) BackupNotify(notification *notification.BackupNotification) error {
	m.logger.Infof("send backup result notification by email")
	var body string
	if notification.IsSucceed {
		body = "Backup successfully."
	} else {
		body = "Backup failed: " + notification.Message + "."
	}
	return m.send("Backup Result Notification", body)
}

func (m *Mail) send(subject, body string) error {
	from := m.username
	password := m.password
	toAddress := strings.Join(m.to, ";")
	smtpHost := m.host
	smtpPort := m.port
	addr := net.JoinHostPort(smtpHost, strconv.Itoa(smtpPort))

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", toAddress))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	buf.WriteString(body)
	message := buf.Bytes()

	auth := smtp.PlainAuth("", from, password, smtpHost)
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	if err = c.Auth(auth); err != nil {
		return err
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, to := range m.to {
		if err = c.Rcpt(to); err != nil {
			return err
		}
		m.logger.Infof("Send email to %s", to)
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
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
