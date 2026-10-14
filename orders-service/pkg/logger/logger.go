package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the configuration for the Logger.
// Values are passed from the application configuration.
type Config struct {
	Level  string
	Folder string
}

// Logger wraps zap.Logger with file-based logging support.
// It writes logs simultaneously to stdout and a timestamped file.
//
// Usage:
//
//	log, err := logger.NewLogger(logger.Config{
//	    Level:  "debug",
//	    Folder: "./out/logs",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer log.Close()
type Logger struct {
	*zap.Logger
	file *os.File
}

// FromContext extracts a Logger from the given context.
// Returns a no-op logger if no logger is found.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		return &Logger{
			Logger: zap.NewNop(),
		}
	}
	return log
}

// NewLogger creates a new Logger that writes to both stdout and a log file.
// The log file is created in config.Folder with a UTC timestamp in its name.
func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	defer func() {
		if err != nil {
			logFile.Close()
		}
	}()

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15-04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// With creates a child logger with additional pre-set fields.
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

// Close flushes and closes the underlying log file.
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}
