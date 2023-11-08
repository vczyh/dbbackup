package s3storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vczyh/dbbackup/client/s3client"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/storage"
	"io"
	"strings"
	"sync"
)

const (
	DefaultDelimiter = "/"
)

type Config struct {
	Logger    log.Logger
	Bucket    string
	Delimiter string
	Prefix    string
}

type S3Storage struct {
	config Config
	c      *s3client.Client
	mu     sync.Mutex
}

func New(client *s3client.Client, config *Config) (storage.BackupStorage, error) {
	bs := new(S3Storage)
	bs.config = *config
	bs.c = client

	if bs.config.Delimiter == "" {
		bs.config.Delimiter = DefaultDelimiter
	}

	return bs, nil
}

func (bs *S3Storage) ListBackups(ctx context.Context, dir string) ([]storage.Backup, error) {
	bs.config.Logger.Infof("list backups in %s", dir)
	c, err := bs.client()
	if err != nil {
		return nil, err
	}

	searchPrefix := objectKey(bs.config.Prefix, dir)
	if searchPrefix != "" {
		searchPrefix += "/"
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    &bs.config.Bucket,
		Prefix:    &searchPrefix,
		Delimiter: &bs.config.Delimiter,
	}

	objects, err := c.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	var subDirs []string
	for _, prefix := range objects.CommonPrefixes {
		subDir := strings.TrimPrefix(*prefix.Prefix, searchPrefix)
		subDir = strings.TrimSuffix(subDir, bs.config.Delimiter)
		subDirs = append(subDirs, subDir)
	}

	backups := make([]storage.Backup, 0, len(subDirs))
	for _, subDir := range subDirs {
		backups = append(backups, &S3Backup{
			logger: bs.config.Logger,
			client: c,
			bs:     bs,
			dir:    dir,
			name:   subDir,
		})
	}
	return backups, nil
}

func (bs *S3Storage) StartBackup(ctx context.Context, dir, name string) (storage.BackupHandler, error) {
	bs.config.Logger.Infof("start backup to %s/%s", dir, name)
	c, err := bs.client()
	if err != nil {
		return nil, err
	}
	return &S3BackupHandler{
		logger: bs.config.Logger,
		client: c,
		bs:     bs,
		dir:    dir,
		name:   name,
	}, nil
}

func (bs *S3Storage) RemoveBackup(ctx context.Context, dir, name string) error {
	bs.config.Logger.Infof("remove backup %s/%s", dir, name)
	c, err := bs.client()
	if err != nil {
		return err
	}

	key := objectKey(bs.config.Delimiter, bs.config.Prefix, dir, name)
	iterator := s3manager.NewDeleteListIterator(c, &s3.ListObjectsInput{
		Bucket: &bs.config.Bucket,
		Prefix: &key,
	})
	if err := s3manager.NewBatchDeleteWithClient(c).Delete(ctx, iterator); err != nil {
		return err
	}
	return nil
}

//func (bs *S3Storage) GetManifest(ctx context.Context, dir string) (storage.BackupManifestHandler, error) {
//	client, err := bs.client()
//	if err != nil {
//		return nil, err
//	}
//	return &S3ManifestHandler{
//		dir:    dir,
//		bs:     bs,
//		client: client,
//	}, nil
//}

func (bs *S3Storage) client() (*s3.S3, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	//if bs._client != nil {
	//	return bs._client, nil
	//}
	return bs.c.AwsClient()

	//// TODO tls self sign certification ignore?
	//
	//awsConfig := &aws.Config{
	//	// TODO  log level
	//	//LogLevel: aws.LogOff,
	//	Region:           &bs.config.Region,
	//	Endpoint:         &bs.config.Endpoint,
	//	S3ForcePathStyle: &bs.config.ForcePathStyle,
	//}
	//accessKeyID := bs.config.AccessKeyID
	//secretAccessKey := bs.config.SecretAccessKey
	//if accessKeyID != "" && secretAccessKey != "" {
	//	staticCredentials := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	//	awsConfig.Credentials = staticCredentials
	//}
	//
	//awsSession, err := session.NewSession(awsConfig)
	//if err != nil {
	//	return nil, err
	//}
	//bs._client = s3.New(awsSession)
	//
	//return bs._client, nil
}

type S3Backup struct {
	logger log.Logger
	client *s3.S3
	bs     *S3Storage
	dir    string
	name   string
}

func (b *S3Backup) Directory() string {
	return b.dir
}

func (b *S3Backup) Name() string {
	return b.name
}

func (b *S3Backup) ReadFile(ctx context.Context, filename string) (io.ReadCloser, error) {
	config := b.bs.config
	object := objectKey(config.Delimiter, config.Prefix, b.dir, b.name, filename)
	out, err := b.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &config.Bucket,
		Key:    &object,
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

type S3BackupHandler struct {
	logger log.Logger
	client *s3.S3
	bs     *S3Storage
	dir    string
	name   string

	wg      sync.WaitGroup
	mu      sync.Mutex
	cancels []context.CancelFunc
	errs    []error
}

func (bh *S3BackupHandler) Directory() string {
	return bh.dir
}

func (bh *S3BackupHandler) Name() string {
	return bh.name
}

func (bh *S3BackupHandler) AddFile(ctx context.Context, filename string, filesize int64) (io.WriteCloser, error) {
	ctx, cancelFunc := context.WithCancel(ctx)

	bh.mu.Lock()
	bh.cancels = append(bh.cancels, cancelFunc)
	bh.mu.Unlock()

	// TODO 分片

	reader, writer := io.Pipe()
	bh.wg.Add(1)
	go func() {
		defer bh.wg.Done()

		config := bh.bs.config
		key := objectKey(config.Delimiter, config.Prefix, bh.dir, bh.name, filename)

		_, err := s3manager.NewUploaderWithClient(bh.client).UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: &config.Bucket,
			Key:    &key,
			Body:   reader,
		})
		if err != nil {
			bh.logger.Errorf("fail upload file %s/%s: %v", bh.dir, bh.name, err)
			_ = reader.CloseWithError(err)
			bh.mu.Lock()
			bh.errs = append(bh.errs, err)
			bh.mu.Unlock()
		}
	}()

	return writer, nil
}

func (bh *S3BackupHandler) WriteManifest(ctx context.Context, manifest any) error {
	config := bh.bs.config
	object := objectKey(config.Delimiter, config.Prefix, bh.dir, bh.name, storage.BackupManifestFilename)
	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	_, err = s3manager.NewUploaderWithClient(bh.client).UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: &config.Bucket,
		Key:    &object,
		Body:   bytes.NewReader(b),
	})
	return err
}

func (bh *S3BackupHandler) ReadManifest(ctx context.Context, manifest any) error {
	config := bh.bs.config
	object := objectKey(config.Delimiter, config.Prefix, bh.dir, bh.name, storage.BackupManifestFilename)

	out, err := bh.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &config.Bucket,
		Key:    &object,
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()
	return json.NewDecoder(out.Body).Decode(manifest)
}

func (bh *S3BackupHandler) Wait(ctx context.Context) error {
	bh.wg.Wait()

	bh.mu.Lock()
	defer bh.mu.Unlock()

	if len(bh.errs) == 0 {
		return nil
	}
	var errs []string
	for _, err := range bh.errs {
		errs = append(errs, err.Error())
	}
	return errors.New(strings.Join(errs, ";"))
}

func (bh *S3BackupHandler) AbortBackup(ctx context.Context) error {
	bh.mu.Lock()
	for _, cancel := range bh.cancels {
		cancel()
	}
	bh.mu.Unlock()
	bh.wg.Wait()
	return bh.bs.RemoveBackup(ctx, bh.dir, bh.name)
}

type S3ManifestHandler struct {
	dir    string
	name   string
	bs     *S3Storage
	client *s3.S3
}

func (smh *S3ManifestHandler) Directory() string {
	return smh.dir
}

func (smh *S3ManifestHandler) Name() string {
	return smh.name
}

//func (smh *S3ManifestHandler) ReadManifest(ctx context.Context, manifest *storage.BackupsManifest) error {
//	config := smh.bs.config
//	object := objectKey(config.Delimiter, config.Prefix, smh.dir, storage.BackupsManifestFilename)
//
//	_, err := smh.client.HeadObject(&s3.HeadObjectInput{
//		Bucket: &config.Bucket,
//		Key:    &object,
//	})
//	if err != nil {
//		if awsErr, ok := err.(awserr.Error); ok {
//			if awsErr.Code() == "NotFound" {
//				return nil
//			}
//		}
//		return err
//	}
//
//	out, err := smh.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
//		Bucket: &config.Bucket,
//		Key:    &object,
//	})
//	if err != nil {
//		return err
//	}
//	defer out.Body.Close()
//	return json.NewDecoder(out.Body).Decode(manifest)
//}

//func (smh *S3ManifestHandler) WriteManifest(ctx context.Context, manifest *storage.BackupsManifest) error {
//	config := smh.bs.config
//	object := objectKey(config.Delimiter, config.Prefix, smh.dir, storage.BackupsManifestFilename)
//
//	b, err := json.MarshalIndent(manifest, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	_, err = s3manager.NewUploaderWithClient(smh.client).UploadWithContext(ctx, &s3manager.UploadInput{
//		Bucket: &config.Bucket,
//		Key:    &object,
//		Body:   bytes.NewReader(b),
//	})
//	return err
//}

func objectKey(delimiter string, parts ...string) string {
	res := strings.Join(parts, delimiter)
	res = strings.TrimPrefix(res, "/")
	res = strings.TrimSuffix(res, "/")
	return res
}
