package mysqlbackup

import (
	"context"
	"github.com/vczyh/dbbackup/client/s3client"
	"github.com/vczyh/dbbackup/storage"
)

type Config struct {
	// storage
	s3Client    *s3client.Client
	storageType string
	bs          storage.BackupStorage
}

func UseXtraBackup(ctx context.Context, bh storage.BackupStorage) {

}
