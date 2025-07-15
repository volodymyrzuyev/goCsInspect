package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	mainhelpers "github.com/volodymyrzuyev/goCsInspect/internal/mainHelpers"
	"github.com/volodymyrzuyev/goCsInspect/internal/web"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/logger"
)

var (
	cfgLocation string
)

func init() {
	flag.StringVar(
		&cfgLocation, "config", config.DefaultConfigLocation, "configuration file used for the api",
	)
}

func main() {
	flag.Parse()

	cfg, err := config.ParseConfig(cfgLocation)
	if err != nil {
		fmt.Println("invalid configuration location, stoping")
		os.Exit(1)
	}

	l := slog.New(logger.NewHandler(cfg.LogLevel, os.Stdout))
	slog.SetDefault(l)
	lt := l.WithGroup("Main")

	cm := mainhelpers.InitDefaultClientManager(cfg, lt, l)

	for _, cli := range cfg.Accounts {
		err := cm.AddClient(cli)
		if err != nil {
			lt.Warn(
				fmt.Sprintf("client %v unable to login, won't be used", cli.Username),
				"error",
				err,
			)
		}
	}

	server := web.NewServer(cm, l)

	server.Run(cfg.BindIP)
}
