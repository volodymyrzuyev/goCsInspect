package detailer

import (
	"fmt"
	"testing"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/common/types"
	"github.com/volodymyrzuyev/goCsInspect/config"
)

type protoTestCase struct {
	input         *protobuf.CEconItemPreviewDataBlock
	expectedItem  *types.Item
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

			assert.Equal(t, input.expectedItem.Proto, respItem.Proto, "Proto should be the same")
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

func uint64Pointer(i uint64) *uint64 {
	return &i
}

func uint32Pointer(i uint32) *uint32 {
	return &i
}

func float32Pointer(f float32) *float32 {
	return &f
}

func float64Pointer(f float64) *float64 {
	return &f
}

func stringPointer(s string) *string {
	return &s
}

func int32Pointer(i int32) *int32 {
	return &i
}
