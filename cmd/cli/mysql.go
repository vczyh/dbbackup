package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbbackup/backup/mysql"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/log/zaplog"
	"os"
	"time"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql",

	Run: func(cmd *cobra.Command, args []string) {
		mysqlCmdRun()
	},
}

var (
	cnf                  string
	xtraBackup           bool
	XtraBackupBinaryPath string
	socket               string
	user                 string
	password             string
	xtraBackupFlags      string
)

func init() {
	mysqlCmd.Flags().BoolVar(&xtraBackup, "xtrabackup", false, "Backup by xtrabackup")
	mysqlCmd.Flags().StringVar(&XtraBackupBinaryPath, "xtrabackup-path", "", "The path of xtrabackup executable")
	mysqlCmd.Flags().StringVar(&cnf, "cnf", "/etc/mysql/my.cnf", "The path of my.cnf")
	mysqlCmd.Flags().StringVar(&socket, "socket", "/var/run/mysqld/mysqld.sock", "Unix socket path")
	mysqlCmd.Flags().StringVar(&user, "user", "root", "MySQL user")
	mysqlCmd.Flags().StringVar(&password, "password", "", "MySQL password")
	mysqlCmd.Flags().StringVar(&xtraBackupFlags, "xtrabackup-flags", "", "XtraBackup extra flags")
}

func mysqlCmdRun() {
	logger := zaplog.Default

	engine, err := mysql.NewXtraBackupEngine(&mysql.Config{
		Logger:               logger,
		CnfPath:              cnf,
		XtraBackupBinaryPath: XtraBackupBinaryPath,
		Socket:               socket,
		User:                 user,
		Password:             password,
		//XtraBackupFlags:      strings.Split(xtraBackupFlags, " "),
	})
	if err != nil {
		exit(logger, err)
	}

	bs := getStorage()

	name := fmt.Sprintf("%d", time.Now().Unix())
	bh, err := bs.StartBackup(context.TODO(), "backup", name)
	if err != nil {
		exit(logger, err)
	}

	if err := engine.ExecuteBackup(context.TODO(), bh); err != nil {
		logger.Errorf("fail backup %v", err)
		if err := bh.AbortBackup(context.TODO()); err != nil {
			exit(logger, err)
		}
	}
	if err = bh.Wait(context.TODO()); err != nil {
		exit(logger, err)
	}
	logger.Infof("backup success")
}

func exit(logger log.Logger, err error) {
	logger.Errorf("%v", err)
	os.Exit(1)
}
