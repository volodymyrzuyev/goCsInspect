package storage

import (
	"context"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/pkg/types"
)

type Storage interface {
	StoreItem(ctx context.Context, inspectParams types.InspectParameters, proto *protobuf.CEconItemPreviewDataBlock) error
	GetItem(ctx context.Context, inspectParams types.InspectParameters) (*protobuf.CEconItemPreviewDataBlock, error)
}
