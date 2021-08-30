package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

/*********************************************
 * Custom logger definition for Gorm to use Logrus
 * Implement gorm Logger iterface
 * Based on : https://github.com/go-gorm/gorm/blob/master/logger/logger.go
 *********************************************/

// GormLogger is a wrapper class that implement Gorm logger interface
type GormLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

var (
	// LogrusGormLogger default GormLogger instance for Gorm logging through Logrus
	LogrusGormLogger = GormLogger{
		LogLevel:      logger.Warn,
		SlowThreshold: 200 * time.Millisecond,
	}
)

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	l.LogLevel = level
	return &newlogger
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logrus.WithContext(ctx).Info(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
	}
}

// Warn print warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logrus.WithContext(ctx).Warn(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
	}
}

// Error print error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logrus.WithContext(ctx).Error(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
	}
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				logrus.WithContext(ctx).Error(utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logrus.WithContext(ctx).Error(utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				logrus.WithContext(ctx).Warn(utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logrus.WithContext(ctx).Warn(utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				logrus.WithContext(ctx).Debug(utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logrus.WithContext(ctx).Debug(utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}
