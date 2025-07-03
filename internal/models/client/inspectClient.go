package client

import (
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
)

type InspectClient interface {
	IsLoggedIn() bool

	LogIn(credentials types.Credentials) error
	LogOff()

	InspectItem(params types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error)
}

type inspectClient struct {
	username string

	exitCh       chan bool
	gcResponseCh chan types.Response
	client       *steam.Client
}

func NewInspectClient() InspectClient {
	return &inspectClient{
		client:       steam.NewClient(),
		gcResponseCh: make(chan types.Response),
		exitCh:       make(chan bool),
	}
}

func (c *inspectClient) LogIn(creds types.Credentials) error {
	slog.Debug("Login attempt", "username", creds.Username)
	logOnDetails, err := creds.GenerateLogOnDetails()
	if err != nil {
		slog.Error("Invalid Credentials", "username", creds.Username, "error", err)
		return err
	}

	logIn := make(chan error)
	go runClientLoop(c.client, logOnDetails, c.exitCh, logIn)

	select {
	case err := <-logIn:
		if err != nil {
			slog.Info("Client got error during connection", "username", logOnDetails.Username, "error", err.Error())
			return err
		}
		slog.Info("Client login complete", "username", logOnDetails.Username)
		c.client.GC.RegisterPacketHandler(NewGcHandler(c.gcResponseCh, logOnDetails.Username))
		slog.Debug("Registered GC handler", "username", logOnDetails.Username)
		return nil
	case <-time.After(config.TimeOutDuration):
		c.LogOff()
		slog.Warn("Client timed out", "username", logOnDetails.Username)
		return errors.ErrClientUnableToConnect
	}
}

func (c *inspectClient) LogOff() {
	if !c.IsLoggedIn() {
		return
	}
	c.client.Disconnect()
	c.exitCh <- true
}

func (c *inspectClient) IsLoggedIn() bool {
	return c.client != nil && c.client.Connected()
}

func (c *inspectClient) InspectItem(params types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error) {
	if !c.IsLoggedIn() {
		return nil, errors.ErrClientUnableToConnect
	}

	requestProto, err := params.GenerateGcRequestProto()
	if err != nil {
		return nil, err
	}

	proto := gamecoordinator.NewGCMsgProtobuf(consts.CsAppID, uint32(consts.InspectRequestProtoID), requestProto)
	c.client.GC.Write(proto)

	select {
	case response := <-c.gcResponseCh:
		if response.Error != nil {
			slog.Debug("Client error when fetching skin", "username", c.username, "skin_id", params.A, "error", response.Error.Error())
			return nil, response.Error
		}

		slog.Debug("Client got skin details", "username", c.username, "skin_id", params.A)
		return response.Response, nil
	case <-time.After(config.TimeOutDuration):
		slog.Warn("Client timed out fetching skin", "username", c.username, "skin_id", params.A)
		return nil, errors.ErrClientTimeout
	}
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
				slog.Debug("Client channel disconnected, leaving client loop", "username", auth.details.Username)
				return
			}

			// debug.HandleEvent(event)
			auth.HandleEvent(event)
			serverList.HandleEvent(event)
			switch e := event.(type) {
			case error:
				slog.Error("Steam client event error", "username", auth.details.Username, "error", e)
			case *steam.LoggedOnEvent:
				client.Social.SetPersonaState(steamlang.EPersonaState_Online)
			}
		}
	}

}
