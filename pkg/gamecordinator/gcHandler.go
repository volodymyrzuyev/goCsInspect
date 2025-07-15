package gamecordinator

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

type Handler interface {
	// Used to handle all packets from the steam game game coordinator
	HandleGCPacket(packet *gamecoordinator.GCPacket)
	// Returns protobuf for item with a specified itemid, timeout can be set using ctx
	GetResponse(
		ctx context.Context,
		itemId uint64,
	) (*csProto.CEconItemPreviewDataBlock, error)
}

type handler struct {
	mu               sync.Mutex
	responses        map[uint64]*csProto.CEconItemPreviewDataBlock
	pendingResponses map[uint64]chan *csProto.CEconItemPreviewDataBlock

	l *slog.Logger
}

func NewGcHandler() Handler {
	return &handler{
		pendingResponses: make(map[uint64]chan *csProto.CEconItemPreviewDataBlock),
		responses:        make(map[uint64]*csProto.CEconItemPreviewDataBlock),

		l: slog.Default().WithGroup("GcHandler"),
	}
}

const (
	csAppID                = 730
	inspectResponseProtoID = 9157
)

func (h *handler) storeResponse(itemInfo *csProto.CEconItemPreviewDataBlock) {
	itemId := itemInfo.GetItemid()

	h.mu.Lock()
	defer h.mu.Unlock()

	// is waiting for response, send it down chanel
	if ch, ok := h.pendingResponses[itemId]; ok {
		ch <- itemInfo
		close(ch)
		delete(h.pendingResponses, itemId)
	} else {
		// update response map
		h.responses[itemId] = itemInfo
	}
}

func (h *handler) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	// vertify right type of packet
	if packet.AppId != csAppID || packet.MsgType != inspectResponseProtoID {
		return
	}

	// unmarshal packet
	var msg protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse
	err := proto.Unmarshal(packet.Body, &msg)
	if err != nil {
		h.l.Debug(
			"could not unmarshal gcPacket",
			"packet_AppId", packet.AppId,
			"packet_MsgType", packet.MsgType,
			"error", err.Error(),
		)
		return
	}

	h.l.Debug(
		"got item preview block",
		"packet_AppId", packet.AppId,
		"packet_MsgType", packet.MsgType,
		"item_id", msg.Iteminfo.GetItemid(),
	)

	h.storeResponse(msg.Iteminfo)
}

func (h *handler) GetResponse(
	ctx context.Context,
	itemId uint64,
) (*csProto.CEconItemPreviewDataBlock, error) {

	h.mu.Lock()
	h.l.Debug("got request for item preview block", "item_id", itemId)

	// if response is in the response map, return it
	if response, ok := h.responses[itemId]; ok {
		delete(h.responses, itemId)
		h.mu.Unlock()
		h.l.Debug("item preview block previously received", "item_id", itemId)
		return response, nil
	}

	// else make a chanel for the response to come to
	ch := make(chan *csProto.CEconItemPreviewDataBlock)
	h.pendingResponses[itemId] = ch
	// unlock mutex
	h.mu.Unlock()
	h.l.Debug("waiting for preview block from gc", "item_id", itemId)

	// wait for response, or timout
	select {
	case response := <-ch:
		h.l.Debug("got item preview block", "item_id", itemId)
		return response, nil

	case <-ctx.Done():
		// if no response in allowed time, clean up the chan map
		h.mu.Lock()
		if chP, ok := h.pendingResponses[itemId]; ok && chP == ch {
			delete(h.pendingResponses, itemId)
		}
		h.mu.Unlock()
		h.l.Error("item preview block request timed out", "item_id", itemId)

		return nil, errors.ErrClientTimeout
	}
}
