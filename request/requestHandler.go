package request

import (
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

type RequestHandler interface {
	RegisterRequest(assetID uint64) *chan types.Response
	ResolverRequest(*protobuf.CEconItemPreviewDataBlock, error)
}
