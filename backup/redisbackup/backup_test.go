package redisbackup

import (
	"context"
	"fmt"
	"github.com/vczyh/dbbackup/client/s3client"
	"github.com/vczyh/dbbackup/log/zaplog"
	"github.com/vczyh/dbbackup/storage/s3storage"
	"io"
	"os"
	"testing"
)

func TestB(t *testing.T) {
	f, err := os.Open("/Users/zhangyuheng/workspace/mine/container-images/redis/redis/data/source/dump.rdb")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	//buf := make([]byte, 20)
	//for {
	//	n, err := f.Read(buf)
	//	if err != nil {
	//		if err == io.EOF {
	//			fmt.Println("COPY EOF")
	//			break
	//		}
	//	}
	//	println(string(buf[:n]))
	//}

	r, w := io.Pipe()

	go func() {
		buf := make([]byte, 20)
		for {
			n, err := r.Read(buf)
			if err != nil {
				panic(err)
				//t.Fatal(err)
			}
			fmt.Println("===", string(buf[:n]))
		}
	}()

	n, err := io.Copy(w, f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("size:", n)

	//_, err = io.Copy(os.Stdout, r)
	//if err != nil {
	//	t.Fatal(err)
	//}
}

func Test2(t *testing.T) {
	ctx := context.TODO()
	sc := s3client.New(
		"QTBELHBAPSf3un1m57mG",
		"EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
		"http://192.168.64.1:9000",
		s3client.WithForcePathStyle(true),
		s3client.WithRegion("auto"),
	)

	bs, err := s3storage.New(sc, &s3storage.Config{
		Logger: zaplog.Default,
		Bucket: "backup",
		//Delimiter:  "/"
		Prefix: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open("/Users/zhangyuheng/workspace/mine/container-images/redis/redis/data/source/dump.rdb")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	bh, err := bs.StartBackup(ctx, "dir1", "t1")
	if err != nil {
		t.Fatal(err)
	}
	w, err := bh.AddFile(ctx, "dump.rdb", -1)
	if err != nil {
		t.Fatal(err)
	}
	n, err := io.Copy(w, f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Log("size:", n)

	if err = bh.Wait(ctx); err != nil {
		t.Fatal(err)
	}

	backups, err := bs.ListBackups(ctx, "t1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(backups)
}
