package detailer_test

import (
	"os"
	"testing"

	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/detailer"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

type testCase struct {
	inputProto   csProto.CEconItemPreviewDataBlock
	inputInspect types.InspectParameters
	output       types.Item
}

func TestDetailer(t *testing.T) {
	log := logger.NewLogger(os.Stdout)

	detailer, err := detailer.NewDetailer(log)
	if err != nil {
		t.FailNow()
	}

	tests := make(map[string]testCase)
	// https://steamcommunity.com/market/listings/730/Special%20Agent%20Ava%20%7C%20FBI
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M663826331106610031A43192124620D5108495161639654660
	tests["Agent"] = testCase{
		inputProto: csProto.CEconItemPreviewDataBlock{
			Itemid:    getUint64pointer(43192124620),
			Defindex:  getUint32pointer(5308),
			Quality:   getUint32pointer(4),
			Rarity:    getUint32pointer(6),
			Inventory: getUint32pointer(16),
			Origin:    getUint32pointer(23),
		},
		inputInspect: types.InspectParameters{A: 43192124620, D: 5108495161639654660, M: 663826331106610031},
		output: types.Item{
			ItemID:       43192124620,
			DefIndex:     5308,
			Rarity:       6,
			Quality:      4,
			Origin:       23,
			A:            43192124620,
			D:            5108495161639654660,
			M:            663826331106610031,
			QualityName:  "Unique",
			ItemName:     "Special Agent Ava | FBI",
			WeaponType:   "Agent",
			RarityName:   "Master",
			WearName:     "",
			FullItemName: "Special Agent Ava | FBI",
		},
	}

	//itemid:26267148  defindex:11  paintindex:6  rarity:2  quality:4  paintwear:1040573991  paintseed:454  inventory:10  origin:0
	// https://steamcommunity.com/market/listings/730/G3SG1%20%7C%20Arctic%20Camo%20%28Minimal%20Wear%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M2128243055175485935A26267148D435200037615186485
	tests["Skin No Stickers Not StatTrak"] = testCase{
		inputProto: csProto.CEconItemPreviewDataBlock{
			Itemid:     getUint64pointer(26267148),
			Defindex:   getUint32pointer(11),
			Paintindex: getUint32pointer(6),
			Rarity:     getUint32pointer(2),
			Quality:    getUint32pointer(4),
			Paintwear:  getUint32pointer(1040573991),
			Paintseed:  getUint32pointer(454),
			Inventory:  getUint32pointer(10),
			Origin:     getUint32pointer(0),
		},
		inputInspect: types.InspectParameters{A: 26267148, D: 435200037615186485, M: 2128243055175485935},
		output: types.Item{
			ItemID:       26267148,
			DefIndex:     11,
			PaintIndex:   6,
			Rarity:       2,
			Quality:      4,
			PaintWear:    1040573991,
			PaintSeed:    454,
			Origin:       0,
			A:            26267148,
			D:            435200037615186485,
			M:            2128243055175485935,
			FloatValue:   0.13076077401638,
			MinFloat:     0.06,
			MaxFloat:     0.8,
			QualityName:  "Unique",
			ItemName:     "Arctic Camo",
			WeaponType:   "G3SG1",
			RarityName:   "Industrial Grade",
			WearName:     "Minimal Wear",
			FullItemName: "G3SG1 | Arctic Camo (Minimal Wear)",
		},
	}

	// itemid:8839911627  defindex:39  paintindex:98  rarity:3  quality:9  paintwear:1060635000  paintseed:369  killeaterscoretype:0  killeatervalue:218  inventory:5  origin:8
	// https://steamcommunity.com/market/listings/730/StatTrak%E2%84%A2%20SG%20553%20%7C%20Ultraviolet%20%28Battle-Scarred%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M640180078099397070A8839911627D9837170024513331701
	tests["Skin No Stickers StatTrak"] = testCase{
		inputProto: csProto.CEconItemPreviewDataBlock{
			Itemid:             getUint64pointer(8839911627),
			Defindex:           getUint32pointer(39),
			Paintindex:         getUint32pointer(98),
			Rarity:             getUint32pointer(3),
			Quality:            getUint32pointer(9),
			Paintwear:          getUint32pointer(1060635000),
			Paintseed:          getUint32pointer(369),
			Killeaterscoretype: getUint32pointer(0),
			Killeatervalue:     getUint32pointer(218),
			Inventory:          getUint32pointer(5),
			Origin:             getUint32pointer(8),
		},
		inputInspect: types.InspectParameters{A: 8839911627, D: 9837170024513331701, M: 640180078099397070},
		output: types.Item{
			ItemID:         8839911627,
			DefIndex:       39,
			PaintIndex:     98,
			Rarity:         3,
			Quality:        9,
			PaintWear:      1060635000,
			PaintSeed:      369,
			Origin:         8,
			KillEaterValue: 218,
			A:              8839911627,
			D:              9837170024513331701,
			M:              640180078099397070,
			FloatValue:     0.718772411346436,
			MinFloat:       0.06,
			MaxFloat:       0.8,
			QualityName:    "StatTrak™",
			ItemName:       "Ultraviolet",
			WeaponType:     "SG 553",
			RarityName:     "Mil-Spec Grade",
			WearName:       "Battle-Scarred",
			FullItemName:   "StatTrak™ SG 553 | Ultraviolet (Battle-Scarred)",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp := detailer.GetDetails(&test.inputProto, test.inputInspect)
			assert.Equal(t, test.output, resp, "Should be same")
		})
	}
}

func getUint64pointer(n int) *uint64 {
	uintp := uint64(n)
	return &uintp
}

func getUint32pointer(n int) *uint32 {
	uintp := uint32(n)
	return &uintp
}
