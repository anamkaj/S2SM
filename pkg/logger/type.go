package logger

type LoggerType interface {
	LoggerBasic(level string, messages string) error
}

