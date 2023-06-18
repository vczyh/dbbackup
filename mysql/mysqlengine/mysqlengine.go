package mysqlengine

var (
	xtraBackupEngine = &XtraBackup{}
)

func GetXtraBackupEngine() *XtraBackup {
	return xtraBackupEngine
}
