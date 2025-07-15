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

type Manager interface {
	AddClient(credentials creds.Credentials) error
	InspectSkin(params inspect.Params) (*item.Item, error)
	InspectSkinWithCtx(ctx context.Context, params inspect.Params) (*item.Item, error)
	GetProto(params inspect.Params) (*protobuf.CEconItemPreviewDataBlock, error)
	GetProtoWithCtx(
		ctx context.Context,
		params inspect.Params,
	) (*protobuf.CEconItemPreviewDataBlock, error)
}

type manager struct {
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
) (Manager, error) {

	lcm := slog.Default().WithGroup("ClientManagment")

	if detailer == nil || storage == nil {
		lcm.Error("Detailler and storage are needed for ClientManager",
			"detailer_exits", detailer != nil, "storage_exists", storage != nil)
		return nil, errors.ErrInvalidManagerConfig
	}

	gcHandler := gcHandler.NewGcHandler()
	clientList := newClientQue(clientCooldown)
	jobQue := newJobQue(clientList)

	return &manager{
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

func (m *manager) AddClient(credentials creds.Credentials) error {
	newClient, err := client.NewInspectClient(credentials, m.clientCooldown, m.gcHandler)
	if err != nil {
		return err
	}

	err = newClient.LogIn()
	if err != nil {
		return err
	}

	m.clientQue.addClient(newClient)
	return nil
}

func (m *manager) GetProtoWithCtx(
	ctx context.Context,
	params inspect.Params,
) (*protobuf.CEconItemPreviewDataBlock, error) {

	ctx, cancel := context.WithTimeout(ctx, m.requestTTl)
	defer cancel()

	proto, err := m.storage.GetItem(ctx, params)
	if err == nil {
		m.l.Debug("item previously stored", "params", fmt.Sprintf("%+v", params))
		return proto, nil
	}

	if m.clientQue.len() == 0 {
		m.l.Error("no avaliable clients to request data")
		return nil, errors.ErrNoAvailableClients
	}

	inspectProto, err := params.GenerateGcRequestProto()
	if err != nil {
		m.l.Error("invalid request params", "params", fmt.Sprintf("%+v", params))
		return nil, err
	}

	responseChan := m.jobQue.registerJob(inspectProto, ctx)

	select {
	case resp := <-responseChan:
		if resp.err != nil {
			return nil, err
		}
		go m.storeToStorage(context.Background(), params, resp.responseProto)
		return resp.responseProto, nil

	case <-ctx.Done():
		slog.Error("client timeout requesting item preview block", "item_id", params.A)
		return nil, errors.ErrClientTimeout
	}
}

func (m *manager) GetProto(
	params inspect.Params,
) (*protobuf.CEconItemPreviewDataBlock, error) {

	return m.GetProtoWithCtx(context.TODO(), params)
}

func (m *manager) InspectSkinWithCtx(
	ctx context.Context,
	params inspect.Params,
) (*item.Item, error) {

	proto, err := m.GetProtoWithCtx(ctx, params)

	if err != nil {
		return nil, err
	}

	return m.detailer.DetailProto(proto)
}

func (m *manager) InspectSkin(params inspect.Params) (*item.Item, error) {
	ctx := context.TODO()
	return m.InspectSkinWithCtx(ctx, params)
}

func (m *manager) storeToStorage(
	ctx context.Context,
	params inspect.Params,
	proto *protobuf.CEconItemPreviewDataBlock,
) {
	err := m.storage.StoreItem(ctx, params, proto)
	if err != nil {
		m.l.Error(
			"item not stored",
			"inspect_params", fmt.Sprintf("%+v", params),
			"error", err,
		)
	}
}
