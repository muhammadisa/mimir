package stdlogger

type LogKey string

const (
	LogError   LogKey = "error"
	LogInfo    LogKey = "info"
	LogWarn    LogKey = "warning"
	LogSuccess LogKey = "success"
	LogCaller  LogKey = "caller"
	LogTook    LogKey = "took"
	LogService LogKey = "service"
	LogTime    LogKey = "time"
	LogData    LogKey = "data"
)