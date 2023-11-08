package mysqlbackup

import (
	"context"
	"github.com/vczyh/dbbackup/storage"
	"time"
)

const (
	backupManifestFileName = "MANIFEST.json"
)

// Manifest describes the metadata of one backup.
type Manifest struct {
	BackupTime   time.Time
	FinishedTime time.Time
	GTID         string
}

// BackupsManifest describe the metadata of all backups.
//type BackupsManifest struct {
//	backups []*Manifest
//}

func (m *Manager) writeManifest(ctx context.Context, bh storage.BackupHandler, manifest *Manifest) error {
	// Write current backup metadata to current backup dir.
	//file, err := bh.AddFile(ctx, backupManifestFileName, -1)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//b, err := json.MarshalIndent(manifest, "", "  ")
	//if err != nil {
	//	return err
	//}
	//if _, err := file.Write(b); err != nil {
	//	return err
	//}
	if err := bh.WriteManifest(ctx, manifest); err != nil {
		return err
	}

	// Write current backup metadata to backups manifest.
	//mh, err := m.bs.GetManifest(ctx, bh.Directory())
	//if err != nil {
	//	return err
	//}
	//var backupsManifest storage.BackupsManifest
	//
	//if err = mh.ReadManifest(ctx, &backupsManifest); err != nil {
	//	return err
	//}
	//backupsManifest.Backups = append(backupsManifest.Backups, &storage.BackupManifest{
	//	SnapshotTime: manifest.FinishedTime,
	//})

	//if err = mh.WriteManifest(ctx, &backupsManifest); err != nil {
	//	return err
	//}
	return nil
}
