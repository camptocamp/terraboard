package db

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestLogMode(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Warn,
		SlowThreshold: 200 * time.Millisecond,
	}

	TestLogger.LogMode(logger.Silent)
	assert.Equal(t, TestLogger.LogLevel, logger.Silent)
	TestLogger.LogMode(logger.Info)
	assert.Equal(t, TestLogger.LogLevel, logger.Info)
	TestLogger.LogMode(logger.Warn)
	assert.Equal(t, TestLogger.LogLevel, logger.Warn)
	TestLogger.LogMode(logger.Error)
	assert.Equal(t, TestLogger.LogLevel, logger.Error)
}

func TestLog_Info(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)

	// Log something on Info level
	TestLogger.Info(context.Background(), "Test")

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
}

func TestLog_Warn(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)

	// Log something on Warn level
	TestLogger.Warn(context.Background(), "Test")

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
}

func TestLog_Error(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}

	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)

	// Log something on Error level
	TestLogger.Error(context.Background(), "Test")

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
}

func TestLog_Trace_Normal(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 0,
	}

	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)
	logrus.SetLevel(logrus.DebugLevel)

	// Log something on Trace level
	TestLogger.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "Test", 1
	}, nil)

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}

	buf.Reset()

	// Log something on Trace level without return rows
	TestLogger.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "Test", -1
	}, nil)

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
}

func TestLog_Trace_Slow(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond,
	}

	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)
	logrus.SetLevel(logrus.DebugLevel)

	// Log something on Trace level
	begin := time.Now()
	time.Sleep(200 * time.Millisecond)
	TestLogger.Trace(context.Background(), begin, func() (string, int64) {
		return "Test", 1
	}, nil)

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
	if !strings.Contains(buf.String(), "SLOW") {
		t.Error("Expected slow to be in the log")
	}

	buf.Reset()

	// Log something on Trace level without return rows
	begin = time.Now()
	time.Sleep(200 * time.Millisecond)
	TestLogger.Trace(context.Background(), begin, func() (string, int64) {
		return "Test", -1
	}, nil)

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
	if !strings.Contains(buf.String(), "SLOW") {
		t.Error("Expected slow to be in the log")
	}
}

func TestLog_Trace_Error(t *testing.T) {
	TestLogger := GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond,
	}

	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)
	logrus.SetLevel(logrus.DebugLevel)

	// Log something on Trace level
	TestLogger.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "Test", 1
	}, gorm.ErrRecordNotFound)

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
	if !strings.Contains(buf.String(), "record not found") {
		t.Error("Expected slow to be in the log")
	}

	buf.Reset()

	// Log something on Trace level
	TestLogger.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "Test", -1
	}, gorm.ErrRecordNotFound)

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}
	if !strings.Contains(buf.String(), "record not found") {
		t.Error("Expected slow to be in the log")
	}
}
