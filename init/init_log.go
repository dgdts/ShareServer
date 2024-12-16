package init

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dgdts/UniversalServer/pkg/config"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(config *config.Log) {
	switch config.LogMode {
	case ConsoleLogMode:
		initConsoleLogger(config)
	case FileLogMode:
		initFileLogger(config)
	default:
		hlog.Errorf("invalid log mode:[%+v]", config.LogMode)
	}
}

func initConsoleLogger(config *config.Log) {
	var opts []hertzzap.Option
	opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())),
		hertzzap.WithZapOptions(zap.AddCaller(), zap.AddCallerSkip(3)))
	logger := hertzzap.NewLogger(opts...)
	fileWriter := io.MultiWriter(os.Stdout)
	logger.SetOutput(fileWriter)
	hlog.SetLevel(GetLogLevel(config.LogLevel))
	hlog.SetLogger(logger)
}

func initFileLogger(config *config.Log) {
	ioWriter := &lumberjack.Logger{
		Filename:   config.LogFileName,
		MaxSize:    config.LogMaxSize,
		MaxBackups: config.LogMaxBackups,
		MaxAge:     config.LogMaxAge,
	}

	var opts []hertzzap.Option
	var output zapcore.WriteSyncer
	if !strings.Contains(strings.ToLower(os.Getenv("GO_ENV")), "prod") {
		opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())))
		output = zapcore.AddSync(ioWriter)
	} else {
		opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())))
		output = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(ioWriter),
			FlushInterval: time.Minute,
		}
	}

	log := hertzzap.NewLogger(opts...)
	hlog.SetLogger(log)
	hlog.SetLevel(GetLogLevel(config.LogLevel))
	hlog.SetOutput(output)
}

func GetLogLevel(level string) hlog.Level {
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
