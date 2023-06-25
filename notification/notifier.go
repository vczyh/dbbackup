package notification

type BackupNotification struct {
	IsSucceed bool
	Message   string
}

type Notifier interface {
	BackupNotify(notification *BackupNotification) error
}
