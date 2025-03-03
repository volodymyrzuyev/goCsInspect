package accounts

import (
	"errors"
	"os"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"
	gt "github.com/volodymyrzuyev/goCsInspect/cmd/globalTypes"
	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
	req "github.com/volodymyrzuyev/goCsInspect/cmd/requests"
	"github.com/volodymyrzuyev/goCsInspect/cmd/storage"
	"google.golang.org/protobuf/proto"
)

const timeOutDuration = time.Second * 30

type Accounts interface {
	AddAccount(Credentials) error
	InspectWeapon(params gt.InspectParams) (gt.Item, error)
}

type account struct {
	client    *steam.Client
	username  string
	lastUsed  time.Time
	avaliable bool
}

type accounts struct {
	clients    []*account
	handler    steam.GCPacketHandler
	reqHandler req.RequestHandler
	db         storage.Storage
}

func (a *accounts) AddAccount(creds Credentials) error {
	err := creds.validate()
	if err != nil {
		logger.ERROR.Printf("Could not validate %v. Err: %v", creds.Username, err)
		return err
	}

	authCode, err := creds.get2FA()
	if err != nil {
		logger.ERROR.Printf("Could not get 2FA code for %v. Err: %v", creds.Username, err)
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
		logger.DEBUG.Printf("%v logged added to accounts", creds.Username)
		return nil
	case <-time.After(timeOutDuration):
		logger.ERROR.Printf("%v login timeout", creds.Username)
		return TimeOut
	}
}

func (a *accounts) InspectWeapon(params gt.InspectParams) (gt.Item, error) {
	clientIdx := a.getNextFreeAccount()
	if clientIdx < 0 {
		logger.ERROR.Print("There are no avaliable account")
		return gt.Item{}, NoAvaliableAccounts
	}

	item, err := a.validateDBItem(params.ParamA)
	if err == nil {
		logger.DEBUG.Printf("Valid item is in DB. ItemID: %v", params.ParamA)
		return item, nil
	}
	if !errors.Is(err, staleItem) && !errors.Is(err, storage.NoItem) {
		logger.ERROR.Printf("Error getting item from DB %v. Err: %v", params.ParamA, err)
		return gt.Item{}, InternalError
	}

	curClient := a.clients[clientIdx]

	pip := &csProto.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest{
		ParamM: proto.Uint64(uint64(params.ParamM)),
		ParamA: proto.Uint64(uint64(params.ParamA)),
		ParamD: proto.Uint64(uint64(params.ParamD)),
		ParamS: proto.Uint64(uint64(params.ParamS)),
	}

	msg := gamecoordinator.NewGCMsgProtobuf(730, 9156, pip)

	ch := make(chan error, 1)

	a.reqHandler.AddRequest(int(params.ParamA), ch)

	curClient.client.GC.Write(msg)
	logger.DEBUG.Printf("%v is requesting skin from GC, ItemID: %v", curClient.username, params.ParamA)

	select {
	case err := <-ch:
		if err != nil {
			logger.ERROR.Printf("Error with getting response for %v. Err: %v", params.ParamA, err)
			return gt.Item{}, InternalError
		}
	case <-time.After(timeOutDuration):
		logger.ERROR.Printf("No response for %v in response window", params.ParamA)
		return gt.Item{}, TimeOut
	}

	return a.db.GetItem(params.ParamA)
}

func NewAccountsList(handler steam.GCPacketHandler, reqHandler req.RequestHandler, store storage.Storage) Accounts {
	if handler == nil {
		panic("No handler func for account")
	}
	return &accounts{handler: handler, reqHandler: reqHandler, db: store}
}

func (a *accounts) getNextFreeAccount() int {
	for i, c := range a.clients {

		if time.Now().After(c.lastUsed.Add(1 * time.Second)) {
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
				lastUsed:  time.Now().Add(-time.Second * 5),
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

func (a accounts) validateDBItem(itemID int64) (gt.Item, error) {
	item, err := a.db.GetItem(itemID)
	if err != nil {
		return item, err
	}
	if !itemNewEnough(item) {
		err = a.db.DeleteItem(int64(item.ItemID))
		if err != nil {
			return item, err
		}

		return item, staleItem
	}

	return item, nil
}

func itemNewEnough(item gt.Item) bool {
	week := time.Hour * 24 * 7
	if time.Since(item.LastModified) < week {
		return false
	}

	return true
}
