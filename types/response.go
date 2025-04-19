package types

import csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"

type Response struct {
	Response *csProto.CEconItemPreviewDataBlock
	Error    error
}
