## Features

- 备份
- 存储


## Install

从[]()下载，或者如果有 go 环境： 

```
go install github.com/vczyh/dbbackup@latest
```

### Usage

使用`XtraBackup`备份`mysql`到`OSS`：

```
 dbbackup mysql \
    --storage s3 \
    --s3-access-key-id QTBELHBAPSf3un1m57mG \
    --s3-secret-access-key EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb \
    --s3-endpoint http://192.168.64.1:9000 \
    --s3-bucket backup \
    --s3-region test \
    --xtrabackup \
    --user bkpuser \
    --password 123 \
    --socket /var/run/mysqld/mysqld.sock \
    --xtrabackup-path /user/bin/xtrabackup \
    --cnf /etc/mysql/my.cnf 
```