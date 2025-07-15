package mainhelpers

import (
	"log/slog"

	"github.com/volodymyrzuyev/goCsInspect/pkg/client"
	"github.com/volodymyrzuyev/goCsInspect/pkg/clientmanagement"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/gamecordinator"
)

func CreateAndEnrollClients(cfg config.Config, cm clientmanagement.Manager) {
	mainLogger := slog.Default().WithGroup("Main")
	gc := gamecordinator.NewGcHandler()
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
