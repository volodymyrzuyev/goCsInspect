package client

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"

	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/common/consts"
	"github.com/volodymyrzuyev/goCsInspect/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/common/types"
	"github.com/volodymyrzuyev/goCsInspect/config"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
)

type InspectClient interface {
	IsLoggedIn() bool
	IsAvaliable() bool

	LogIn(credentials types.Credentials) error
	LogOff()

	InspectItem(params types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error)
}

type inspectClient struct {
	username string
	lastUsed time.Time

	exitCh chan bool
	client *steam.Client

	gcHandler gcHandler.GcHandler
}

func NewInspectClient(gcHandler gcHandler.GcHandler) InspectClient {
	return &inspectClient{
		client:    steam.NewClient(),
		exitCh:    make(chan bool),
		gcHandler: gcHandler,
	}
}

func (c *inspectClient) LogIn(creds types.Credentials) error {
	slog.Debug("Login attempt", "username", creds.Username)

	logOnDetails, err := creds.GenerateLogOnDetails()
	if err != nil {
		slog.Error("Invalid credentials", "username", creds.Username, "error", err.Error())
		return err
	}

	logIn := make(chan error)

	go runClientLoop(c.client, logOnDetails, c.exitCh, logIn)

	c.username = logOnDetails.Username

	select {
	case err := <-logIn:
		if err != nil {
			slog.Error("Client got error during connection", "username", c.username, "error", err.Error())
			return err
		}
		slog.Info("Client login complete", "username", c.username)
		c.client.GC.RegisterPacketHandler(c.gcHandler)
		c.lastUsed = time.Now().Add(-config.RequestCooldown * 2)
		return nil
	case <-time.After(config.TimeOutDuration):
		c.LogOff()
		slog.Error("Client timed out during login", "username", c.username)
		return errors.ErrClientUnableToConnect
	}
}

func (c *inspectClient) LogOff() {
	slog.Info("Stopping client", "username", c.username)
	if !c.IsLoggedIn() {
		return
	}
	c.client.Disconnect()
	c.exitCh <- true
}

func (c *inspectClient) IsLoggedIn() bool {
	return c.client != nil && c.client.Connected()
}
func (c *inspectClient) IsAvaliable() bool {
	willBeAvaliable := c.lastUsed.Add(config.RequestCooldown)
	return c.IsLoggedIn() && time.Now().After(willBeAvaliable)
}

func (c *inspectClient) InspectItem(params types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error) {
	slog.Debug("Client requested to inspect skin", "username", c.username, "lastUsed", c.lastUsed.Format(time.TimeOnly), "inspect_params", fmt.Sprintf("%+v", params))
	if !c.IsAvaliable() {
		slog.Error("Client not avaliable to inspect skin", "username", c.username)
		return nil, errors.ErrClientUnavaliable
	}

	requestProto, err := params.GenerateGcRequestProto()
	if err != nil {
		slog.Error("Client unable to parse inspect link", "username", c.username, "inspect_params", params)
		return nil, err
	}

	proto := gamecoordinator.NewGCMsgProtobuf(consts.CsAppID, uint32(consts.InspectRequestProtoID), requestProto)
	slog.Debug("Sending inspect request", "username", c.username, "inspect_params", fmt.Sprintf("%+v", params))
	c.client.GC.Write(proto)
	c.lastUsed = time.Now()

	return c.gcHandler.GetResponse(params.A)
}

func runClientLoop(client *steam.Client, creds steam.LogOnDetails, exitCh <-chan bool, loginCh chan<- error) {
	auth := newAuth(client, &creds, loginCh)
	serverList := newServerList(client, "servers/list.json")
	debug := newDebug(creds.Username, config.DebugLogger, config.DebugLogger)
	if config.IsDebug {
		client.RegisterPacketHandler(debug)
	}

	serverList.Connect()

	for {
		select {
		case <-exitCh:
			slog.Info("Stopping client loop", "username", auth.details.Username)
			return
		case event, ok := <-client.Events():
			if !ok {
				slog.Debug("Client channel disconnected, leaving client loop", "username", auth.details.Username)
				return
			}

			auth.HandleEvent(event)
			serverList.HandleEvent(event)
			if config.IsDebug {
				debug.HandleEvent(event)
			}

			switch e := event.(type) {
			case error:
				slog.Warn("Steam client event error", "username", auth.details.Username, "error", e.Error())
			case *steam.LoggedOnEvent:
				client.Social.SetPersonaState(steamlang.EPersonaState_Online)
				client.GC.SetGamesPlayed(consts.CsAppID)
			}
		}
	}

}
