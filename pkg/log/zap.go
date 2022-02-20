package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(appName string, level zapcore.LevelEnabler) *zap.Logger {
	encConfig := zap.NewProductionEncoderConfig()
	encConfig.EncodeLevel = levelEncoder
	encConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	encoder := zapcore.NewJSONEncoder(encConfig)
	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= zapcore.ErrorLevel
				}),
			),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return zapcore.DebugLevel <= lvl && lvl < zapcore.ErrorLevel
				}),
			),
		),
		zap.AddCaller(),
		zap.Fields(zap.String("applicationName", appName)),
		zap.IncreaseLevel(level),
	)
	return logger
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("DEBUG")
	case zapcore.WarnLevel:
		enc.AppendString("DEBUG")
	case zapcore.ErrorLevel:
		enc.AppendString("DEBUG")
	case zapcore.PanicLevel, zapcore.FatalLevel, zapcore.DPanicLevel:
		enc.AppendString("DEBUG")
	}
}
