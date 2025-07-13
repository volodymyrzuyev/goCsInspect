package storage

import (
	"context"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/pkg/inspectParams"
)

type Storage interface {
	StoreItem(
		ctx context.Context,
		inspectParams inspectParams.InspectParameters,
		proto *protobuf.CEconItemPreviewDataBlock,
	) error
	GetItem(
		ctx context.Context,
		inspectParams inspectParams.InspectParameters,
	) (*protobuf.CEconItemPreviewDataBlock, error)
}
