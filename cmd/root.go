package cmd

import (
	"github.com/vczyh/dbbackup/log/zaplog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "dbbackup",
	Short:        "A database backup tool",
	SilenceUsage: true,
}

func Execute(aversion string) {
	version = aversion
	err := rootCmd.Execute()
	if err != nil {
		zaplog.Default.Errorf("%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(mysqlCmd)
	rootCmd.AddCommand(redisCmd)
	rootCmd.AddCommand(mailNotificationCmd)
	rootCmd.AddCommand(s3StorageCmd)
}
