package mysqlbackup

import (
	"context"
	"github.com/vczyh/dbbackup/client/s3client"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/mysql/mysqlengine"
	"github.com/vczyh/dbbackup/storage"
	"time"
)

type BackupType string

const (
	backupTypeXtraBackup BackupType = "xtraBackup"
)

type Manager struct {
	logger   log.Logger
	backType BackupType

	s3Client *s3client.Client
	bs       storage.BackupStorage

	cnfPath              string
	xtraBackupBinaryPath string
	socket               string
	user                 string
	password             string
	//xtraBackupFlags      []string
}

func New(opts ...Option) (*Manager, error) {
	b := new(Manager)
	for _, opt := range opts {
		opt.apply(b)
	}

	//  TODO validate

	return b, nil
}

func (b *Manager) ExecuteBackup(ctx context.Context, bh storage.BackupHandler) error {
	//name := fmt.Sprintf("%d", time.Now().Unix())
	//bh, err := b.bs.StartBackup(ctx, "backup", name)
	//if err != nil {
	//	return err
	//}

	backupTime := time.Now()
	xtraBackupEngine := mysqlengine.GetXtraBackupEngine()
	backupInfo, err := xtraBackupEngine.Backup(ctx, b.logger, bh, &mysqlengine.BackupConfig{
		CnfPath:              b.cnfPath,
		XtraBackupBinaryPath: b.xtraBackupBinaryPath,
		Socket:               b.socket,
		User:                 b.user,
		Password:             b.password,
		//XtraBackupFlags:      ,
	})
	if err != nil {
		return err
	}

	manifest := &Manifest{
		BackupTime:   backupTime,
		FinishedTime: backupInfo.FinishedTime,
		GTID:         backupInfo.GTID,
	}
	if err := b.writeManifest(ctx, bh, manifest); err != nil {
		return err
	}

	if err = bh.Wait(ctx); err != nil {
		return err
	}
	return nil
}

type Option interface {
	apply(*Manager)
}

type optionFunc func(*Manager)

func (f optionFunc) apply(b *Manager) {
	f(b)
}

func WithLogger(logger log.Logger) Option {
	return optionFunc(func(b *Manager) {
		b.logger = logger
	})
}

func WithS3Client(s3Client *s3client.Client) Option {
	return optionFunc(func(b *Manager) {
		b.s3Client = s3Client
	})
}

func WithBackupStorage(bs storage.BackupStorage) Option {
	return optionFunc(func(b *Manager) {
		b.bs = bs
	})
}

func WithXtraBackup() Option {
	return optionFunc(func(b *Manager) {
		b.backType = backupTypeXtraBackup
	})
}

func WithCnf(cnf string) Option {
	return optionFunc(func(b *Manager) {
		b.cnfPath = cnf
	})
}

func WithXtraBackupBinaryPath(xtraBackupBinaryPath string) Option {
	return optionFunc(func(b *Manager) {
		b.xtraBackupBinaryPath = xtraBackupBinaryPath
	})
}
func WithSocket(socket string) Option {
	return optionFunc(func(b *Manager) {
		b.socket = socket
	})
}
func WithUser(user string) Option {
	return optionFunc(func(b *Manager) {
		b.user = user
	})
}

func WithPassword(password string) Option {
	return optionFunc(func(b *Manager) {
		b.password = password
	})
}
