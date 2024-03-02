package log

import (
	"context"
	"fmt"
	"os"

	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sirupsen/logrus"
)

var (
	ctxFieldsKey = "loggerFields"
	logger       *logrus.Logger
	instance     *Logger

	fieldAppName   = "app_name"
	fieldName      = "name"
	fieldRequestId = "request_id"
)

type Logger struct {
	logger *logrus.Logger
	ctx    context.Context
}

func NewLogger(name string) *Logger {
	return NewLoggerWithContext(name, context.Background())
}

func NewLoggerWithContext(name string, ctx context.Context) *Logger {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = os.Stdout

	instance = &Logger{logger: logger, ctx: ctx}
	instance.Infof("Initialized new logger")
	return instance.WithAppName(getAppName()).WithName(name)
}

func Get() *Logger {
	if instance == nil {
		instance = NewLogger(getAppName())
		instance.Infof("Got new instance of logger")
	}
	return instance
}

// WithContext sets the context of this logger.
// All logger fields will be copied
func (l *Logger) WithContext(ctx context.Context) *Logger {
	loggerFields := l.getLoggerFields()
	l.ctx = ctx
	l.addLoggerFields(loggerFields)
	return l
}

// WithName sets the name for this logger
func (l *Logger) WithAppName(value string) *Logger {
	return l.WithField(fieldAppName, value)
}

// WithName sets the name for this logger
func (l *Logger) WithName(value string) *Logger {
	return l.WithField(fieldName, value)
}

func (l *Logger) WithLevel(level logrus.Level) *Logger {
	l.logger.SetLevel(level)
	return l
}

func (l *Logger) WithField(key string, value any) *Logger {
	l.addLoggerFields(map[string]any{key: value})
	return l
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	l.addLoggerFields(fields)
	return l
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.WithFields(l.getLoggerFields()).Fatalf(format, args...)
}

func (l *Logger) addLoggerFields(fields map[string]any) {
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

func (l *Logger) getLoggerFields() map[string]any {
	if l.ctx != nil && l.ctx.Value(ctxFieldsKey) != nil {
		return l.ctx.Value(ctxFieldsKey).(map[string]any)
	}
	return map[string]any{}
}

func getAppName() string {
	return fmt.Sprintf("%s_%s", env.GetEnv("APPNAME"), env.GetEnv("ENV"))
}
