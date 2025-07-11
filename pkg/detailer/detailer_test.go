package detailer

import (
	"fmt"
	"testing"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
	"github.com/volodymyrzuyev/goCsInspect/pkg/item"
)

type protoTestCase struct {
	input         *protobuf.CEconItemPreviewDataBlock
	expectedItem  *item.Item
	expectedError error
}

func TestDetailSkin(t *testing.T) {
	detailer := NewDetailer(config.GetEnglishFile(), config.GetGameItems())

	tests := getTestCases()

	for testName, input := range tests {
		t.Run(testName, func(t *testing.T) {
			respItem, respErr := detailer.DetailProto(input.input)

			assert.Equal(t, respErr, input.expectedError, "errors should be the same. Resp error %v", respErr)

			if respErr != nil {
				return
			}

			// proto comparison
			assert.Equal(t, input.expectedItem.Accountid, respItem.Accountid, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Itemid, respItem.Itemid, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Defindex, respItem.Defindex, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Paintindex, respItem.Paintindex, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Rarity, respItem.Rarity, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Quality, respItem.Quality, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Paintwear, respItem.Paintwear, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Paintseed, respItem.Paintseed, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Killeaterscoretype, respItem.Killeaterscoretype, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Killeatervalue, respItem.Killeatervalue, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Customname, respItem.Customname, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Inventory, respItem.Inventory, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Origin, respItem.Origin, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Questid, respItem.Questid, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Dropreason, respItem.Dropreason, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Musicindex, respItem.Musicindex, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Entindex, respItem.Entindex, "Proto should be the same")
			assert.Equal(t, input.expectedItem.Petindex, respItem.Petindex, "Proto should be the same")

			assert.Equal(t, fmt.Sprintf("%.15f", input.expectedItem.FloatValue), fmt.Sprintf("%.15f", respItem.FloatValue), "Float values should be same")
			assert.Equal(t, input.expectedItem.MinFloat, respItem.MinFloat, "MinFloat should be the same")
			assert.Equal(t, input.expectedItem.MaxFloat, respItem.MaxFloat, "MaxFloat should be the same")
			assert.Equal(t, input.expectedItem.ItemName, respItem.ItemName, "ItemName should be the same")
			assert.Equal(t, input.expectedItem.QualityName, respItem.QualityName, "QualityName should be the same")
			assert.Equal(t, input.expectedItem.WeaponType, respItem.WeaponType, "WeaponType should be the same")
			assert.Equal(t, input.expectedItem.RarityName, respItem.RarityName, "RarityName should be the same")
			assert.Equal(t, input.expectedItem.WearName, respItem.WearName, "WearName should be the same")
			assert.Equal(t, input.expectedItem.FullItemName, respItem.FullItemName, "FullItemName should be the same")
			assert.Equal(t, input.expectedItem.Stickers, respItem.Stickers, "Stickers should be the same")
			assert.Equal(t, input.expectedItem.Keychains, respItem.Keychains, "Keychains should be the same")
		})
	}

}

var (
	uint64Pointer  = common.Uint64Pointer
	uint32Pointer  = common.Uint32Pointer
	float32Pointer = common.Float32Pointer
	float64Pointer = common.Float64Pointer
	stringPointer  = common.StringPointer
	int32Pointer   = common.Int32Pointer
)
