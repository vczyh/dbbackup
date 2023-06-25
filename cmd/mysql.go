package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbbackup/backup"
	"github.com/vczyh/dbbackup/backup/mysqlbackup"
	"github.com/vczyh/dbbackup/log/zaplog"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql",

	RunE: func(cmd *cobra.Command, args []string) error {
		return mysqlCmdRun()
	},
}

var (
	mc = &mysqlConfig{}
)

type mysqlConfig struct {
	cnf                  string
	xtraBackup           bool
	mysqlDump            bool
	xtraBackupBinaryPath string
	socket               string
	user                 string
	password             string
	xtraBackupFlags      string
}

func init() {
	mysqlCmd.Flags().BoolVar(&mc.xtraBackup, "xtrabackup", false, "Backup by xtrabackup")
	mysqlCmd.Flags().BoolVar(&mc.mysqlDump, "mysqldump", false, "Backup by mysqldump")
	mysqlCmd.Flags().StringVar(&mc.xtraBackupBinaryPath, "xtrabackup-path", "", "The path of xtrabackup executable")
	mysqlCmd.Flags().StringVar(&mc.cnf, "cnf", "/etc/mysql/my.cnf", "The path of my.cnf")
	mysqlCmd.Flags().StringVar(&mc.socket, "socket", "/var/run/mysqld/mysqld.sock", "Unix socket path")
	mysqlCmd.Flags().StringVar(&mc.user, "user", "root", "MySQL user")
	mysqlCmd.Flags().StringVar(&mc.password, "password", "", "MySQL password")
	mysqlCmd.Flags().StringVar(&mc.xtraBackupFlags, "xtrabackup-flags", "", "XtraBackup extra flags")
}

func mysqlCmdRun() error {
	logger := zaplog.Default
	s3Client := GetS3Client()
	backupStorage, err := GetStorage(s3Client)
	if err != nil {
		return err
	}

	var manager *mysqlbackup.Manager
	if mc.xtraBackup {
		manager, err = mysqlbackup.New(
			mysqlbackup.WithLogger(logger),
			mysqlbackup.WithS3Client(s3Client),
			mysqlbackup.WithBackupStorage(backupStorage),
			mysqlbackup.WithXtraBackup(),
			mysqlbackup.WithCnf(mc.cnf),
			mysqlbackup.WithXtraBackupBinaryPath(mc.xtraBackupBinaryPath),
			mysqlbackup.WithSocket(mc.socket),
			mysqlbackup.WithUser(mc.user),
			mysqlbackup.WithPassword(mc.password),
		)
		if err != nil {
			return err
		}
	}
	if manager == nil {
		return fmt.Errorf("please specify the backup type")
	}

	notifiers, err := GetNotifiers()
	if err != nil {
		return err
	}

	if err := backup.Execute(context.Background(), logger, manager, notifiers); err != nil {
		return err
	}
	return nil
}
