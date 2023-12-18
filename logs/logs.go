package logs

import (
	"github.com/preceeder/gobase/env"
	"github.com/preceeder/gobase/utils"
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

func InitLogWithViper(config viper.Viper) {
	//logConfig := ReadLogConfig(config)
	logConfig := LogConfig{}
	utils.ReadViperConfig(config, "log", &logConfig)
	initSlog(logConfig)
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

	opt := &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			d := a.Value.Any().(*slog.Source)
			d.File = ""
			a.Value = slog.AnyValue(d)
		}
		return a
	}}
	var writeBuild io.Writer
	if cfg.StdOut == "1" {
		writeBuild = io.MultiWriter(os.Stdout, lumberJackInfo)
	} else {
		writeBuild = io.MultiWriter(lumberJackInfo)
	}
	var log *slog.Logger
	log = slog.New(slog.NewJSONHandler(writeBuild, opt))
	slog.SetDefault(log)
	if env.GetEnv("env") == "local" {
		//slog.Default()
		log = slog.New(slog.NewTextHandler(writeBuild, opt))
		slog.SetDefault(log)
	}
}
