package mysql

import (
	"context"
	"fmt"
	"github.com/vczyh/dbbackup/log/zaplog"
	"github.com/vczyh/dbbackup/storage/s3storage"
	"testing"
	"time"
)

func TestXtraBackup_ExecuteBackup(t *testing.T) {

	engine, err := NewXtraBackupEngine(&Config{
		Logger:               zaplog.Default,
		CnfPath:              "/etc/mysql/my.cnf",
		XtraBackupBinaryPath: "/tmp/xtrabackup/bin/xtrabackup",
		Socket:               "/var/run/mysqld/mysqld.sock",
		User:                 "bkpuser",
		Password:             "123",
		XtraBackupFlags:      nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	bs, err := s3storage.New(&s3storage.Config{
		Logger:          zaplog.Default,
		Endpoint:        "http://192.168.64.1:9000",
		AccessKeyID:     "QTBELHBAPSf3un1m57mG",
		SecretAccessKey: "EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
		ForcePathStyle:  true,
		Bucket:          "backup",
		Region:          "test",
		Prefix:          "",
	})
	if err != nil {
		t.Fatal(err)
	}

	name := fmt.Sprintf("%d", time.Now().Unix())
	bh, err := bs.StartBackup(context.TODO(), "backup", name)
	if err != nil {
		t.Fatal(err)
	}
	time.AfterFunc(time.Second*1, func() {
		if err = bh.AbortBackup(context.TODO()); err != nil {
			t.Fatal(err)
		}
	})

	if err := engine.ExecuteBackup(context.TODO(), bh); err != nil {
		t.Fatal(err)
	}

	//if err := bh.Wait(context.TODO()); err != nil {
	//	t.Fatal(err)
	//}
	//if err = bh.AbortBackup(context.TODO()); err != nil {
	//	t.Fatal(err)
	//}

}
