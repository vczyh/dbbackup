package storage

import (
	"context"
	"io"
	"time"
)

const (
	BackupsManifestFilename = "backups-manifest.json"
)

type BackupStorage interface {
	ListBackups(ctx context.Context, dir string) ([]Backup, error)
	StartBackup(ctx context.Context, dir, name string) (BackupHandler, error)
	RemoveBackup(ctx context.Context, dir, name string) error
	GetManifest(ctx context.Context, dir string) (BackupManifestHandler, error)
}

type Backup interface {
	Directory() string
	Name() string
	ReadFile(ctx context.Context, filename string) (io.ReadCloser, error)
}

type BackupHandler interface {
	Directory() string
	Name() string
	AddFile(ctx context.Context, filename string, filesize int64) (io.WriteCloser, error)
	Wait(ctx context.Context) error
	AbortBackup(ctx context.Context) error
}

type BackupManifestHandler interface {
	Directory() string
	UnmarshalManifest(ctx context.Context, manifest *BackupsManifest) error
	MarshalManifest(ctx context.Context, manifest *BackupsManifest) error
}

type BackupsManifest struct {
	Backups []*BackupManifest
}

type BackupManifest struct {
	SnapshotTime time.Time
	Meta         map[string]any
}
