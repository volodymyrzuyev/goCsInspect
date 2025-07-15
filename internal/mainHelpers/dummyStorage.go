package mainhelpers

import (
	"context"
	"errors"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
)

type dummyStorage struct{}

func (d *dummyStorage) StoreItem(
	ctx context.Context,
	inspectParams inspect.Params,
	proto *protobuf.CEconItemPreviewDataBlock,
) error {
	return nil
}

func (d *dummyStorage) GetItem(
	ctx context.Context,
	inspectParams inspect.Params,
) (*protobuf.CEconItemPreviewDataBlock, error) {
	return nil, errors.New("I am dummy")
}
