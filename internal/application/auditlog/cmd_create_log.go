package auditlog

// CreateLogCommand 创建审计日志命令
type CreateLogCommand struct {
	UserID     uint
	Username   string
	Action     string
	Resource   string
	ResourceID string
	IPAddress  string
	UserAgent  string
	Details    string
	Status     string
}
