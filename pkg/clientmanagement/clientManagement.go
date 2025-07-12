package clientmanagement

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/volodymyrzuyev/goCsInspect/config"
	"github.com/volodymyrzuyev/goCsInspect/internal/client"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
	"github.com/volodymyrzuyev/goCsInspect/pkg/item"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage"
	"github.com/volodymyrzuyev/goCsInspect/pkg/types"
)

type ClientManager struct {
	gcHandler  gcHandler.GcHandler
	clientList *clientQue
	jobQue     *jobQue

	storage      storage.Storage
	clientConfig config.ClientConfig
	detailer     detailer.Detailer
}

func NewClientManager(detailer detailer.Detailer, clientConfig config.ClientConfig, storage storage.Storage) (*ClientManager, error) {
	if detailer == nil || storage == nil {
		slog.Error("Detailler and storage are needed for ClientManager",
			"detailer_exits", detailer != nil, "storage_exists", storage != nil)
		return nil, errors.ErrInvalidManagerConfig
	}

	clientList := newClientQue(clientConfig.TimeOutDuration)
	return &ClientManager{
		gcHandler:  gcHandler.NewGcHandler(clientConfig.TimeOutDuration),
		clientList: clientList,
		jobQue:     newJobQue(clientList),

		storage:      storage,
		clientConfig: clientConfig,
		detailer:     detailer,
	}, nil
}

func (c *ClientManager) AddClient(credentials types.Credentials) error {
	newClient, err := client.NewInspectClient(c.clientConfig, c.gcHandler, credentials)
	if err != nil {
		return err
	}

	err = newClient.LogIn()
	if err != nil {
		return err
	}

	c.clientList.addClient(newClient)
	return nil
}

func (c *ClientManager) InspectSkin(params types.InspectParameters) (*item.Item, error) {
	proto, err := c.storage.GetItem(context.TODO(), params)
	if err == nil {
		slog.Debug("item stored in db", "params", params)
		return c.detailer.DetailProto(proto)
	}

	if c.clientList.len() == 0 {
		slog.Error("No available clients")
		return nil, errors.ErrNoAvailableClients
	}

	inspectProto, err := params.GenerateGcRequestProto()
	if err != nil {
		slog.Error("Invalid inspect parameters", "inspect_parameters", params)
		return nil, err
	}

	responseChan := c.jobQue.registerJob(inspectProto)
	select {
	case resp := <-responseChan:
		if resp.err != nil {
			return nil, err
		}
		err = c.storage.StoreItem(context.TODO(), params, resp.responseProto)
		if err != nil {
			slog.Error("Error string item", "inspect_params", fmt.Sprintf("%+v", params), "error", err)
		}
		return c.detailer.DetailProto(resp.responseProto)

	case <-time.After(c.clientList.clientCooldown * 25):
		return nil, errors.ErrClientTimeout
	}
}
