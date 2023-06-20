## 功能

- 备份
- 存储

## 安装

从 [Releases](https://github.com/vczyh/dbbackup/releases) 下载，或者如果有 `go` 环境：

```
go install github.com/vczyh/dbbackup@latest
```

## 使用

通过 s3，使用 xtrabackup 备份 mysql 到 oss；

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

### Flag 说明

### 存储

| 名称                     | 默认 | 说明                         |
|------------------------|----|----------------------------|
| --storage              | "" | 存储类型，支持 `s3`               |
| --prefix               | "" | 路径前缀                       |
| --s3-access-key-id     | "" | s3 access key id           |
| --s3-secret-access-key | "" | s3 secret access key       |
| --s3-endpoint          | "" | s3 endpoint                |
| --s3-region            | "" | s3 region                  |
| --s3-bucket            | "" | s3 bucket                  |
| --s3-force-path-style  | "" | s3 enable force path style |

### MySQL

| 名称                | 默认                          | 说明                   |
|-------------------|-----------------------------|----------------------|
| --xtrabackup      | false                       | 使用xtrabackup备份       |
| --xtrabackup-path | ""                          | xtrabackup 可执行程序路径   |
| --cnf             | /etc/mysql/my.cnf           | mysql 配置文件路径         |
| --socket          | /var/run/mysqld/mysqld.sock | mysql unix socket 路径 |
| --user            | root                        | mysql 用户名            |
| --password        | ""                          | mysql 密码             |


