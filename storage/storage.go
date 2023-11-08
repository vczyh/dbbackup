package storage

import (
	"context"
	"io"
)

const (
	BackupsManifestFilename = "backups-manifest.json"
	BackupManifestFilename  = "manifest.json"
)

type BackupStorage interface {
	ListBackups(ctx context.Context, dir string) ([]Backup, error)
	StartBackup(ctx context.Context, dir, name string) (BackupHandler, error)
	RemoveBackup(ctx context.Context, dir, name string) error
	//GetManifest(ctx context.Context, dir string) (BackupManifestHandler, error)
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
	WriteManifest(ctx context.Context, manifest any) error
	ReadManifest(ctx context.Context, manifest any) error
	Wait(ctx context.Context) error
	AbortBackup(ctx context.Context) error
}

//type BackupManifestHandler interface {
//	Directory() string
//	ReadManifest(ctx context.Context, manifest *BackupsManifest) error
//	WriteManifest(ctx context.Context, manifest *BackupsManifest) error
//}
//
//type BackupsManifest struct {
//	Backups []*BackupManifest
//}
//
//type BackupManifest struct {
//	SnapshotTime time.Time
//}
