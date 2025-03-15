package response

import (
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/types"
	"google.golang.org/protobuf/proto"
)

type ResponseHandler interface {
	HandleGCPacket(*gamecoordinator.GCPacket)
}

type responseHandler struct {
	log            logger.Logger
	responseChan   *chan types.Response
	clientUsername string
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
	err := proto.Unmarshal(packet.Body, &msg)
	*r.responseChan <- types.Response{Response: msg.GetIteminfo(), Error: err}
}

func NewResponseHandler(l logger.Logger, r *chan types.Response, clientName string) ResponseHandler {
	if l == nil || r == nil {
		panic("Logger and RequestHandler are needed for ResponseHandler")
	}
	l.Debug("Initializing responseHandler for Client %v", clientName)
	return responseHandler{log: l, responseChan: r, clientUsername: clientName}
}
