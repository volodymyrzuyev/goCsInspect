package sqlite

import (
	"log/slog"
	"testing"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage/sqlite/sql/sqlc"
	"github.com/volodymyrzuyev/goCsInspect/tests/testdata"
)

func getModDBItem(protos []*protobuf.CEconItemPreviewDataBlock_Sticker) []sqlc.Mod {
	mods := make([]sqlc.Mod, len(protos))

	for i, proto := range protos {
		mods[i] = sqlc.Mod{
			Slot:          common.NullInt64Uint32Ptr(proto.Slot),
			Stickerid:     common.NullInt64Uint32Ptr(proto.StickerId),
			Wear:          common.NullFloat64Float32Ptr(proto.Wear),
			Scale:         common.NullFloat64Float32Ptr(proto.Scale),
			Rotation:      common.NullFloat64Float32Ptr(proto.Rotation),
			Tintid:        common.NullInt64Uint32Ptr(proto.TintId),
			Offsetx:       common.NullFloat64Float32Ptr(proto.OffsetX),
			Offsety:       common.NullFloat64Float32Ptr(proto.OffsetY),
			Offsetz:       common.NullFloat64Float32Ptr(proto.OffsetZ),
			Pattern:       common.NullInt64Uint32Ptr(proto.Pattern),
			Highlightreel: common.NullInt64Uint32Ptr(proto.HighlightReel),
		}
	}

	return mods
}

func TestFetch(t *testing.T) {
	protos := testdata.GetResponseProtos()
	params := testdata.GetInspectParams()

	s := Sqlite{l: slog.New(slog.DiscardHandler)}

	for name, proto := range protos {

		t.Run(name, func(t *testing.T) {
			param, ok := params[name]
			if !ok {
				t.Fatalf("No inspect link for proto name:%v", name)
			}

			dbItem := getDBProto(param, proto)
			stickers := getModDBItem(proto.GetStickers())
			chains := getModDBItem(proto.GetKeychains())

			newProto, err := s.assembleItem(sqlc.Item(dbItem), chains, stickers)

			assert.Nil(t, nil, "should be no err")
			if err != nil {
				return
			}
			assert.Equal(t, proto, newProto, "should be same")
		})
	}

}
