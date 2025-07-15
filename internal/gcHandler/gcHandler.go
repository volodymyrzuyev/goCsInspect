package gcHandler

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"google.golang.org/protobuf/proto"
)

type GcHandler interface {
	HandleGCPacket(packet *gamecoordinator.GCPacket)
	GetResponse(
		ctx context.Context,
		itemId uint64,
	) (*csProto.CEconItemPreviewDataBlock, error)
}

type gcHandler struct {
	mu               sync.Mutex
	responses        map[uint64]*csProto.CEconItemPreviewDataBlock
	pendingResponses map[uint64]chan *csProto.CEconItemPreviewDataBlock

	l *slog.Logger
}

func NewGcHandler() GcHandler {
	return &gcHandler{
		pendingResponses: make(map[uint64]chan *csProto.CEconItemPreviewDataBlock),
		responses:        make(map[uint64]*csProto.CEconItemPreviewDataBlock),

		l: slog.Default().WithGroup("GcHandler"),
	}
}

const (
	csAppID                = 730
	inspectResponseProtoID = 9157
)

func (g *gcHandler) storeResponse(itemInfo *csProto.CEconItemPreviewDataBlock) {
	itemId := itemInfo.GetItemid()

	g.mu.Lock()
	defer g.mu.Unlock()

	// is waiting for response, send it down chanel
	if ch, ok := g.pendingResponses[itemId]; ok {
		ch <- itemInfo
		close(ch)
		delete(g.pendingResponses, itemId)
	} else {
		// update response map
		g.responses[itemId] = itemInfo
	}
}

func (g *gcHandler) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	// vertify right type of packet
	if packet.AppId != csAppID || packet.MsgType != inspectResponseProtoID {
		return
	}

	// unmarshal packet
	var msg protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse
	err := proto.Unmarshal(packet.Body, &msg)
	if err != nil {
		g.l.Debug(
			"could not unmarshal gcPacket",
			"packet_AppId", packet.AppId,
			"packet_MsgType", packet.MsgType,
			"error", err.Error(),
		)
		return
	}

	g.l.Debug(
		"got item preview block",
		"packet_AppId", packet.AppId,
		"packet_MsgType", packet.MsgType,
		"item_id", msg.Iteminfo.GetItemid(),
	)

	g.storeResponse(msg.Iteminfo)
}

func (g *gcHandler) GetResponse(
	ctx context.Context,
	itemId uint64,
) (*csProto.CEconItemPreviewDataBlock, error) {

	g.mu.Lock()
	g.l.Debug("got request for item preview block", "item_id", itemId)

	// if response is in the response map, return it
	if response, ok := g.responses[itemId]; ok {
		delete(g.responses, itemId)
		g.mu.Unlock()
		g.l.Debug("item preview block previously received", "item_id", itemId)
		return response, nil
	}

	// else make a chanel for the response to come to
	ch := make(chan *csProto.CEconItemPreviewDataBlock)
	g.pendingResponses[itemId] = ch
	// unlock mutex
	g.mu.Unlock()
	g.l.Debug("waiting for preview block from gc", "item_id", itemId)

	// wait for response, or timout
	select {
	case response := <-ch:
		g.l.Debug("got item preview block", "item_id", itemId)
		return response, nil

	case <-ctx.Done():
		// if no response in allowed time, clean up the chan map
		g.mu.Lock()
		if chP, ok := g.pendingResponses[itemId]; ok && chP == ch {
			delete(g.pendingResponses, itemId)
		}
		g.mu.Unlock()
		g.l.Error("item preview block request timed out", "item_id", itemId)

		return nil, errors.ErrClientTimeout
	}
}
