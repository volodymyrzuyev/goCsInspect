package client

import (
	"errors"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/response"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

var (
	UnableToLogin      = errors.New("Unable to login")
	TimeoutPassed      = errors.New("Request wait passed")
	ClientNotAvaliable = errors.New("Client is not avaliable")
)

type Client interface {
	LogIn(types.Credentials) error
	LogOut()
	RequestSkin(types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error)
	Avaliable() bool
}

type client struct {
	log logger.Logger

	client       *steam.Client
	username     string
	lastUsed     time.Time
	exitChan     *chan bool
	responseChan *chan types.Response
	disconected  bool
}

var eventLoopRunner = runEventLoop

func (c *client) LogIn(creds types.Credentials) error {
	c.log.Debug("Login request for %v", creds.Username)
	logInInfo, err := getLoginDetails(creds)
	if err != nil {
		c.log.Error("%v credentials are invalid", creds.Username)
		return err
	}
	c.setPreLoginState(creds.Username)

	c.log.Debug("Starting event loop for %v", creds.Username)
	loginStatus := make(chan bool)
	go eventLoopRunner(c.client, logInInfo, loginStatus, c.log, c.exitChan)

	handler := response.NewResponseHandler(c.log, c.responseChan, c.username)

	select {
	case <-loginStatus:
		c.client.GC.RegisterPacketHandler(handler)
		c.setPostSuccessfulLoginState()
		c.log.Debug("Client: %v properly logged in", c.username)
		return nil
	case <-time.After(types.TimeOutDuration):
		c.log.Debug("%v timed out during login", c.username)
		*c.exitChan <- true
		return UnableToLogin
	}
}

func (c *client) setPreLoginState(username string) {
	c.username = username
	curClient := steam.NewClient()
	c.client = curClient

	resp := make(chan types.Response)
	c.responseChan = &resp
	exit := make(chan bool, 1)
	c.exitChan = &exit

}

func (c *client) setPostSuccessfulLoginState() {
	c.disconected = false
	c.lastUsed = time.Now().Add(-5 * time.Second)
}

func (c *client) LogOut() {
	c.log.Debug("Requesting to stop event loop for Client %v", c.username)
	*c.exitChan <- true
	c.disconected = true
}

func (c *client) RequestSkin(inspectParams types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error) {
	c.log.Debug("Client: %v is requested to inspect %v", c.username, inspectParams.A)
	if !c.Avaliable() {
		c.log.Debug("Client: %v is not avaliable", c.username)
		return nil, ClientNotAvaliable
	}

	c.lastUsed = time.Now()

	requestProto, err := getInspectDetails(inspectParams)
	if err != nil {
		c.log.Debug("Client: %v, item: %v invalid link", c.username, inspectParams.A)
		return nil, err
	}

	c.log.Debug("Client: %v is sending GC message to inspect %v", c.username, requestProto.ParamA)
	proto := gamecoordinator.NewGCMsgProtobuf(730, uint32(types.InspectRequestProtoID), requestProto)
	c.client.GC.Write(proto)

	select {
	case resp := <-*c.responseChan:
		c.log.Debug("Client: %v successfully got response for %v", c.username, inspectParams.A)
		c.lastUsed = time.Now()
		return resp.Response, resp.Error
	case <-time.After(types.TimeOutDuration):
		c.log.Error("Client: %v timed out when requesting %v", c.username, inspectParams.A)
		return nil, TimeoutPassed
	}
}

func (c *client) Avaliable() bool {
	select {
	case <-*c.exitChan:
		c.disconected = true
	default:
	}
	willBeAvaliable := c.lastUsed.Add(types.RequestCooldown)
	return !c.disconected && time.Now().After(willBeAvaliable)
}

func getLoginDetails(creds types.Credentials) (steam.LogOnDetails, error) {
	err := creds.Validate()
	if err != nil {
		return steam.LogOnDetails{}, err
	}

	twoFA, err := creds.Get2FC()
	if err != nil {
		return steam.LogOnDetails{}, err
	}

	logInInfo := steam.LogOnDetails{
		Username:      creds.Username,
		Password:      creds.Password,
		TwoFactorCode: twoFA,
	}

	return logInInfo, nil
}

func runEventLoop(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit *chan bool) {
	curClient.Connect()
	select {
	case <-*exit:
		*exit <- true
		log.Debug("Stopping event loop for %v", logInInfo.Username)
		return
	default:
		for event := range curClient.Events() {
			switch e := event.(type) {
			case *steam.ConnectedEvent:
				log.Debug("Connection event Username: %v", logInInfo.Username)
				curClient.Auth.LogOn(&logInInfo)
			case *steam.LoggedOnEvent:
				curClient.GC.SetGamesPlayed(730)
				log.Debug("Client: %v fully connected", logInInfo.Username)
				login <- true
			case *steam.DisconnectedEvent:
				log.Error("Client: %v disconnected", logInInfo.Username)
				*exit <- true
			case steam.FatalErrorEvent:
				log.Error("Client: %v disconected due to error: %v", logInInfo.Username, e)
				*exit <- true
			case error:
				log.Error("Client: %v go an error: %v", logInInfo.Username, e)
			}
		}
	}

}

func getInspectDetails(params types.InspectParameters) (*csProto.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	requestProto := csProto.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest{
		ParamS: &params.S,
		ParamA: &params.A,
		ParamD: &params.D,
		ParamM: &params.M,
	}

	return &requestProto, nil
}

func NewClient(log logger.Logger) Client {
	if log == nil {
		panic("Log is needed to run client")
	}

	return &client{
		log:         log,
		disconected: true,
	}
}
