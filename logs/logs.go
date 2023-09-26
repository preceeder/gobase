package logs

import (
	"fmt"
	"github.com/preceeder/gobase/env"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
)

type LogConfig struct {
	DebugFileName string `json:"debugFileName"`
	InfoFileName  string `json:"infoFileName"`
	WarnFileName  string `json:"warnFileName"`
	MaxSize       int    `json:"maxsize"`
	MaxAge        int    `json:"max_age"`
	MaxBackups    int    `json:"max_backups"`
	StdOut        string `json:"stdOut"`
}

//func init() {
//	config := *ConfigObj.viper
//	logConfig := ReadLogConfig(config)
//	initSlog(logConfig)
//}

func InitLogWithViper(config viper.Viper) {
	logConfig := ReadLogConfig(config)
	initSlog(logConfig)
}

func ReadLogConfig(config viper.Viper) LogConfig {
	LCf := config.Sub("log")
	if LCf == nil {
		fmt.Printf("log config is nil")
		os.Exit(1)
	}
	//从配置中读取日志配置，初始化日志
	return LogConfig{
		DebugFileName: LCf.GetString("DebugFileName"),
		InfoFileName:  LCf.GetString("InfoFileName"),
		WarnFileName:  LCf.GetString("WarnFileName"),
		MaxSize:       LCf.GetInt("MaxSize"),
		MaxAge:        LCf.GetInt("MaxAge"),
		MaxBackups:    LCf.GetInt("MaxBackups"),
		StdOut:        LCf.GetString("StdOut"),
	}
}

func InitLogWithStruct(cfg LogConfig) {
	initSlog(cfg)
}

func initSlog(cfg LogConfig) {
	lumberJackInfo := &lumberjack.Logger{
		Filename:   cfg.InfoFileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	}
	opt := &slog.HandlerOptions{Level: slog.LevelInfo}
	var writeBuild io.Writer
	if cfg.StdOut == "1" {
		writeBuild = io.MultiWriter(os.Stdout, lumberJackInfo)
	} else {
		writeBuild = io.MultiWriter(lumberJackInfo)
	}
	var log *slog.Logger
	if env.GetEnv("env") == "product" || env.GetEnv("env") == "test" {
		log = slog.New(slog.NewJSONHandler(writeBuild, opt))
		slog.SetDefault(log)
	} else if env.GetEnv("env") == "local" {
		//log = slog.New(slog.NewTextHandler(writeBuild, nil))
		//slog.SetDefault(log)

	}
}

//
//// InitLogger 初始化Logger
//func initLogger(cfg LogConfig) (err error) {
//	writeSyncerDebug := getLogWriter(cfg.DebugFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
//	writeSyncerInfo := getLogWriter(cfg.InfoFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
//	writeSyncerWarn := getLogWriter(cfg.WarnFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
//	encoder := getEncoder()
//	//文件输出
//	debugCore := zapcore.NewCore(encoder, writeSyncerDebug, zapcore.DebugLevel)
//	infoCore := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel)
//	warnCore := zapcore.NewCore(encoder, writeSyncerWarn, zapcore.WarnLevel)
//	//标准输出
//	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
//	std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
//	core := zapcore.NewTee(debugCore, infoCore, warnCore, std)
//	LG := zap.New(core, zap.AddCaller())
//	zap.ReplaceGlobals(LG) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
//	return
//}
//
//func getEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	encoderConfig.TimeKey = "time"
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
//	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
//	return zapcore.NewJSONEncoder(encoderConfig)
//}
//
//func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
//	lumberJackLogger := &lumberjack.Logger{
//		Filename:   filename,
//		MaxSize:    maxSize,
//		MaxBackups: maxBackup,
//		MaxAge:     maxAge,
//	}
//	return zapcore.AddSync(lumberJackLogger)
//}
