package log

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var ctxFieldsKey = "loggerFields"

type LogrusLogger struct {
	logger *logrus.Logger
	ctx    context.Context
}

func NewLogger(loggerType string) *LogrusLogger {
	return NewLoggerWithContext(loggerType, nil)
}

func NewLoggerWithContext(loggerType string, ctx context.Context) *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = os.Stdout

	return &LogrusLogger{logger: logger, ctx: ctx}
}

func (l *LogrusLogger) WithContext(ctx context.Context) *LogrusLogger {
	loggerFields := l.getLoggerFields()
	l.ctx = ctx
	l.addLoggerFields(loggerFields)
	return l
}

func (l *LogrusLogger) WithLevel(level logrus.Level) *LogrusLogger {
	l.logger.SetLevel(level)
	return l
}

func (l *LogrusLogger) WithField(key string, value any) *LogrusLogger {
	l.addLoggerFields(map[string]any{key: value})
	return l
}

func (l *LogrusLogger) WithFields(fields map[string]any) *LogrusLogger {
	l.addLoggerFields(fields)
	return l
}

func (l *LogrusLogger) Debugf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Debugf(format, args...)
}

func (l *LogrusLogger) Infof(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Infof(format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Warnf(format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Errorf(format, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Fatalf(format, args...)
}

func (l *LogrusLogger) addLoggerFields(fields map[string]any) {
	if l.ctx != nil {
		if l.ctx.Value(ctxFieldsKey) == nil {
			l.ctx = context.WithValue(l.ctx, ctxFieldsKey, map[string]any{})
		}

		ctxFields := l.ctx.Value(ctxFieldsKey).(map[string]any)

		for k, v := range fields {
			ctxFields[k] = v
		}

		l.ctx = context.WithValue(l.ctx, ctxFieldsKey, ctxFields)
	}
}

func (l *LogrusLogger) getLoggerFields() map[string]any {
	if l.ctx != nil && l.ctx.Value(ctxFieldsKey) != nil {
		return l.ctx.Value(ctxFieldsKey).(map[string]any)
	}
	return map[string]any{}
}
