package main

import (
	"fmt"
	"github.com/vczyh/dbbackup/log/zaplog"
	"github.com/vczyh/dbbackup/storage"
	"github.com/vczyh/dbbackup/storage/s3storage"
)

var (
	prefix      string
	storageType string
	//dir         string

	s3AccessKeyID     string
	s3SecretAccessKey string
	s3ForcePathStyle  bool
	s3Bucket          string
	s3Endpoint        string
	s3Region          string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&storageType, "storage", "", "Backup storage type")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "Backup storage common path")
	//rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "Backup storage dir")
	rootCmd.PersistentFlags().StringVar(&s3AccessKeyID, "s3-access-key-id", "", "S3 storage access key ID")
	rootCmd.PersistentFlags().StringVar(&s3SecretAccessKey, "s3-secret-access-key", "", "S3 storage secret access key")
	rootCmd.PersistentFlags().BoolVar(&s3ForcePathStyle, "s3-force-path-style", true, "Enable s3 storage force path style")
	rootCmd.PersistentFlags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket")
	rootCmd.PersistentFlags().StringVar(&s3Endpoint, "s3-endpoint", "", "S3 endpoint")
	rootCmd.PersistentFlags().StringVar(&s3Region, "s3-region", "", "S3 region")
}

func getStorage() storage.BackupStorage {
	logger := zaplog.Default

	switch storageType {
	case "s3":
		bs, err := s3storage.New(&s3storage.Config{
			Logger:          zaplog.Default,
			Endpoint:        s3Endpoint,
			AccessKeyID:     s3AccessKeyID,
			SecretAccessKey: s3SecretAccessKey,
			ForcePathStyle:  s3ForcePathStyle,
			Bucket:          s3Bucket,
			Region:          s3Region,
			Prefix:          prefix,
		})
		if err != nil {
			exit(logger, err)
		}
		return bs
	default:
		exit(logger, fmt.Errorf("unsupported storage type"))
	}
	return nil
}
