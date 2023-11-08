package redisbackup

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/vczyh/dbbackup/log"
	"github.com/vczyh/dbbackup/storage"
	"github.com/vczyh/redis-lib/replica"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	backupName = "backup.rdb"
)

var (
	configNotFound = errors.New("config not found")
)

type Manager struct {
	config *Config
	client *redis.Client
	logger log.Logger
}

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	UseRemote bool
}

func New(logger log.Logger, config *Config) (*Manager, error) {
	return &Manager{
		logger: logger,
		config: config,
	}, nil
}

func (m *Manager) ExecuteBackup(ctx context.Context, bh storage.BackupHandler) error {
	backupWriter, err := bh.AddFile(ctx, backupName, -1)
	if err != nil {
		return err
	}
	defer backupWriter.Close()

	m.createClient()
	if m.config.UseRemote {
		return m.remoteBackup(backupWriter)
	}
	return m.localBackup(ctx, backupWriter)
}

func (m *Manager) localBackup(ctx context.Context, backupWriter io.WriteCloser) error {
	if err := m.client.Save(ctx).Err(); err != nil {
		return err
	}
	dataDir, err := m.configGetOne(ctx, "dir")
	if err != nil {
		return err
	}
	rdbFilename, err := m.configGetOne(ctx, "dbfilename")
	if err != nil {
		return err
	}
	rdbFile := path.Join(dataDir, rdbFilename)

	f, err := os.Open(rdbFile)
	if err != nil {
		return err
	}

	n, err := io.Copy(backupWriter, f)
	if err != nil {
		return err
	}
	m.logger.Infof("backup rdb file size: %d", n)

	return nil
}

func (m *Manager) remoteBackup(backupWriter io.WriteCloser) error {
	repl, err := replica.NewReplica(&replica.Config{
		MasterIP:              m.config.Host,
		MasterPort:            m.config.Port,
		MasterUser:            m.config.User,
		MasterPassword:        m.config.Password,
		RdbWriter:             backupWriter,
		ContinueAfterFullSync: false,
	})
	if err != nil {
		return err
	}
	return repl.SyncWithMaster()
}

func (m *Manager) createClient() {
	if m.client != nil {
		return
	}
	m.client = redis.NewClient(&redis.Options{
		Addr:                  net.JoinHostPort(m.config.Host, strconv.Itoa(m.config.Port)),
		Username:              m.config.User,
		Password:              m.config.Password,
		DialTimeout:           2 * time.Second,
		WriteTimeout:          2 * time.Second,
		ReadTimeout:           2 * time.Second,
		ContextTimeoutEnabled: true,
	})
}

func (m *Manager) configGetOne(ctx context.Context, name string) (string, error) {
	cm, err := m.client.ConfigGet(ctx, name).Result()
	if err != nil {
		return "", err
	}
	if v, ok := cm[name]; ok {
		return v, nil
	}
	return "", configNotFound
}
