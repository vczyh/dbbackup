package backup

import (
	"context"
	"github.com/vczyh/dbbackup/storage"
)

type Backup interface {
	ExecuteBackup(ctx context.Context, bh storage.BackupHandler) error
	//AbortBackup(ctx context.Context) error
}
