## 功能

- [备份 MySQL](#mysql)
- [备份 Redis](#redis)
- [上传到 OSS](#存储)
- [邮件通知](#通知)

## 安装

从 [Releases](https://github.com/vczyh/dbbackup/releases) 下载，或者如果有 `go` 环境：

```
go install github.com/vczyh/dbbackup@latest
```

## 使用

通过 s3，使用 xtrabackup 备份 mysql 到 oss：

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

备份后发送邮件通知：

```
 /tmp/dbbackup  mysql \
    --storage s3 \
    --s3-access-key-id QTBELHBAPSf3un1m57mG \
    --s3-secret-access-key EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb \
    --s3-endpoint http://192.168.64.1:9000 \
    --s3-bucket backup \
    --s3-region test \
    --email-username xxx@163.com \
    --email-password xxxx \
    --email-host smtp.163.com \
    --email-port 465 \
    --email-to xxxxx \
    --xtrabackup \
    --user bkpuser \
    --password 123 \
    --socket /var/run/mysqld/mysqld.sock \
    --xtrabackup-path /home/ubuntu/xtrabackup/bin/xtrabackup \
    --cnf /etc/mysql/my.cnf
```

## 存储

目前支持 S3。

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

## 通知

目前支持 email。

| 名称               | 类型 | 默认 | 说明                                   |
|------------------|----|----|--------------------------------------|
| --email-username | 邮件 | "" | `SMTP`用户名                            |
| --email-password | 邮件 | "" | `SMTP`密码                             |
| --email-host     | 邮件 | "" | `SMTP`host                           |
| --email-port     | 邮件 | 25 | `SMTP`端口                             |
| --email-to       | 邮件 | "" | 接收人，多个接收人使用 `--email-to mail1,mail2` |

## Redis

Redis 支持本地备份（直接拷贝本地RDB文件），和远程备份（通过网络获取RDB文件）。

| 名称         | 默认          | 说明       |
|------------|-------------|----------|
| --user     | ""          | 用户       |
| --password | ""          | 密码       |
| --host     | "127.0.0.1" | host     |
| --port     | 6379        | 端口       |
| --remote   | false       | 是否使用远程备份 |

## MySQL

| 名称                | 默认                          | 说明                   |
|-------------------|-----------------------------|----------------------|
| --xtrabackup      | false                       | 使用xtrabackup备份       |
| --xtrabackup-path | ""                          | xtrabackup 可执行程序路径   |
| --cnf             | /etc/mysql/my.cnf           | mysql 配置文件路径         |
| --socket          | /var/run/mysqld/mysqld.sock | mysql unix socket 路径 |
| --user            | root                        | mysql 用户名            |
| --password        | ""                          | mysql 密码             |


