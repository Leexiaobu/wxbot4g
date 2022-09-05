package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
	"wxbot4g/config"
)

// Log 日志工具
var Log *zap.SugaredLogger

func init() {
	logConfig := config.Config.LogConfig
	// 配置 sugaredLogger
	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
	}
	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间戳的格式
	encoderConfig.EncodeTime = customTimeEncoder
	// 日志级别使用大写带颜色显示
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.FileName,   //日志文件存放目录
		MaxSize:    logConfig.MaxSize,    //文件大小限制,单位MB
		MaxBackups: logConfig.MaxBackups, //最大保留日志文件数量
		MaxAge:     logConfig.MaxAge,     //日志文件保留天数
		Compress:   logConfig.Compress,   //是否压缩处理
	})
	logOut := zapcore.NewMultiWriteSyncer(os.Stdout, fileWriteSyncer)
	var level zapcore.LevelEnabler
	level, err := zapcore.ParseLevel(logConfig.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	// 将日志级别设置为 INFO
	core := zapcore.NewCore(encoder, logOut, level)
	// 增加 caller 信息
	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
}
