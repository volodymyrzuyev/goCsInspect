package clientmanagement

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
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
	gcHandler gcHandler.GcHandler
	clientQue *clientQue
	jobQue    *jobQue

	storage      storage.Storage
	clientConfig config.ClientConfig
	detailer     detailer.Detailer
	l            *slog.Logger
}

func NewClientManager(detailer detailer.Detailer, clientConfig config.ClientConfig, storage storage.Storage, l *slog.Logger) (*ClientManager, error) {
	lcm := l.WithGroup("ClientManagment")

	if detailer == nil || storage == nil {
		l.Error("Detailler and storage are needed for ClientManager",
			"detailer_exits", detailer != nil, "storage_exists", storage != nil)
		return nil, errors.ErrInvalidManagerConfig
	}

	gcHandler := gcHandler.NewGcHandler(l)
	clientList := newClientQue(clientConfig.TimeOutDuration, l)
	jobQue := newJobQue(clientList, l)

	return &ClientManager{
		gcHandler: gcHandler,
		clientQue: clientList,
		jobQue:    jobQue,

		storage:      storage,
		clientConfig: clientConfig,
		detailer:     detailer,
		l:            lcm,
	}, nil
}

func (c *ClientManager) AddClient(credentials types.Credentials) error {
	newClient, err := client.NewInspectClient(c.clientConfig, c.gcHandler, credentials, c.l)
	if err != nil {
		return err
	}

	err = newClient.LogIn()
	if err != nil {
		return err
	}

	c.clientQue.addClient(newClient)
	return nil
}

func (c *ClientManager) InspectSkin(params types.InspectParameters) (*item.Item, error) {
	ctx := context.TODO()
	return c.InspectSkinWithCtx(ctx, params)
}

func (c *ClientManager) InspectSkinWithCtx(ctx context.Context, params types.InspectParameters) (*item.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, c.clientConfig.TimeOutDuration)
	defer cancel()

	proto, err := c.storage.GetItem(ctx, params)
	if err == nil {
		c.l.Debug("item previously stored", "params", fmt.Sprintf("%+v", params))
		return c.detailer.DetailProto(proto)
	}

	if c.clientQue.len() == 0 {
		c.l.Error("no avaliable clients to request data")
		return nil, errors.ErrNoAvailableClients
	}

	inspectProto, err := params.GenerateGcRequestProto()
	if err != nil {
		c.l.Error("invalid request params", "params", fmt.Sprintf("%+v", params))
		return nil, err
	}

	responseChan := c.jobQue.registerJob(inspectProto, ctx)

	select {
	case resp := <-responseChan:
		if resp.err != nil {
			return nil, err
		}
		go c.storeToStorage(context.Background(), params, resp.responseProto)
		return c.detailer.DetailProto(resp.responseProto)

	case <-ctx.Done():
		slog.Error("client timeout requesting item preview block", "item_id", params.A)
		return nil, errors.ErrClientTimeout
	}
}

func (c *ClientManager) storeToStorage(ctx context.Context, params types.InspectParameters, proto *protobuf.CEconItemPreviewDataBlock) {
	err := c.storage.StoreItem(ctx, params, proto)
	if err != nil {
		c.l.Error(
			"item not stored",
			"inspect_params", fmt.Sprintf("%+v", params),
			"error", err,
		)
	}
}
