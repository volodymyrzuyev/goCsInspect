package client

import (
	"context"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
)

type Client interface {
	IsLoggedIn() bool
	IsAvailable() bool
	Username() string

	LogOff()
	LogIn() error
	Reconnect() error

	InspectItem(
		ctx context.Context,
		params *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest,
	) (*protobuf.CEconItemPreviewDataBlock, error)
}
