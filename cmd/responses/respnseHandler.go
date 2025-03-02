package responses

import (
	"fmt"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/volodymyrzuyev/goCsInspect/cmd/globalTypes"
	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
	req "github.com/volodymyrzuyev/goCsInspect/cmd/requests"
	"github.com/volodymyrzuyev/goCsInspect/cmd/storage"
	"google.golang.org/protobuf/proto"
)

type ResponseHandler interface {
	HandleGCPacket(pack *gamecoordinator.GCPacket)
}

type responseHandler struct {
	reqHandler req.RequestHandler
	db         storage.Storage
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

		err = r.db.InsertItem(globalTypes.Item{ItemID: int(*msg.GetIteminfo().Itemid)})
		if err != nil {
			fmt.Println(err)
			panic("DB error")
		} else {
			fmt.Println("NO err")
		}

		r.reqHandler.FinishRequest(int(*msg.GetIteminfo().Itemid))
	}
}

func NewReponseHandler(reqHandler req.RequestHandler, store storage.Storage) ResponseHandler {
	return responseHandler{reqHandler, store}
}
