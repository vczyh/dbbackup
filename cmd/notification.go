package cmd

import (
	"errors"
	"github.com/vczyh/dbbackup/notification"
)

var (
	NotifierFlagNotSetError = errors.New("notifier flag not set")
	notifierNames           []string
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&notifierNames, "notifiers", []string{}, "Notifiers: support email")
}

func GetNotifiers() (notifiers []notification.Notifier, err error) {
	if len(notifierNames) == 0 {
		return nil, NotifierFlagNotSetError
	}
	for _, notifierName := range notifierNames {
		switch notifierName {
		case "email":
			mailNotifier, err := GetMailNotifier()
			if err != nil {
				return nil, err
			}
			notifiers = append(notifiers, mailNotifier)
		}
	}
	return notifiers, nil
}
