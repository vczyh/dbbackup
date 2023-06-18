package mysql

import (
	"bufio"
	"context"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/storage"
	"io"
	"os"
	"os/exec"
	"path"
)

const (
	xtraBackupBinaryName        = "xtrabackup"
	xtraBackupTempTargetDirName = "xtrabackup-target"
)

type Config struct {
	Logger               log.Logger
	CnfPath              string
	XtraBackupBinaryPath string
	Socket               string
	User                 string
	Password             string
	XtraBackupFlags      []string
}

type XtraBackup struct {
	config     Config
	cancelFunc context.CancelFunc
}

func NewXtraBackupEngine(config *Config) (*XtraBackup, error) {
	e := new(XtraBackup)
	e.config = *config
	return e, nil
}

func (b *XtraBackup) ExecuteBackup(ctx context.Context, bh storage.BackupHandler) error {
	logger := b.config.Logger

	backupWriter, err := bh.AddFile(ctx, "backup.xbstream", -1)
	if err != nil {
		return err
	}
	defer func() {
		logger.Infof("closing backup file")
		if err := backupWriter.Close(); err != nil {
			logger.Errorf("fail close backup file")
		}
	}()

	xtraBackupExecutable := xtraBackupBinaryName
	if b.config.XtraBackupBinaryPath != "" {
		xtraBackupExecutable = b.config.XtraBackupBinaryPath
	}

	targetDir := path.Join(os.TempDir(), xtraBackupTempTargetDirName)
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		return err
	}

	xtraBackupFlags := []string{
		"--defaults-file=" + b.config.CnfPath,
		"--backup",
		"--stream=xbstream",
		"--socket=" + b.config.Socket,
		"--user=" + b.config.User,
		"--password=" + b.config.Password,
		"--target-dir=" + targetDir,
	}
	if len(b.config.XtraBackupFlags) > 0 {
		xtraBackupFlags = append(xtraBackupFlags, b.config.XtraBackupFlags...)
	}

	cmd := exec.CommandContext(ctx, xtraBackupExecutable, xtraBackupFlags...)
	backStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	backStderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
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
		return err
	}
	logger.Infof("backup file size: %d", n)

	<-backStderrDone
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
