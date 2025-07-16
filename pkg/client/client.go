package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"

	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/creds"
	gameC "github.com/volodymyrzuyev/goCsInspect/pkg/gamecordinator"
)

// Client interface is used to interface with steam game gamecoordinator to
// fetch data about items. Responsible for logging in into steam and handling
// all events that are sent by the game gamecoordinator as well as other steam
// services
type Client interface {
	// Used to check if client is healthy
	// True means the client is safe to use
	IsLoggedIn() bool
	// Used to check if client is ready to handle a new request
	// True means a new request can be sent
	IsAvailable() bool
	// Returns the username of client
	Username() string

	// Stops the client
	LogOff()
	// Logins into the steam network
	// expects the implementation handles creds on it's own
	LogIn() error
	// Stops, after which re logins into steam
	Reconnect() error

	// Fetches an item from steam
	InspectItem(
		ctx context.Context,
		params *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest,
	) (*protobuf.CEconItemPreviewDataBlock, error)
}

type inspectClient struct {
	lastUsed time.Time

	exitCh chan bool
	client *steam.Client

	cooldown  time.Duration
	creds     creds.Account
	gcHandler gameC.Handler
	l         *slog.Logger
}

func New(creds creds.Account, cooldown time.Duration, gc gameC.Handler) (Client, error) {
	if _, err := creds.GenerateLogOnDetails(); err != nil {
		slog.Error("invalid client credentials")
		return nil, err
	}

	newClient := &inspectClient{
		client: steam.NewClient(),
		exitCh: make(chan bool),

		cooldown:  cooldown,
		creds:     creds,
		gcHandler: gc,
		l:         slog.Default().WithGroup("Client." + creds.Username),
	}

	return newClient, nil
}

func (c *inspectClient) IsLoggedIn() bool {
	return c.client.Connected()
}

func (c *inspectClient) IsAvailable() bool {
	willBeAvaliable := c.lastUsed.Add(c.cooldown)
	return c.IsLoggedIn() && time.Now().After(willBeAvaliable)
}

func (c *inspectClient) Username() string {
	return c.creds.Username
}

func (c *inspectClient) LogOff() {
	c.l.Info("stopping client")
	if !c.IsLoggedIn() {
		return
	}
	c.client.Disconnect()
	c.exitCh <- true
}

const csAppID = 730

func runClientLoop(
	client *steam.Client,
	creds steam.LogOnDetails,
	exitCh <-chan bool,
	loginCh chan<- error,
	l *slog.Logger,
) {

	err := serverManager.Connect(client)
	if err != nil {
		loginCh <- err
		return
	}

	for {
		select {
		case <-exitCh:
			l.Info("stopping client loop")
			return

		case event, ok := <-client.Events():
			if !ok {
				l.Error("leaving client loop")
				return
			}

			switch e := event.(type) {
			case error:
				l.Debug("client event error", "error", e.Error())
			case *steam.ConnectedEvent:
				client.Auth.LogOn(&creds)

			case *steam.LoggedOnEvent:
				loginCh <- nil
				client.Social.SetPersonaState(steamlang.EPersonaState_Online)
				client.GC.SetGamesPlayed(csAppID)

			case *steam.LogOnFailedEvent:
				loginCh <- fmt.Errorf("%s", e.Result.String())
			}

		}
	}
}

const timeoutDuration = 10 * time.Second

func (c *inspectClient) LogIn() error {
	c.l.Debug("Login attempt", "username", c.creds.Username)
	logOnDetails, err := c.creds.GenerateLogOnDetails()
	if err != nil {
		c.l.Error("invalid client credentials", "error", err.Error())
		return err
	}

	logIn := make(chan error)

	go runClientLoop(c.client, logOnDetails, c.exitCh, logIn, c.l)

	select {
	case err := <-logIn:
		if err != nil {
			c.l.Error("error during steam login", "error", err.Error())
			return err
		}
		c.l.Info("login complete, client ready")
		c.client.GC.RegisterPacketHandler(c.gcHandler)
		c.lastUsed = time.Now().Add(-c.cooldown * 2)
		return nil
	case <-time.After(timeoutDuration):
		c.LogOff()
		c.l.Error("timed out during login")
		return errors.ErrClientUnableToConnect
	}
}

func (c *inspectClient) Reconnect() error {
	c.l.Info("relogin attempt")
	c.LogOff()
	return c.LogIn()
}

const inspectRequestProtoID = 9156

func (c *inspectClient) InspectItem(
	ctx context.Context,
	params *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest,
) (*protobuf.CEconItemPreviewDataBlock, error) {

	c.l.Debug("item previw block requested", "proto", fmt.Sprintf("%+v", params))

	if params == nil {
		c.l.Error("invalid itemp previw block request params")
		return nil, errors.ErrInvalidInspectLink
	}

	if !c.IsAvailable() {
		c.l.Error("client not available")
		return nil, errors.ErrClientUnavailable
	}

	select {
	case <-ctx.Done():
		return nil, errors.ErrClientTimeout
	default:
		proto := gamecoordinator.NewGCMsgProtobuf(csAppID, uint32(inspectRequestProtoID), params)
		c.l.Debug("sent item preview block packet")
		c.client.GC.Write(proto)
		c.lastUsed = time.Now()

		return c.gcHandler.GetResponse(ctx, params.GetParamA())
	}
}
