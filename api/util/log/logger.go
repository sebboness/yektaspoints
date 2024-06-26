package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

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
	logger     *logrus.Logger
	ctx        context.Context
	tempFields map[string]any
}

func NewLogger(name string) *Logger {
	return NewLoggerWithContext(name, context.Background())
}

func NewLoggerWithContext(name string, ctx context.Context) *Logger {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = os.Stdout

	instance = &Logger{
		logger:     logger,
		ctx:        ctx,
		tempFields: map[string]any{},
	}

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
	return l.AddField(fieldAppName, value)
}

// WithName sets the name for this logger
func (l *Logger) WithName(value string) *Logger {
	return l.AddField(fieldName, value)
}

func (l *Logger) WithLevel(level logrus.Level) *Logger {
	l.logger.SetLevel(level)
	return l
}

// AddField adds the given key and value to the context to be output with the logger
func (l *Logger) AddField(key string, value any) *Logger {
	l.addLoggerFields(map[string]any{key: value})
	return l
}

// AddFields adds the given fields to the context to be output with the logger
func (l *Logger) AddFields(fields map[string]any) *Logger {
	l.addLoggerFields(fields)
	return l
}

var mapMutex = sync.Mutex{}

// WithField uses a temporary map of fields that is cleared after each log output
func (l *Logger) WithField(key string, value any) *Logger {
	mapMutex.Lock()
	l.tempFields[key] = value
	mapMutex.Unlock()
	return l
}

// WithFields uses a temporary map of fields that is cleared after each log output
func (l *Logger) WithFields(fields map[string]any) *Logger {
	mapMutex.Lock()
	for k, v := range fields {
		l.tempFields[k] = v
	}
	mapMutex.Unlock()
	return l
}

func (l *Logger) Debugf(format string, args ...any) {
	if file, lineNum, ok := getCaller(); ok {
		l.WithField("caller", fmt.Sprintf("%v:%v", file, lineNum))
	}
	l.logger.WithFields(l.getLoggerFields()).Debugf(format, args...)
	l.tempFields = map[string]any{}
}

func (l *Logger) Infof(format string, args ...any) {
	if file, lineNum, ok := getCaller(); ok {
		l.WithField("caller", fmt.Sprintf("%v:%v", file, lineNum))
	}
	l.logger.WithFields(l.getLoggerFields()).Infof(format, args...)
	l.tempFields = map[string]any{}
}

func (l *Logger) Warnf(format string, args ...any) {
	if file, lineNum, ok := getCaller(); ok {
		l.WithField("caller", fmt.Sprintf("%v:%v", file, lineNum))
	}
	l.logger.WithFields(l.getLoggerFields()).Warnf(format, args...)
	l.tempFields = map[string]any{}
}

func (l *Logger) Errorf(format string, args ...any) {
	if file, lineNum, ok := getCaller(); ok {
		l.WithField("caller", fmt.Sprintf("%v:%v", file, lineNum))
	}
	l.logger.WithFields(l.getLoggerFields()).Errorf(format, args...)
	l.tempFields = map[string]any{}
}

func (l *Logger) Fatalf(format string, args ...any) {
	if file, lineNum, ok := getCaller(); ok {
		l.WithField("caller", fmt.Sprintf("%v:%v", file, lineNum))
	}
	l.logger.WithFields(l.getLoggerFields()).Fatalf(format, args...)
	l.tempFields = map[string]any{}
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
		mapMutex.Lock()
		addedFields := l.ctx.Value(ctxFieldsKey).(map[string]any)
		if len(l.tempFields) > 0 {
			for k, v := range l.tempFields {
				addedFields[k] = v
			}
		}
		mapMutex.Unlock()

		return addedFields
	}
	return map[string]any{}
}

func getAppName() string {
	return fmt.Sprintf("%s_%s", env.GetEnv("APPNAME"), env.GetEnv("ENV"))
}

// gets caller of logger log func
func getCaller() (string, int, bool) {
	_, file, no, ok := runtime.Caller(2)
	if ok {
		fileParts := strings.Split(file, "/")
		return fileParts[len(fileParts)-1], no, true
	}
	return "", 0, false
}
