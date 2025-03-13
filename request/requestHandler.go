package request

import "github.com/volodymyrzuyev/goCsInspect/types"

type RequestHandler interface {
	RegisterRequest(assetID uint64) *chan types.Response
}
