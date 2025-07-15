package clientmanagement

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/internal/client"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/creds"
	"github.com/volodymyrzuyev/goCsInspect/pkg/detailer"
	inspect "github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
	"github.com/volodymyrzuyev/goCsInspect/pkg/item"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage"
)

type ClientManager interface {
	AddClient(credentials creds.Credentials) error
	InspectSkin(params inspect.Parameters) (*item.Item, error)
	InspectSkinWithCtx(ctx context.Context, params inspect.Parameters) (*item.Item, error)
	GetProto(params inspect.Parameters) (*protobuf.CEconItemPreviewDataBlock, error)
	GetProtoWithCtx(
		ctx context.Context,
		params inspect.Parameters,
	) (*protobuf.CEconItemPreviewDataBlock, error)
}

type clientManager struct {
	clientCooldown time.Duration
	requestTTl     time.Duration

	detailer detailer.Detailer
	storage  storage.Storage
	l        *slog.Logger

	gcHandler gcHandler.GcHandler
	clientQue *clientQue
	jobQue    *jobQue
}

func NewClientManager(
	requestTTL time.Duration,
	clientCooldown time.Duration,
	detailer detailer.Detailer,
	storage storage.Storage,
) (ClientManager, error) {

	lcm := slog.Default().WithGroup("ClientManagment")

	if detailer == nil || storage == nil {
		lcm.Error("Detailler and storage are needed for ClientManager",
			"detailer_exits", detailer != nil, "storage_exists", storage != nil)
		return nil, errors.ErrInvalidManagerConfig
	}

	gcHandler := gcHandler.NewGcHandler()
	clientList := newClientQue(clientCooldown)
	jobQue := newJobQue(clientList)

	return &clientManager{
		clientCooldown: clientCooldown,
		requestTTl:     requestTTL,

		l:        lcm,
		storage:  storage,
		detailer: detailer,

		gcHandler: gcHandler,
		clientQue: clientList,
		jobQue:    jobQue,
	}, nil
}

func (c *clientManager) AddClient(credentials creds.Credentials) error {
	newClient, err := client.NewInspectClient(credentials, c.clientCooldown, c.gcHandler)
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

func (c *clientManager) GetProtoWithCtx(
	ctx context.Context,
	params inspect.Parameters,
) (*protobuf.CEconItemPreviewDataBlock, error) {

	ctx, cancel := context.WithTimeout(ctx, c.requestTTl)
	defer cancel()

	proto, err := c.storage.GetItem(ctx, params)
	if err == nil {
		c.l.Debug("item previously stored", "params", fmt.Sprintf("%+v", params))
		return proto, nil
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
		return resp.responseProto, nil

	case <-ctx.Done():
		slog.Error("client timeout requesting item preview block", "item_id", params.A)
		return nil, errors.ErrClientTimeout
	}
}

func (c *clientManager) GetProto(
	params inspect.Parameters,
) (*protobuf.CEconItemPreviewDataBlock, error) {

	return c.GetProtoWithCtx(context.TODO(), params)
}

func (c *clientManager) InspectSkinWithCtx(
	ctx context.Context,
	params inspect.Parameters,
) (*item.Item, error) {

	proto, err := c.GetProtoWithCtx(ctx, params)

	if err != nil {
		return nil, err
	}

	return c.detailer.DetailProto(proto)
}

func (c *clientManager) InspectSkin(params inspect.Parameters) (*item.Item, error) {
	ctx := context.TODO()
	return c.InspectSkinWithCtx(ctx, params)
}

func (c *clientManager) storeToStorage(
	ctx context.Context,
	params inspect.Parameters,
	proto *protobuf.CEconItemPreviewDataBlock,
) {
	err := c.storage.StoreItem(ctx, params, proto)
	if err != nil {
		c.l.Error(
			"item not stored",
			"inspect_params", fmt.Sprintf("%+v", params),
			"error", err,
		)
	}
}
