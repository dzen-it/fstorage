package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	c := newConfig(zapcore.InfoLevel)
	l, _ := c.Build(zap.AddCallerSkip(1))
	l = l.Named("FileStorage")
	logger = l.Sugar()
}

func newConfig(lvl zapcore.Level) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			TimeKey:        "T",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     customTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func EnableDebugMode() {
	c := newConfig(zapcore.DebugLevel)
	l, _ := c.Build(zap.AddCallerSkip(1))
	l = l.Named("FileStorage")
	logger = l.Sugar()
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}
