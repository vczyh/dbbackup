package cmd

import (
	"github.com/vczyh/dbbackup/notification"
	"github.com/vczyh/dbbackup/notification/mailnotification"
)

var mnc = &mailNotificationConfig{}

type mailNotificationConfig struct {
	username string
	password string
	host     string
	port     int
	to       []string
}

func init() {
	rootCmd.PersistentFlags().StringVar(&mnc.username, "mail-username", "", "Mail username")
	rootCmd.PersistentFlags().StringVar(&mnc.password, "mail-password", "", "Mail password")
	rootCmd.PersistentFlags().StringVar(&mnc.host, "mail-host", "", "Mail SMTP server host")
	rootCmd.PersistentFlags().IntVar(&mnc.port, "mail-port", 25, "Mail SMTP server port")
	rootCmd.PersistentFlags().StringArrayVar(&mnc.to, "mail-to", nil, "Mail recipients")
}

func GetNotifiers() (notifiers []notification.Notifier, err error) {
	if mnc.username != "" {
		mailNotifier, err := mailnotification.New(
			mnc.username,
			mnc.password,
			mnc.host,
			mnc.to,
		)
		if err != nil {
			return nil, err
		}
		notifiers = append(notifiers, mailNotifier)
	}

	return notifiers, nil
}
