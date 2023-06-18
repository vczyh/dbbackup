package mysqlengine

import (
	"bufio"
	"context"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/storage"
	"io"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	xtraBackupBinaryName        = "xtrabackup"
	xtraBackupTempTargetDirName = "xtrabackup-target"
)

type BackupConfig struct {
	CnfPath              string
	XtraBackupBinaryPath string
	Socket               string
	User                 string
	Password             string
	XtraBackupFlags      []string
}

type BackupInfo struct {
	GTID         string
	FinishedTime time.Time
}

type XtraBackup struct{}

func (e *XtraBackup) Backup(ctx context.Context, logger log.Logger, bh storage.BackupHandler, config *BackupConfig) (*BackupInfo, error) {
	backupName := e.backupName()

	backupWriter, err := bh.AddFile(ctx, backupName, -1)
	if err != nil {
		return nil, err
	}
	defer func() {
		logger.Infof("closing backup file")
		if err := backupWriter.Close(); err != nil {
			logger.Errorf("fail close backup file")
		}
	}()

	xtraBackupExecutable := xtraBackupBinaryName
	if config.XtraBackupBinaryPath != "" {
		xtraBackupExecutable = config.XtraBackupBinaryPath
	}

	targetDir := path.Join(os.TempDir(), xtraBackupTempTargetDirName)
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		return nil, err
	}

	xtraBackupFlags := []string{
		"--defaults-file=" + config.CnfPath,
		"--backup",
		"--stream=xbstream",
		"--socket=" + config.Socket,
		"--user=" + config.User,
		"--password=" + config.Password,
		"--target-dir=" + targetDir,
	}
	if len(config.XtraBackupFlags) > 0 {
		xtraBackupFlags = append(xtraBackupFlags, config.XtraBackupFlags...)
	}

	cmd := exec.CommandContext(ctx, xtraBackupExecutable, xtraBackupFlags...)
	backStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	backStderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	backStderrDone := make(chan struct{})
	go func() {
		defer close(backStderrDone)

		scanner := bufio.NewScanner(backStderr)
		for scanner.Scan() {
			line := scanner.Text()
			logger.Infof("[XtraBackup] %s", line)
			// TODO 查找 binlog position
		}

		if err := scanner.Err(); err != nil {
			logger.Errorf("fail scan xtrabackup stderr: %v", err)
		}
	}()

	n, err := io.Copy(backupWriter, backStdout)
	if err != nil {
		return nil, err
	}
	logger.Infof("backup file size: %d", n)

	<-backStderrDone
	finishedTime := time.Now()
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return &BackupInfo{
		// TODO GTID
		GTID:         "",
		FinishedTime: finishedTime,
	}, nil
}

func (e *XtraBackup) backupName() string {
	return "backup.xbstream"
}
