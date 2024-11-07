package logger

import (
	"github.com/jmoiron/sqlx"
)

const createLogTable = `
CREATE TABLE IF NOT EXISTS public.logs (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    level VARCHAR(50),
    message TEXT
);`

const (
	INFO_LOG  string = "INFO"
	ERROR_LOG string = "ERROR"
)

type Logger struct {
	db *sqlx.DB
}

func NewLogger(db *sqlx.DB) (*Logger, error) {
	_, err := db.Exec(createLogTable)
	if err != nil {
		return nil, err
	}
	return &Logger{db}, nil
}

func (l *Logger) LoggerBasic(level string, messages string) error {
	_, err := l.db.Exec("INSERT INTO logs (level, message) VALUES ($1, $2)", level, messages)

	return err
}
