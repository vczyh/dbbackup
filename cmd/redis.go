package cmd

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbbackup/backup"
	"github.com/vczyh/dbbackup/backup/redisbackup"
	"github.com/vczyh/dbbackup/log/zaplog"
)

var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Backup redis",

	RunE: func(cmd *cobra.Command, args []string) error {
		return redisCmdRun()
	},
}

var (
	rc = &redisConfig{}
)

type redisConfig struct {
	host      string
	port      int
	user      string
	password  string
	useRemote bool
}

func init() {
	redisCmd.Flags().StringVar(&rc.user, "user", "", "Redis user")
	redisCmd.Flags().StringVar(&rc.password, "password", "", "Redis password")
	redisCmd.Flags().StringVar(&rc.host, "host", "127.0.0.1", "Redis host")
	redisCmd.Flags().IntVar(&rc.port, "port", 6379, "Redis host")
	redisCmd.Flags().BoolVar(&rc.useRemote, "remote", false, "Remote backup by replica receiving rdb")
}

func redisCmdRun() error {
	logger := zaplog.Default

	backupStorage, err := GetStorage()
	if err != nil {
		return err
	}

	m, err := redisbackup.New(logger, &redisbackup.Config{
		Host:      rc.host,
		Port:      rc.port,
		User:      rc.user,
		Password:  rc.password,
		UseRemote: rc.useRemote,
	})
	if err != nil {
		return err
	}

	notifiers, err := GetNotifiers()
	if err != nil {
		if errors.Is(err, NotifierFlagNotSetError) {
			notifiers = nil
		} else {
			return err
		}
	}

	if err := backup.Execute(context.Background(), logger, m, backupStorage, notifiers); err != nil {
		return err
	}
	return nil
}
