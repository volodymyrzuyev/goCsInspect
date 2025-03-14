package client

import (
	"errors"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/request"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

var (
	UnableToLogin      = errors.New("Unable to login")
	TimeoutPassed      = errors.New("Request wait passed")
	ClientNotAvaliable = errors.New("Client is not avaliable")
)

type Client interface {
	LogIn(types.Credentials, steam.GCPacketHandler) error
	LogOut()
	RequestSkin(types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error)
	Avaliable() bool
}

type client struct {
	requestHandler request.RequestHandler
	log            logger.Logger

	client      *steam.Client
	username    string
	lastUsed    time.Time
	exitChan    *chan bool
	avaliable   bool
	disconected bool
}

var eventLoopRunner = runEventLoop

func (c *client) LogIn(creds types.Credentials, handler steam.GCPacketHandler) error {
	c.log.Debug("Login request for %v", creds.Username)
	logInInfo, err := getLoginDetails(creds)
	if err != nil {
		c.log.Error("%v credentials are invalid", creds.Username)
		return err
	}
	c.username = creds.Username

	curClient := steam.NewClient()

	loginStatus := make(chan bool)
	exitChan := make(chan bool, 1)
	c.log.Debug("Starting event loop for %v", creds.Username)
	go eventLoopRunner(curClient, logInInfo, loginStatus, c.log, exitChan)

	c.client = curClient
	c.lastUsed = time.Now().Add(-5 * time.Second)
	c.exitChan = &exitChan

	select {
	case <-loginStatus:
		curClient.GC.RegisterPacketHandler(handler)
		c.avaliable = true
		c.disconected = false
		c.log.Debug("Client: %v properly logged in", c.username)
		return nil
	case <-time.After(types.TimeOutDuration):
		c.log.Debug("%v timed out", c.username)
		exitChan <- true
		return UnableToLogin
	}
}

func (c *client) LogOut() {
	c.log.Debug("Requesting to stop event loop for Client %v", c.username)
	*c.exitChan <- true
	c.avaliable = false
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

	respChan := c.sendGCRequest(requestProto)

	select {
	case resp := <-*respChan:
		c.log.Debug("Client: %v successfully got response for %v", c.username, inspectParams.A)
		return resp.Response, resp.Error
	case <-time.After(types.TimeOutDuration):
		c.log.Error("Client: %v timed out when requesting %v", c.username, inspectParams.A)
		return nil, TimeoutPassed
	}
}

func (c *client) sendGCRequest(msg *csProto.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest) *chan types.Response {
	c.log.Debug("Client: %v is sending GC message to inspect %v", c.username, msg.ParamA)
	proto := gamecoordinator.NewGCMsgProtobuf(730, uint32(types.InspectRequestProtoID), msg)
	c.client.GC.Write(proto)
	return c.requestHandler.RegisterRequest(*msg.ParamA)
}

func (c *client) Avaliable() bool {
	select {
	case <-*c.exitChan:
		c.disconected = true
	default:
		break
	}
	willBeAvaliable := c.lastUsed.Add(types.RequestCooldown)
	return c.avaliable && !c.disconected && c.lastUsed.After(willBeAvaliable)
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

func runEventLoop(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
	curClient.Connect()
	select {
	case <-exit:
		exit <- true
		log.Debug("Stopping event loop for %v", logInInfo.Username)
		return
	default:
		for event := range curClient.Events() {
			switch e := event.(type) {
			case *steam.ConnectedEvent:
				log.Debug("Connection event Username: %v", logInInfo.Username)
				curClient.Auth.LogOn(&logInInfo)
			case *steam.LoggedOnEvent:
				curClient.Social.SetPersonaState(steamlang.EPersonaState_Busy)
				log.Debug("Account %v fully connected", logInInfo.Username)
				login <- true
			case steam.FatalErrorEvent:
				log.Error("Account: %v disconected due to error: %v", logInInfo.Username, e)
				exit <- true
			case error:
				log.Error("Account: %v go an error: %v", logInInfo.Username, e)
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

func NewClient(requestHandler request.RequestHandler, log logger.Logger) Client {
	if requestHandler == nil || log == nil {
		panic("Reqest handler and log can not be null")
	}

	return &client{
		requestHandler: requestHandler,
		log:            log,
		disconected:    true,
	}
}
