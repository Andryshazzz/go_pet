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

// Logger wraps zap.Logger with file-based logging support.
// It writes logs simultaneously to stdout and a timestamped file.
//
// Usage:
//
//	log, err := logger.New(logger.NewConfigMust())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer log.Close()
//
//	log.Info("server started", zap.Int("port", 8080))
type Logger struct {
	*zap.Logger
	file *os.File
}

// FromContext extracts a Logger from the given context.
// Panics if no logger is found.
//
// Usage:
//
//	func handler(ctx context.Context) {
//	    log := logger.FromContext(ctx)
//	    log.Info("handling request")
//	}
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
//
// Usage:
//
//	cfg := Config{
//	    Level:  "debug",
//	    Folder: "/var/log/myapp",
//	}
//	log, err := logger.NewLogger(cfg)
//	if err != nil {
//	    // handle error
//	}
//	defer log.Close()
func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()

	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("Unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("Mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, fmt.Errorf("Open log file: %w", err)
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
// All log messages from the returned logger will include these fields.
//
// Usage:
//
//	requestLogger := log.With(
//	    zap.String("request_id", requestID),
//	    zap.String("method", "GET"),
//	)
//	requestLogger.Info("processing request")
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

// Close flushes and closes the underlying log file.
// Should be called once when the application shuts down.
//
// Usage:
//
//	defer log.Close()
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("Failed to close application logger:", err)
	}
}
