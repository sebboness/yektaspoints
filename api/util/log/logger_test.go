package log

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Logger_Get(t *testing.T) {
	logger := Get()
	assert.NotNil(t, logger)
}

func Test_Logger_NewLogger(t *testing.T) {
	NewLogger("test")
}

func Test_Logger_addLoggerFields(t *testing.T) {
	type state struct {
		fields map[string]any
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{map[string]any{"user_id": "123", "id": "456"}}, want{}},
	}

	for _, c := range cases {

		logger := NewLoggerWithContext("test", context.Background())

		// now add some fields
		logger.addLoggerFields(c.fields)
		loggerFields := logger.getLoggerFields()

		assert.GreaterOrEqual(t, len(loggerFields), 2)
		assert.Equal(t, "123", loggerFields["user_id"])
		assert.Equal(t, "456", loggerFields["id"])
	}
}

func Test_Logger_WithContext(t *testing.T) {
	logger := NewLoggerWithContext("test", context.Background())

	logger = logger.AddField("id", "xyz")

	logger = logger.WithContext(context.Background())
	loggerFields := logger.getLoggerFields()

	assert.GreaterOrEqual(t, len(loggerFields), 1)
	assert.Equal(t, "xyz", loggerFields["id"])
}

func Test_Logger_WithField(t *testing.T) {
	logger := NewLoggerWithContext("test", context.Background())

	logger = logger.AddField("id", "xyz")
	loggerFields := logger.getLoggerFields()

	assert.GreaterOrEqual(t, len(loggerFields), 1)
	assert.Equal(t, "xyz", loggerFields["id"])
}

func Test_Logger_WithFields(t *testing.T) {
	logger := NewLoggerWithContext("test", context.Background())

	logger = logger.AddFields(map[string]any{"user_id": "123", "id": "456"})
	loggerFields := logger.getLoggerFields()

	assert.GreaterOrEqual(t, len(loggerFields), 2)
	assert.Equal(t, "123", loggerFields["user_id"])
	assert.Equal(t, "456", loggerFields["id"])
}

func Test_Logger_Log(t *testing.T) {
	loggerFields := map[string]any{"user_id": "123"}
	logger := NewLoggerWithContext("test", context.Background()).WithLevel(logrus.DebugLevel).AddFields(loggerFields)

	logger.Debugf("debug message %s:%s", "a", "b")
	logger.Infof("info message %s:%s", "a", "b")
	logger.Warnf("warn message %s:%s", "a", "b")
	logger.AddField("test", 456).Errorf("error message %s:%s", "a", "b")

	loggerFields = logger.getLoggerFields()
	assert.GreaterOrEqual(t, len(loggerFields), 2)
	assert.Equal(t, "123", loggerFields["user_id"])
	assert.Equal(t, 456, loggerFields["test"])
}
