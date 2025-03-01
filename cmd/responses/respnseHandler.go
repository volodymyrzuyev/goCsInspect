package responses

import (
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
	"github.com/volodymyrzuyev/goCsInspect/cmd/requests"
	"google.golang.org/protobuf/proto"
)

type ResponseHandler interface {
	HandleGCPacket(pack *gamecoordinator.GCPacket)
}

type responseHandler struct {
	reqHandler requests.RequestHandler
}

func (r responseHandler) HandleGCPacket(pack *gamecoordinator.GCPacket) {
	logger.DEBUG.Printf("Got GC packet AppId: %v MsgType: %v", pack.AppId, pack.MsgType)
	switch pack.AppId {
	case 730:
		r.csResponse(pack)
	}
}

func (r responseHandler) csResponse(pack *gamecoordinator.GCPacket) {
	switch pack.MsgType {
	case 9157:
		var msg protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse
		err := proto.Unmarshal(pack.Body, &msg)
		if err != nil {
			return
		}

		logger.INFO.Printf("Got resp for %v", int(*msg.GetIteminfo().Itemid))
		r.reqHandler.FinishRequest(int(*msg.GetIteminfo().Itemid))
	}
}

func NewReponseHandler(reqHandler requests.RequestHandler) ResponseHandler {
	return responseHandler{reqHandler}
}
