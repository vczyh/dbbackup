package backup

import (
	"context"
	"fmt"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/notification"
	"github.com/vczyh/dbbackup/storage"
	"time"
)

type Executor interface {
	ExecuteBackup(ctx context.Context, storage storage.BackupHandler) error
}

func Execute(ctx context.Context,
	logger log.Logger,
	executor Executor,
	bs storage.BackupStorage,
	notifiers []notification.Notifier) (finalErr error) {

	notice := &notification.BackupNotification{
		IsSucceed: true,
		Message:   "Backup successfully",
	}
	assignErr := func(err error) {
		notice.IsSucceed = false
		notice.Message = err.Error()
	}
	defer func() {
		for _, notifier := range notifiers {
			if err := notifier.BackupNotify(notice); err != nil {
				logger.Errorf("fail send mail: %v", err)
				finalErr = err
			}
		}
	}()

	name := fmt.Sprintf("%d", time.Now().Unix())
	bh, err := bs.StartBackup(ctx, "backup", name)
	if err != nil {
		return err
	}

	if err := executor.ExecuteBackup(ctx, bh); err != nil {
		logger.Errorf("backup failed: %v", err)
		assignErr(err)
		return err
	}
	if err := bh.Wait(ctx); err != nil {
		assignErr(err)
		return err
	}

	return finalErr
}
