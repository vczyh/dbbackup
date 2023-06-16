package s3storage

//func TestS3BackupStorage_ListBackups(t *testing.T) {
//	bs, err := New(&Config{
//		Endpoint:        "http://127.0.0.1:9000",
//		AccessKeyID:     "QTBELHBAPSf3un1m57mG",
//		SecretAccessKey: "EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
//		ForcePathStyle:  true,
//		Bucket:          "backup",
//		Region:          "test",
//		Prefix:          "",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	backups, err := bs.ListBackups(context.TODO(), "backup")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(backups)
//}
//
//func TestS3Backup_ReadFile(t *testing.T) {
//	bs, err := New(&Config{
//		Endpoint:        "http://127.0.0.1:9000",
//		AccessKeyID:     "QTBELHBAPSf3un1m57mG",
//		SecretAccessKey: "EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
//		ForcePathStyle:  true,
//		Bucket:          "backup",
//		Region:          "test",
//		Prefix:          "",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	backups, err := bs.ListBackups(context.TODO(), "backup")
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, backup := range backups {
//		if backup.Name() == "b1" {
//			file, err := backup.ReadFile(context.TODO(), "test")
//			if err != nil {
//				t.Fatal(err)
//			}
//			defer file.Close()
//			b, err := io.ReadAll(file)
//			if err != nil {
//				t.Fatal(err)
//			}
//			//t.Log(string(b))
//			fmt.Print(string(b))
//		}
//	}
//}
//
//func TestS3BackupStorage_RemoveBackup(t *testing.T) {
//	bs, err := New(&Config{
//		Endpoint:        "http://127.0.0.1:9000",
//		AccessKeyID:     "QTBELHBAPSf3un1m57mG",
//		SecretAccessKey: "EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
//		ForcePathStyle:  true,
//		Bucket:          "backup",
//		Region:          "test",
//		Prefix:          "",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if err = bs.RemoveBackup(context.TODO(), "backup", "11"); err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestS3BackupStorage_StartBackup(t *testing.T) {
//	bs, err := New(&Config{
//		Endpoint:        "http://127.0.0.1:9000",
//		AccessKeyID:     "QTBELHBAPSf3un1m57mG",
//		SecretAccessKey: "EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
//		ForcePathStyle:  true,
//		Bucket:          "backup",
//		Region:          "test",
//		Prefix:          "",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	name := fmt.Sprintf("%d", time.Now().Unix())
//	bh, err := bs.StartBackup(context.TODO(), "backup", name)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	file, err := bh.AddFile(context.TODO(), "file1", -1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer file.Close()
//	_, err = file.Write([]byte("abcdefg"))
//	if err != nil {
//		t.Fatal(err)
//	}
//	file.Close()
//
//	if err = bh.Wait(context.TODO()); err != nil {
//		t.Fatal(err)
//	}
//}
