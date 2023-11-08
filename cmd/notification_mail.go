package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbbackup/notification"
	"github.com/vczyh/dbbackup/notification/mailnotification"
)

var (
	mailNotificationCmd = &cobra.Command{
		Use:   "mail",
		Short: "Send a test mail",

		RunE: func(cmd *cobra.Command, args []string) error {
			return sendMail()
		},
	}

	mnc = &struct {
		username string
		password string
		host     string
		port     int
		to       []string
	}{}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&mnc.username, "email-username", "", "Mail username")
	rootCmd.PersistentFlags().StringVar(&mnc.password, "email-password", "", "Mail password")
	rootCmd.PersistentFlags().StringVar(&mnc.host, "email-host", "", "Mail SMTP server host")
	rootCmd.PersistentFlags().IntVar(&mnc.port, "email-port", 25, "Mail SMTP server port")
	rootCmd.PersistentFlags().StringSliceVar(&mnc.to, "email-to", nil, "Mail recipients")
}

func GetMailNotifier() (notifier notification.Notifier, err error) {
	if mnc.username == "" {
		return nil, fmt.Errorf("mail flag required")
	}
	mailNotifier, err := mailnotification.New(
		mnc.username,
		mnc.password,
		mnc.host,
		mnc.to,
		mailnotification.WithPort(mnc.port),
	)
	if err != nil {
		return nil, err
	}
	return mailNotifier, nil

}

func sendMail() error {
	mailNotifier, err := GetMailNotifier()
	if err != nil {
		return err
	}
	return mailNotifier.Test()
}
