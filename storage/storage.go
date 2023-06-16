package storage

import (
	"context"
	"io"
)

type BackupStorage interface {
	ListBackups(ctx context.Context, dir string) ([]Backup, error)
	StartBackup(ctx context.Context, dir, name string) (BackupHandler, error)
	RemoveBackup(ctx context.Context, dir, name string) error
}

type BackupMeta interface {
	Directory() string
	Name() string
}

type Backup interface {
	BackupMeta
	ReadFile(ctx context.Context, filename string) (io.ReadCloser, error)
}

// BackupHandler represents the return value of BackupStorage.StartBackup.
type BackupHandler interface {
	BackupMeta
	AddFile(ctx context.Context, filename string, filesize int64) (io.WriteCloser, error)
	Wait(ctx context.Context) error
	AbortBackup(ctx context.Context) error
}
