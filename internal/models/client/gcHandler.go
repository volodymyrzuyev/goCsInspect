package client

import (
	"log/slog"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/common/consts"
	"github.com/volodymyrzuyev/goCsInspect/common/types"
	"google.golang.org/protobuf/proto"
)

type GcHandler interface {
	HandleGCPacket(*gamecoordinator.GCPacket)
}

type gcHandler struct {
	responseChan chan types.Response
	username     string
}

func NewGcHandler(r chan types.Response, username string) GcHandler {
	return gcHandler{responseChan: r, username: username}
}

func (r gcHandler) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	if packet.AppId != consts.CsAppID || packet.MsgType != consts.InspectResponseProtoID {
		slog.Debug("Got non CS inpsect packet", "username", r.username, "packet_AppId", packet.AppId, "packet_MsgType", packet.MsgType)
		return
	}

	var msg protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse
	err := proto.Unmarshal(packet.Body, &msg)
	if err != nil {
		slog.Debug("Error unmarshaling proto", "username", r.username, "packet_AppId", packet.AppId, "packet_MsgType", packet.MsgType, "error", err.Error())
		r.responseChan <- types.Response{Response: nil, Error: err}
		return
	}
	r.responseChan <- types.Response{Response: msg.GetIteminfo(), Error: err}
}
