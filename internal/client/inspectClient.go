package client

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"

	"github.com/volodymyrzuyev/goCsInspect/config"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/types"
)

type InspectClient interface {
	IsLoggedIn() bool
	IsAvailable() bool
	Username() string

	LogOff()
	LogIn() error
	Reconnect() error

	InspectItem(params *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest) (*protobuf.CEconItemPreviewDataBlock, error)
}

type inspectClient struct {
	lastUsed time.Time

	exitCh chan bool
	client *steam.Client

	creds     types.Credentials
	config    config.ClientConfig
	gcHandler gcHandler.GcHandler
}

func NewInspectClient(config config.ClientConfig, gcHandler gcHandler.GcHandler, creds types.Credentials) (InspectClient, error) {
	_, err := creds.GenerateLogOnDetails()
	if err != nil {
		return nil, err
	}

	newClient := &inspectClient{
		client: steam.NewClient(),
		exitCh: make(chan bool),

		creds:     creds,
		config:    config,
		gcHandler: gcHandler,
	}

	return newClient, nil
}

func (c *inspectClient) IsLoggedIn() bool {
	return c.client.Connected()
}

func (c *inspectClient) IsAvailable() bool {
	willBeAvaliable := c.lastUsed.Add(c.config.RequestCooldown)
	return c.IsLoggedIn() && time.Now().After(willBeAvaliable)
}

func (c *inspectClient) Username() string {
	return c.creds.Username
}

func (c *inspectClient) LogOff() {
	slog.Info("Stopping client", "username", c.creds.Username)
	if !c.IsLoggedIn() {
		return
	}
	c.client.Disconnect()
	c.exitCh <- true
}

const csAppID = 730

func runClientLoop(c config.ClientConfig, client *steam.Client, creds steam.LogOnDetails, exitCh <-chan bool, loginCh chan<- error) {
	auth := newAuth(client, &creds, loginCh)
	serverList := newServerList(client, "servers/list.json")
	debug := newDebug(creds.Username, c.DebugLogger, c.DebugLogger)
	if c.IsDebug {
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
				slog.Debug("Client channel disconnected, leaving client loop",
					"username", auth.details.Username)
				return
			}

			auth.HandleEvent(event)
			serverList.HandleEvent(event)
			if c.IsDebug {
				debug.HandleEvent(event)
			}

			switch e := event.(type) {
			case error:
				slog.Debug("Steam client event error",
					"username", auth.details.Username, "error", e.Error())

			case *steam.LoggedOnEvent:
				client.Social.SetPersonaState(steamlang.EPersonaState_Online)
				client.GC.SetGamesPlayed(csAppID)
			}
		}
	}
}

func (c *inspectClient) LogIn() error {
	slog.Debug("Login attempt", "username", c.creds.Username)
	logOnDetails, err := c.creds.GenerateLogOnDetails()
	if err != nil {
		slog.Error("Invalid credentials",
			"username", c.creds.Username, "error", err.Error())
		return err
	}

	logIn := make(chan error)

	go runClientLoop(c.config, c.client, logOnDetails, c.exitCh, logIn)

	select {
	case err := <-logIn:
		if err != nil {
			slog.Error("Client got error during connection",
				"username", c.creds.Username, "error", err.Error())
			return err
		}
		slog.Info("Client login complete", "username", c.creds.Username)
		c.client.GC.RegisterPacketHandler(c.gcHandler)
		c.lastUsed = time.Now().Add(-c.config.RequestCooldown * 2)
		return nil
	case <-time.After(c.config.TimeOutDuration):
		c.LogOff()
		slog.Error("Client timed out during login", "username", c.creds.Username)
		return errors.ErrClientUnableToConnect
	}
}

func (c *inspectClient) Reconnect() error {
	slog.Warn("Relogin into client", "username", c.creds.Username)
	c.LogOff()
	return c.LogIn()
}

const inspectRequestProtoID = 9156

func (c *inspectClient) InspectItem(params *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest) (*protobuf.CEconItemPreviewDataBlock, error) {
	slog.Debug("Client requested to inspect skin",
		"username", c.creds.Username,
		"lastUsed", c.lastUsed.Format(time.TimeOnly),
		"proto", fmt.Sprintf("%+v", params))

	if params == nil {
		slog.Error("Invalid inspect proto", "proto", fmt.Sprintf("%+v", params))
		return nil, errors.ErrInvalidInspectLink
	}

	if !c.IsAvailable() {
		slog.Error("Client not available to inspect skin",
			"username", c.creds.Username)
		return nil, errors.ErrClientUnavailable
	}

	proto := gamecoordinator.NewGCMsgProtobuf(csAppID, uint32(inspectRequestProtoID), params)
	slog.Debug("Sending inspect request",
		"username", c.creds.Username, "proto", fmt.Sprintf("%+v", params))
	c.client.GC.Write(proto)
	c.lastUsed = time.Now()

	return c.gcHandler.GetResponse(params.GetParamA())
}
