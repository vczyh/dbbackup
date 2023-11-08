package notification

type BackupNotification struct {
	IsSucceed bool
	Message   string
}

type Notifier interface {
	Test() error
	BackupNotify(notification *BackupNotification) error
}
