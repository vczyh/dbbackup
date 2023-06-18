package mysqlbackup

import (
	"context"
	"fmt"
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

type Backup struct {
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

func New(opts ...Option) (*Backup, error) {
	b := new(Backup)
	for _, opt := range opts {
		opt.apply(b)
	}

	//  TODO validate

	return b, nil
}

func (b *Backup) ExecuteBackup(ctx context.Context) error {
	name := fmt.Sprintf("%d", time.Now().Unix())
	bh, err := b.bs.StartBackup(ctx, "backup", name)
	if err != nil {
		return err
	}

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
	b.logger.Infof("backup info: %v", backupInfo)

	if err = bh.Wait(ctx); err != nil {
		return err
	}
	return nil
}

type Option interface {
	apply(*Backup)
}

type optionFunc func(*Backup)

func (f optionFunc) apply(b *Backup) {
	f(b)
}

func WithLogger(logger log.Logger) Option {
	return optionFunc(func(b *Backup) {
		b.logger = logger
	})
}

func WithS3Client(s3Client *s3client.Client) Option {
	return optionFunc(func(b *Backup) {
		b.s3Client = s3Client
	})
}

func WithBackupStorage(bs storage.BackupStorage) Option {
	return optionFunc(func(b *Backup) {
		b.bs = bs
	})
}

func WithXtraBackup() Option {
	return optionFunc(func(b *Backup) {
		b.backType = backupTypeXtraBackup
	})
}

func WithCnf(cnf string) Option {
	return optionFunc(func(b *Backup) {
		b.cnfPath = cnf
	})
}

func WithXtraBackupBinaryPath(xtraBackupBinaryPath string) Option {
	return optionFunc(func(b *Backup) {
		b.xtraBackupBinaryPath = xtraBackupBinaryPath
	})
}
func WithSocket(socket string) Option {
	return optionFunc(func(b *Backup) {
		b.socket = socket
	})
}
func WithUser(user string) Option {
	return optionFunc(func(b *Backup) {
		b.user = user
	})
}

func WithPassword(password string) Option {
	return optionFunc(func(b *Backup) {
		b.password = password
	})
}
