package backup

import (
	"context"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/notification"
)

type Executor interface {
	ExecuteBackup(ctx context.Context) error
}

func Execute(ctx context.Context, logger log.Logger, executor Executor, notifiers []notification.Notifier) error {
	var finalErr error
	backupNotification := &notification.BackupNotification{
		IsSucceed: true,
		Message:   "Backup successfully",
	}

	if err := executor.ExecuteBackup(ctx); err != nil {
		logger.Errorf("backup failed: %v", err)
		finalErr = err
		backupNotification.IsSucceed = false
		backupNotification.Message = err.Error()
	}

	for _, notifier := range notifiers {
		if err := notifier.BackupNotify(backupNotification); err != nil {
			logger.Errorf("fail send mail: %v", err)
			if finalErr != nil {
				finalErr = err
			}
		}
	}

	return finalErr
}
