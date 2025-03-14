package response

import (
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/request"
	"github.com/volodymyrzuyev/goCsInspect/types"
	"google.golang.org/protobuf/proto"
)

type ResponseHandler interface {
	HandleGCPacket(*gamecoordinator.GCPacket)
}

type responseHandler struct {
	log        logger.Logger
	reqHandler request.RequestHandler
}

func (r responseHandler) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	r.log.Debug("Processing response for appid: %v", packet.AppId)
	if packet.AppId != types.CsGameID {
		return
	}

	if packet.MsgType != types.InspectResponseProtoID {
		r.log.Debug("Packet is not inspect response packet")
		return
	}

	var msg protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse
	if err := proto.Unmarshal(packet.Body, &msg); err != nil {
		r.log.Debug("Error decoding message")
		r.reqHandler.ResolverRequest(nil, err)
		return
	}

	r.log.Debug("Successfully got response for %v", msg.GetIteminfo().Itemid)
	r.reqHandler.ResolverRequest(msg.GetIteminfo(), nil)
}
