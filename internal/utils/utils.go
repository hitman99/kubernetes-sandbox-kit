package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/hitman99/kubernetes-sandbox/internal/config"
	"os"
	"runtime"
	"strings"
)

func SetupLogger() *log.Logger {
	logger := log.New()
	dir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Panic("cannot get current directory")
	}
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		DataKey:          "",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			file = strings.ReplaceAll(f.File, dir, "")
			function = strings.TrimPrefix(f.Function, f.Function[:strings.LastIndex(f.Function, "/")+1])
			return
		},
		PrettyPrint: false,
	})

	logger.SetReportCaller(true)
	cfg, _ := config.Get()
	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = log.InfoLevel
	}
	logger.SetLevel(lvl)
	return logger
}
