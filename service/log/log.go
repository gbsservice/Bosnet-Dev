package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gormlogger "gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	FormatJson    = "json"
	FormatConsole = "console"

	TimeKey       = "time"
	LevelKey      = "level"
	NameKey       = "logger"
	CallerKey     = "caller"
	MessageKey    = "msg"
	StackTraceKey = "stacktrace"

	MaxSize    = 1
	MaxBackups = 5
	MaxAge     = 7
)

type Logger struct {
	ZapLogger        *zap.Logger
	LogLevel         gormlogger.LogLevel
	SlowThreshold    time.Duration
	SkipCallerLookup bool
}

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)

func New(zapLogger *zap.Logger) Logger {
	return Logger{
		ZapLogger:        zapLogger,
		LogLevel:         gormlogger.Warn,
		SlowThreshold:    100 * time.Millisecond,
		SkipCallerLookup: false,
	}
}

func (l Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		ZapLogger:        l.ZapLogger,
		SlowThreshold:    l.SlowThreshold,
		LogLevel:         level,
		SkipCallerLookup: l.SkipCallerLookup,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger().Sugar().Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger().Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger().Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		sql, rows := fc()
		l.logger().Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.logger().Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		l.logger().Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

func (l Logger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}

func Config(logLevel zapcore.Level, logFormat, fileName string) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        TimeKey,
		LevelKey:       LevelKey,
		NameKey:        NameKey,
		CallerKey:      CallerKey,
		MessageKey:     MessageKey,
		StacktraceKey:  StackTraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	var encoder zapcore.Encoder
	switch logFormat {
	case FormatJson:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	hook := lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   true,
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)),
		zap.NewAtomicLevelAt(logLevel),
	)
	development := zap.Development()
	logger := zap.New(core, development)
	//caller := zap.AddCaller()
	//logger := zap.New(core, caller, development)
	zap.ReplaceGlobals(logger)
	return logger
}

/*
func Log() (Logger *logs.Logger) {
	e, err := os.OpenFile("./foo.logs", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	Logger = logs.New(e, "", logs.Ldate|logs.Ltime)
	Logger.SetOutput(&lumberjack.Logger{
		Filename:   "./foo.logs",
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
	})
	return Logger
}
*/
