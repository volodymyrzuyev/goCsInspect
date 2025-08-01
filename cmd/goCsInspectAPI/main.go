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

	l := slog.New(logger.NewHandler(cfg.GetLogLevel(), os.Stdout))
	slog.SetDefault(l)

	cm := mainhelpers.InitDefaultClientManager(cfg)

	mainhelpers.CreateAndEnrollClients(cfg, cm)

	server := web.NewServer(cm)

	server.Run(cfg.BindIP)
}
