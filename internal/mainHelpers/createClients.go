package mainhelpers

import (
	"log/slog"

	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/client"
	"github.com/volodymyrzuyev/goCsInspect/pkg/clientmanagement"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
)

func CreateAndEnrollClients(cfg config.Config, cm clientmanagement.Manager) {
	mainLogger := slog.Default().WithGroup("Main")
	gc := gcHandler.NewGcHandler()
	for _, cli := range cfg.Accounts {
		cli, err := client.New(cli, cfg.ClientCooldown, gc)
		if err != nil {
			mainLogger.Warn("unable to create client", "username", cli.Username(), "error", err)
		}
		err = cm.AddClient(cli)
		if err != nil {
			mainLogger.Warn("unable to login into client", "username", cli.Username(), "error", err)
		}
	}
}
