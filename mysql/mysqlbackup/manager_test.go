package mysqlbackup

import (
	"context"
	"github.com/vczyh/dbbackup/client/s3client"
	"github.com/vczyh/dbbackup/log/zaplog"
	"github.com/vczyh/dbbackup/storage/s3storage"
	"testing"
)

func TestManager_ExecuteBackup(t *testing.T) {
	s3Client := s3client.New(
		"QTBELHBAPSf3un1m57mG",
		"EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
		"https://192.168.64.1:9000",
		s3client.WithForcePathStyle(true),
		s3client.WithRegion("test"),
	)

	bs, err := s3storage.New(s3Client, &s3storage.Config{
		Logger: zaplog.Default,
		Bucket: "backup",
		//Delimiter:  "/"
		//Prefix: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	m, err := New(
		WithLogger(zaplog.Default),
		WithS3Client(s3Client),
		WithBackupStorage(bs),
		WithXtraBackup(),
		WithCnf("/etc/mysql/my.cnf"),
		WithXtraBackupBinaryPath("/home/ubuntu/xtrabackup/bin/xtrabackup"),
		WithSocket("/var/run/mysqld/mysqld.sock"),
		WithUser("bkpuser"),
		WithPassword("123"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = m.ExecuteBackup(context.TODO()); err != nil {
		t.Fatal(err)
	}

}
