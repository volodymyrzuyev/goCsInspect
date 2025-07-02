package client

import (
	"log/slog"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"
	"github.com/volodymyrzuyev/goCsInspect/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/common/types"
	"github.com/volodymyrzuyev/goCsInspect/config"
)

type InspectClient interface {
	IsLoggedIn() bool

	LogIn(credentials types.Credentials) error
	LogOff()
}

type inspectClient struct {
	username string

	exitCh *chan bool
	client *steam.Client
}

func NewInspectClient() InspectClient {
	logOff := make(chan bool)

	inspectClient := inspectClient{
		client: steam.NewClient(),
		exitCh: &logOff,
	}

	return &inspectClient
}

func (c *inspectClient) LogIn(creds types.Credentials) error {
	slog.Debug("Login attempt", "username", creds.Username)
	logOnDetails, err := creds.GenerateLogOnDetails()
	if err != nil {
		slog.Error("Invalid Credentials", "username", creds.Username, "error", err)
		return err
	}
	logIn := make(chan error)

	go runClientLoop(c.client, logOnDetails, *c.exitCh, logIn)

	select {
	case err := <-logIn:
		if err != nil {
			slog.Info("Client got error during connection", "username", creds.Username, "error", err.Error())
			return err
		}

		slog.Info("Client login complete", "username", logOnDetails.Username)
		return nil

	case <-time.After(config.TimeOutDuration):
		c.LogOff()
		slog.Warn("Client timedout", "username", logOnDetails.Username)
		return errors.UnableToConnect

	}
}

func (c *inspectClient) LogOff() {
	if !c.client.Connected() {
		return
	}

	c.client.Disconnect()
	*c.exitCh <- true
}

func (c *inspectClient) IsLoggedIn() bool {
	return c.client.Connected()
}

func runClientLoop(client *steam.Client, creds steam.LogOnDetails, exitCh <-chan bool, loginCh chan<- error) {
	auth := newAuth(client, &creds, loginCh)
	serverList := newServerList(client, "servers/list.json")
	debug := newDebug(creds.Username, slog.Default(), slog.Default())
	client.RegisterPacketHandler(debug)

	serverList.Connect()

	for {
		select {
		case <-exitCh:
			slog.Debug("Stopping client loop", "username", auth.details.Username)
			return

		case event, ok := <-client.Events():
			if !ok {
				slog.Debug("Client chanel disconected, leaving client loop", "username", auth.details.Username)
				return
			}

			// debug.HandleEvent(event)
			auth.HandleEvent(event)
			serverList.HandleEvent(event)

			switch event.(type) {
			case error:
			case *steam.LoggedOnEvent:
				client.Social.SetPersonaState(steamlang.EPersonaState_Online)
			}
		}
	}

}
