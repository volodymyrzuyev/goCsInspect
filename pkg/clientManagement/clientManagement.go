package clientManagement

import (
	"context"
	"log/slog"

	"github.com/volodymyrzuyev/goCsInspect/config"
	"github.com/volodymyrzuyev/goCsInspect/internal/client"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
	"github.com/volodymyrzuyev/goCsInspect/pkg/item"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage"
	"github.com/volodymyrzuyev/goCsInspect/pkg/types"
)

type clientList struct {
	clients  []client.InspectClient
	lastUsed int
}

func (c *clientList) getNextClient() client.InspectClient {
	c.lastUsed++
	c.lastUsed = c.lastUsed % len(c.clients)
	return c.clients[c.lastUsed]
}

type ClientManager struct {
	gcHandler  gcHandler.GcHandler
	clientList *clientList

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

	return &ClientManager{
		gcHandler:  gcHandler.NewGcHandler(clientConfig.TimeOutDuration),
		clientList: &clientList{clients: make([]client.InspectClient, 0)},

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

	c.clientList.clients = append(c.clientList.clients, newClient)

	return nil
}

func (c *ClientManager) InspectSkin(params types.InspectParameters) (*item.Item, error) {
	proto, err := c.storage.GetItem(context.TODO(), params)
	if err == nil {
		slog.Debug("item stored in db", "params", params)
		return c.detailer.DetailProto(proto)
	}

	var curClient client.InspectClient
	for range len(c.clientList.clients) {
		curClient = c.clientList.getNextClient()
		if !curClient.IsLoggedIn() {
			curClient.Reconnect()
		}

		if curClient.IsAvailable() {
			break
		}

		curClient = nil
	}
	if curClient == nil {
		slog.Error("No available clients")
		return nil, errors.ErrNoAvailableClients
	}

	inspectProto, err := params.GenerateGcRequestProto()
	if err != nil {
		slog.Error("Invalid inspect parameters", "inspect_parameters", params)
		return nil, err
	}

	responseProto, err := curClient.InspectItem(inspectProto)
	if err != nil {
		return nil, err
	}
	c.storage.StoreItem(context.TODO(), params, responseProto)

	return c.detailer.DetailProto(responseProto)
}
