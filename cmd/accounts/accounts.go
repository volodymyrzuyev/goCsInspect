package accounts

import (
	"fmt"
	"os"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"
	gt "github.com/volodymyrzuyev/goCsInspect/cmd/globalTypes"
	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
	"google.golang.org/protobuf/proto"
)

type Accounts interface {
	AddAccount(Credentials) error
	InspectWeapon(params gt.InspectParams) (gt.Skin, error)
}

type account struct {
	client    *steam.Client
	username  string
	lastUsed  time.Time
	avaliable bool
}

type accounts struct {
	clients   []*account
	handler   steam.GCPacketHandler
	responses chan bool
}

func (a *accounts) AddAccount(creds Credentials) error {
	err := creds.validate()
	if err != nil {
		return err
	}

	authCode, err := creds.get2FA()
	if err != nil {
		return err
	}

	logInInfo := steam.LogOnDetails{
		Username:      creds.Username,
		Password:      creds.Password,
		TwoFactorCode: authCode,
	}

	client := steam.NewClient()
	loginComplete := make(chan bool, 1)

	go a.handleEvents(loginComplete, client, logInInfo)

	select {
	case <-loginComplete:
		logger.INFO.Printf("%v logged added to accounts", creds.Username)
		return nil
	case <-time.After(30 * time.Second):
		logger.ERROR.Printf("%v login timeout", creds.Username)
		return LoginTimeOut
	}
}

func (a *accounts) InspectWeapon(params gt.InspectParams) (gt.Skin, error) {
	clientIdx := a.getNextFreeAccount()
	if clientIdx < 0 {
		return gt.Skin{}, NoAvaliableAccounts
	}

	curClient := a.clients[clientIdx]

	pip := &protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest{
		ParamM: proto.Uint64(uint64(params.ParamM)),
		ParamA: proto.Uint64(uint64(params.ParamA)),
		ParamD: proto.Uint64(uint64(params.ParamD)),
		ParamS: proto.Uint64(uint64(params.ParamS)),
	}

	// crafting the message
	msg := gamecoordinator.NewGCMsgProtobuf(730, 9156, pip)

	curClient.client.GC.Write(msg)
	logger.DEBUG.Printf("%v is requesting skin, itemID: %v", curClient.username, params.ParamA)
	for i := 0; i < 1; i++ {
		select {
		case kik := <-a.responses:
			if kik {
				fmt.Println("Got response")
			}
		}
	}

	return gt.Skin{}, nil
}

func (a *accounts) getNextFreeAccount() int {
	for i, c := range a.clients {

		if time.Now().Before(c.lastUsed.Add(1 * time.Second)) {
			c.avaliable = true
			logger.DEBUG.Printf("Cooldown passed on %v", c.username)
		}

		if c.avaliable {
			return i
		}
	}

	return -1
}

func (a *accounts) handleEvents(loginComplete chan bool, client *steam.Client, logInInfo steam.LogOnDetails) {
	client.Connect()
	for event := range client.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			client.Auth.LogOn(&logInInfo)
			logger.DEBUG.Printf("account: %v sent login info", logInInfo.Username)

		case *steam.MachineAuthUpdateEvent:
			os.WriteFile("sentry", e.Hash, 0666)
			logger.DEBUG.Printf("account: %v wrote sentry file", logInInfo.Username)

		case *steam.LoggedOnEvent:
			client.Social.SetPersonaState(steamlang.EPersonaState_Online)
			client.GC.SetGamesPlayed(730)
			newAccount := account{
				client:    client,
				lastUsed:  time.Now(),
				avaliable: true,
				username:  logInInfo.Username,
			}
			a.clients = append(a.clients, &newAccount)
			logger.DEBUG.Printf("account: %v logged in", logInInfo.Username)
			client.GC.RegisterPacketHandler(a.handler)
			loginComplete <- true

		case steam.FatalErrorEvent:
			logger.ERROR.Printf("FatalEvent: %v", e)

		case error:
			logger.ERROR.Printf("Error: %v", e)
		}
	}
}

func NewAccountsList(handler steam.GCPacketHandler, responses chan bool) Accounts {
	if handler == nil {
		panic("No handler func for account")
	}
	return &accounts{handler: handler, responses: responses}
}
