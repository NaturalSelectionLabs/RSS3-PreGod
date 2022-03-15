package logger

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var _ gormlogger.Interface = &Logger{}

type Logger struct {
	CoreLogger                *zap.SugaredLogger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func New() gormlogger.Interface {
	return &Logger{
		CoreLogger:                logger.Logger,
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

func (l *Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &Logger{
		CoreLogger:                l.CoreLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}

	l.logger().Infof(str, args...)
}

func (l *Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}

	l.logger().Warnf(str, args...)
}

func (l *Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}

	l.logger().Errorf(str, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error &&
		(!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.logger().Error("trace",
			zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.logger().Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		l.logger().Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

func (l *Logger) logger() *zap.SugaredLogger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)

		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.HasSuffix(file, "dblogger.go"):
		default:
			return l.CoreLogger.Desugar().WithOptions(zap.AddCallerSkip(i)).Sugar()
		}
	}

	return l.CoreLogger
}
