package loggerfx

import (
	"go.uber.org/zap"
)

type Logger struct {
	zap.SugaredLogger
	s Storager
}

func New() *Logger {
	logger, _ := zap.NewProduction()
	slogger := logger.Sugar()

	return &Logger{*slogger, nil}
}

type Storager interface {
	WriteLog(msg []byte) error
}

func (l *Logger) SetStorager(s Storager) {
	l.s = s
}

func (l *Logger) InfoS(message []byte) {
	go func() {
		if l.s == nil {
			l.With(zap.String("Message", string(message))).
				Warn("Attempted to store log in nil storage")
			return
		}

		if err := l.s.WriteLog(message); err != nil {
			l.With(zap.String("Message", string(message))).
				Warn("Failed to write log to storage")
		}
	}()
}
