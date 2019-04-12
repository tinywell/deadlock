package zaplog

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogfile = "./logs/default.log"
)

var (
	loggerCache  []*zap.Logger
	parentLogger *zap.Logger
	initOnce     sync.Once
)

// LogConfig config for zap log
type LogConfig struct {
	Level        string
	FilePath     string
	IsProduction bool
}

// InitLog init log module
func InitLog(cfg LogConfig) {
	var level zapcore.Level
	switch strings.ToLower(cfg.Level) {
	case "info":
		level = zap.InfoLevel
	case "debug":
		level = zap.DebugLevel
	case "warn", "warning":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}

	lvl := zap.NewAtomicLevelAt(level)

	http.HandleFunc("/handle/level", lvl.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(":9090", nil); err != nil {
			panic(err)
		}
	}()

	var filePath string
	if len(strings.TrimSpace(cfg.FilePath)) == 0 {
		filePath = defaultLogfile
	} else {
		filePath = strings.TrimSpace(cfg.FilePath)
	}

	_, err := os.Create(filePath)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			dir := filepath.Dir(filePath)
			merr := os.MkdirAll(dir, 0766)
			if merr != nil {
				panic(merr)
			}
			_, ferr := os.Create(filePath)
			if ferr != nil {
				panic(ferr)
			}
		} else {
			panic(err)
		}
	}

	encoder := zapcore.EncoderConfig{
		TimeKey:        "Time",
		LevelKey:       "Level",
		NameKey:        "Logger",
		CallerKey:      "Caller",
		MessageKey:     "Msg",
		StacktraceKey:  "StackTrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	zapConfig := zap.Config{
		Level:            lvl,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    encoder,
		OutputPaths:      []string{"stderr", filePath},
		ErrorOutputPaths: []string{"stderr"},
	}
	parentLogger, err = zapConfig.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(parentLogger)
	loggerCache = append(loggerCache, parentLogger)
}

// MustGetLogger return a logger named by `name`
func MustGetLogger(name string) *zap.Logger {
	if parentLogger == nil {
		// InitLog(LogConfig{})
		// panic(fmt.Errorf("log need init"))
		fmt.Println("Get Logger ", name)
		return zap.L()
	}
	logger := zap.L().Named(name)
	loggerCache = append(loggerCache, logger)
	return logger
}

// Sync call all logger's method `Sync`
func Sync() {
	for _, logger := range loggerCache {
		logger.Sync()
	}
}
