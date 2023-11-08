package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbbackup/client/s3client"
	"github.com/vczyh/dbbackup/log/zaplog"
	"strings"
)

var (
	s3StorageCmd = &cobra.Command{
		Use:   "s3",
		Short: "Test s3",

		RunE: func(cmd *cobra.Command, args []string) error {
			return listBuckets()
		},
	}

	s3AccessKeyID     string
	s3SecretAccessKey string
	s3ForcePathStyle  bool
	s3Bucket          string
	s3Endpoint        string
	s3Region          string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&s3AccessKeyID, "s3-access-key-id", "", "S3 storage access key ID")
	rootCmd.PersistentFlags().StringVar(&s3SecretAccessKey, "s3-secret-access-key", "", "S3 storage secret access key")
	rootCmd.PersistentFlags().BoolVar(&s3ForcePathStyle, "s3-force-path-style", true, "Enable s3 storage force path style")
	rootCmd.PersistentFlags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket")
	rootCmd.PersistentFlags().StringVar(&s3Endpoint, "s3-endpoint", "", "S3 endpoint")
	rootCmd.PersistentFlags().StringVar(&s3Region, "s3-region", "", "S3 region")
}

func GetS3Client() *s3client.Client {
	return s3client.New(
		s3AccessKeyID,
		s3SecretAccessKey,
		s3Endpoint,
		s3client.WithForcePathStyle(s3ForcePathStyle),
		s3client.WithRegion(s3Region),
	)
}

func listBuckets() error {
	s3Client := GetS3Client()
	buckets, err := s3Client.ListBuckets(context.Background())
	if err != nil {
		return err
	}
	logger := zaplog.Default
	logger.Infof("buckets: %s", strings.Join(buckets, ", "))
	return nil
}
