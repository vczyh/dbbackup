package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Default = New()

type Logger struct {
	zapLogger *zap.SugaredLogger
}

func New() *Logger {
	logger := new(Logger)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.DebugLevel,
	)
	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		//zap.AddStacktrace(zapcore.ErrorLevel),
	)
	logger.zapLogger = zapLogger.Sugar()

	return logger
}

func (l *Logger) Debugf(format string, args ...any) {
	l.zapLogger.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.zapLogger.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.zapLogger.Warnf(format, args...)

}

func (l *Logger) Errorf(format string, args ...any) {
	l.zapLogger.Errorf(format, args...)
}
