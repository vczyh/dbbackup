package cmd

import (
	"fmt"
	"github.com/vczyh/dbbackup/log/zaplog"
	"github.com/vczyh/dbbackup/storage"
	"github.com/vczyh/dbbackup/storage/s3storage"
)

const (
	StorageTypeS3 = "s3"
)

var (
	storageType string
	prefix      string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&storageType, "storage", "", "Backup storage type, support s3")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "Backup storage common path")
	//rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "Backup storage dir")
}

func GetStorage() (bs storage.BackupStorage, err error) {
	switch storageType {
	case StorageTypeS3:
		s3Client := GetS3Client()
		bs, err = s3storage.New(s3Client, &s3storage.Config{
			Logger: zaplog.Default,
			Bucket: s3Bucket,
			//Delimiter:  "/"
			Prefix: prefix,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported storage type")
	}

	return bs, nil
}
