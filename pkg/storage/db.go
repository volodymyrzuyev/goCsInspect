package storage

import (
	"context"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
)

type Storage interface {
	// stores an item protobuf
	StoreItem(
		ctx context.Context,
		inspectParams inspect.Params,
		proto *protobuf.CEconItemPreviewDataBlock,
	) error
	// fetches an item protobuf
	GetItem(
		ctx context.Context,
		inspectParams inspect.Params,
	) (*protobuf.CEconItemPreviewDataBlock, error)
}
