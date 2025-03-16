package accounts

import (
	"errors"

	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/accounts/client"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

type ClientManager interface {
	AddClient(creds types.Credentials) error
	InspectSkin(params types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error)
}

type clientManager struct {
	curIdx     int
	clientList []client.Client
	log        logger.Logger
}

func (c *clientManager) AddClient(creds types.Credentials) error {
	cli := client.NewClient(c.log)
	err := cli.LogIn(creds)
	if err != nil {
		return err
	}
	c.log.Debug("Adding Client: %v to client list", creds.Username)
	c.clientList = append(c.clientList, cli)

	return nil
}

var NoValidAccount = errors.New("No valid accounts to send request")

func (c *clientManager) InspectSkin(params types.InspectParameters) (*csProto.CEconItemPreviewDataBlock, error) {
	var data *csProto.CEconItemPreviewDataBlock
	var err error

	for range c.clientList {
		curClient := c.clientList[c.getNextClient()]
		data, err = curClient.RequestSkin(params)

		if err != nil {
			c.log.Debug("Client: %v got error \"%v\" when inspecting %v", curClient.Username(), err, params.A)

			if errors.Is(err, client.ClientNotAvaliable) || errors.Is(err, client.TimeoutPassed) {
				continue
			}

			return nil, err
		}

		break
	}

	if err != nil || data == nil {
		c.log.Error("No account to inspect %v", params.A)
		return nil, NoValidAccount
	}

	c.log.Error("Got item information for %v", params.A)
	return data, nil

}

func (c *clientManager) getNextClient() int {
	c.curIdx += 1
	return c.curIdx % len(c.clientList)
}

func NewClientManager(l logger.Logger) ClientManager {
	if l == nil {
		panic("Logger is requeried for Client manager")
	}

	return &clientManager{log: l}
}
