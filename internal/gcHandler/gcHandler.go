package gcHandler

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/config"
	"google.golang.org/protobuf/proto"
)

type GcHandler interface {
	HandleGCPacket(packet *gamecoordinator.GCPacket)
	GetResponse(itemId uint64) (*csProto.CEconItemPreviewDataBlock, error)
}

type gcHandler struct {
	mu               sync.Mutex
	responses        map[uint64]*csProto.CEconItemPreviewDataBlock
	pendingResponses map[uint64]chan *csProto.CEconItemPreviewDataBlock
}

func NewGcHandler() GcHandler {
	return &gcHandler{
		pendingResponses: make(map[uint64]chan *csProto.CEconItemPreviewDataBlock),
		responses:        make(map[uint64]*csProto.CEconItemPreviewDataBlock),
	}
}

const (
	csAppID                = 730
	inspectResponseProtoID = 9157
)

func (r *gcHandler) storeResponse(itemInfo *csProto.CEconItemPreviewDataBlock) {
	itemId := itemInfo.GetItemid()

	r.mu.Lock()
	defer r.mu.Unlock()

	// is waiting for response, send it down chanel
	if ch, ok := r.pendingResponses[itemId]; ok {
		ch <- itemInfo
		close(ch)
		delete(r.pendingResponses, itemId)
		slog.Debug("Sending skin fetching response to in-flight recipient",
			"item_id", itemId)
	} else {
		// update response map
		slog.Debug("Storing response for later reply", "item_id", itemId)
		r.responses[itemId] = itemInfo
	}
}

func (r *gcHandler) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	// vertify right type of packet
	if packet.AppId != csAppID || packet.MsgType != inspectResponseProtoID {
		slog.Debug("Got non CS inpsect packet",
			"packet_AppId", packet.AppId, "packet_MsgType",
			packet.MsgType, "body", packet.Body)
		return
	}

	// unmarshal packet
	var msg protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse
	err := proto.Unmarshal(packet.Body, &msg)
	if err != nil {
		slog.Debug("Error unmarshaling proto",
			"packet_AppId", packet.AppId, "packet_MsgType", packet.MsgType,
			"body", packet.Body, "error", err.Error())
		return
	}

	slog.Debug("Got skin fetching response",
		"packet_AppId", packet.AppId,
		"packet_MsgType", packet.MsgType, "item_id", msg.Iteminfo.GetItemid())

	r.storeResponse(msg.Iteminfo)
}

func (r *gcHandler) GetResponse(itemId uint64) (*csProto.CEconItemPreviewDataBlock, error) {
	r.mu.Lock()
	slog.Debug("Got request for fetched skin details", "item_id", itemId)
	// if response is in the response map, return it
	if response, ok := r.responses[itemId]; ok {
		delete(r.responses, itemId)
		r.mu.Unlock()
		slog.Debug("Request for fetched skin details already stored",
			"item_id", itemId)
		return response, nil
	}

	// else make a chanel for the response to come to
	ch := make(chan *csProto.CEconItemPreviewDataBlock)
	r.pendingResponses[itemId] = ch
	// unlock mutex
	r.mu.Unlock()
	slog.Debug("Waiting for an in-flight fetched skin details", "item_id", itemId)

	// wait for response, or timout
	select {
	case response := <-ch:
		slog.Debug("Get in-flight response", "item_id", itemId)
		return response, nil
	case <-time.After(config.TimeOutDuration):
		// if no response in allowed time, clean up the chan map
		r.mu.Lock()
		if chP, ok := r.pendingResponses[itemId]; ok && chP == ch {
			delete(r.pendingResponses, itemId)
		}
		r.mu.Unlock()
		slog.Error("Request for fetched skin details timed out", "item_id", itemId)
		return nil, errors.ErrClientTimeout
	}
}
