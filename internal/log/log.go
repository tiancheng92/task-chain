package log

import (
	"io"
	"time"

	"github.com/Yostardev/gf"
	"github.com/bytedance/sonic"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

func init() {
	level := zapcore.InfoLevel

	logConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level), // 日志级别
		Development:       false,                       // 开发模式，堆栈跟踪
		DisableStacktrace: true,                        // 关闭自动堆栈捕获
		Encoding:          "console",                   // 输出格式 console 或 json
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			CallerKey:     "caller",
			LineEnding:    zapcore.DefaultLineEnding,
			NewReflectedEncoder: func(writer io.Writer) zapcore.ReflectedEncoder {
				enc := sonic.ConfigDefault.NewEncoder(writer)
				enc.SetEscapeHTML(false)
				return enc
			},
			EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
				c := new(color.Color)
				switch level {
				case zapcore.InfoLevel:
					c = color.New(color.FgBlue, color.Bold)
				case zapcore.WarnLevel:
					c = color.New(color.FgYellow, color.Bold)
				case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
					c = color.New(color.FgRed, color.Bold)
				default:
					c = color.New(color.FgWhite, color.Bold)
				}
				c.EnableColor()
				encoder.AppendString(c.Sprintf("[%s]", level.CapitalString()))
			},
			EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(gf.StringJoin("[", t.Format(time.DateTime), "]"))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(gf.StringJoin("[", caller.TrimmedPath(), "]:"))
			},
			ConsoleSeparator: " ",
		}, // 编码器配置
		InitialFields:    nil,                // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"}, // 错误输出到指定文件
	}

	// 构建日志
	l, err := logConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger = l
	sugaredLogger = l.Sugar()
}

func DebugWithArg(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func InfoWithArg(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func WarnWithArg(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func ErrorWithArg(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func FatalWithArg(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Debug(args ...any) {
	sugaredLogger.Debug(args...)
}

func Info(args ...any) {
	sugaredLogger.Info(args...)
}

func Warn(args ...any) {
	sugaredLogger.Warn(args...)
}

func Error(args ...any) {
	sugaredLogger.Error(args...)
}

func Fatal(args ...any) {
	sugaredLogger.Fatal(args...)
}

func Debugf(template string, args ...any) {
	sugaredLogger.Debugf(template, args...)
}

func Infof(template string, args ...any) {
	sugaredLogger.Infof(template, args...)
}

func Warnf(template string, args ...any) {
	sugaredLogger.Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	sugaredLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...any) {
	sugaredLogger.Fatalf(template, args...)
}
