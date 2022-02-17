package engine

import (
	"log/syslog"
	"net/url"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sysLogSink struct {
	sysLogWriter *syslog.Writer
}

func (s sysLogSink) Sync() error {
	return nil
}

func (s sysLogSink) Close() error {
	return nil
}

func (s sysLogSink) Write(p []byte) (n int, err error) {
	s.sysLogWriter.Write(p)

	return len(p), nil
}

type ZapLogger struct {
	zap.SugaredLogger
}

func InitZapLogger(loggerConfig config.LoggerStruct) (*zap.SugaredLogger, error) {
	outputPaths, err := parseOutputPaths(loggerConfig)
	if err != nil {
		return nil, err
	}

	if len(outputPaths) == 0 {
		outputPaths = append(outputPaths, "stdout")
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(getLogLevel(loggerConfig.Level)),
		Development:      true,
		Encoding:         loggerConfig.Encoding,
		EncoderConfig:    getDefaultEncoderCfg(),
		InitialFields:    map[string]interface{}{"prefix_tag": loggerConfig.PrefixTag},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{},
	}

	zaplog, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	sugaredLogger := zaplog.Sugar()

	return sugaredLogger, nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.DebugLevel
	}
}

func getDefaultEncoderCfg() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func parseOutputPaths(loggerConfig config.LoggerStruct) ([]string, error) {
	var outputPaths []string

	for _, outputConfig := range loggerConfig.Output {
		if outputConfig.Type == "stdout" ||
			outputConfig.Type == "stderr" {
			outputPaths = append(outputPaths, outputConfig.Type)
		} else if outputConfig.Type == "file" {
			outputPaths = append(outputPaths, outputConfig.Filepath)
		} else if outputConfig.Type == "syslog" {
			outputPaths = append(outputPaths, "Syslog://127.0.0.1")
			sysLogFactory := func(u *url.URL) (zap.Sink, error) {
				w, err := util.GetSysLogger(outputConfig, loggerConfig.PrefixTag)

				if err != nil {
					return nil, err
				}

				s := sysLogSink{
					sysLogWriter: w,
				}

				return s, nil
			}
			if err := zap.RegisterSink("Syslog", sysLogFactory); err != nil {
				return nil, err
			}
		}
	}

	return outputPaths, nil
}
