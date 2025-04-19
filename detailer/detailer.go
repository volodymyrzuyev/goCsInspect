package detailer

import (
	"fmt"
	"math"
	"os"
	"strconv"

	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

type Detailer interface {
	GetDetails(*csProto.CEconItemPreviewDataBlock, types.InspectParameters) types.Item
}

func NewDetailer(log logger.Logger) (Detailer, error) {
	detailer := detailer{log: log}

	fmt.Println(os.Getwd())

	english, err := os.Open("./detailer/csgo_english.txt")
	if err != nil {
		return nil, err
	}

	items, err := os.Open("./detailer/items_game.txt")
	if err != nil {
		return nil, err
	}

	files := files{
		csgoEnglish: english,
		itemsGame:   items,
	}

	detailer.populateAssets(files)

	return &detailer, nil
}

type detailer struct {
	log logger.Logger

	englishTokens      map[string]string
	raritieDefenitions map[uint32]raritie
	itemDefenitions    map[string]any
	stickerDefenitions map[string]any
	paintKits          map[uint32]paintKit
	qualities          map[uint32]qualitie
	paintKitRarities   map[string]any
	keyChains          map[string]any
	itemSets           map[string]any
	prefabs            map[string]any
}

func (d *detailer) GetDetails(proto *csProto.CEconItemPreviewDataBlock, params types.InspectParameters) types.Item {
	stickers := d.parseSticker(proto.Stickers)
	keychains := d.parseKeychians(proto.Stickers)

	item := types.Item{
		ItemID:         uint64OrZero(proto.Itemid),
		DefIndex:       uint32OrZero(proto.Defindex),
		PaintIndex:     uint32OrZero(proto.Paintindex),
		Rarity:         uint32OrZero(proto.Rarity),
		Quality:        uint32OrZero(proto.Quality),
		PaintWear:      uint32OrZero(proto.Paintwear),
		PaintSeed:      uint32OrZero(proto.Paintseed),
		KillEaterValue: uint32OrZero(proto.Killeatervalue),
		CustomName:     stringOrNull(proto.Customname),
		Origin:         uint32OrZero(proto.Origin),
		Stickers:       stickers,
		Keychains:      keychains,
		S:              params.S,
		A:              params.A,
		D:              params.D,
		M:              params.M,
		FloatValue:     math.Float32frombits(proto.GetPaintwear()),
		MinFloat:       0.06,
		MaxFloat:       0.8,
		ItemName:       "",
	}

	item.WearName = getWearName(item.FloatValue)

	item.QualityName = d.qualities[item.Quality].englishName

	if defIndexInfo, ok := d.itemDefenitions[strconv.FormatUint(uint64(item.DefIndex), 10)].(map[string]interface{}); ok {
		if prefabName, ok := defIndexInfo["prefab"].(string); ok {
			if prefabName == "customplayertradable" {
				prefabName = "customplayer"
			}

			if itemPrefab, ok := d.prefabs[prefabName].(map[string]any); ok {
				tokenName := itemPrefab["item_name"].(string)
				if prefabName == "customplayer" {
					tokenName = itemPrefab["item_type_name"].(string)
				}

				item.WeaponType = d.englishTokens[tokenName[1:]]

				itemNameToken := defIndexInfo["item_name"]

				// non weapon item name can be gotten from the prefab
				if item.PaintWear == 0 {
					item.ItemName = d.englishTokens[itemNameToken.(string)[1:]]
					item.FullItemName = item.ItemName
					if prefabName != "customplayer" {
						item.FullItemName = fmt.Sprintf("%v | %v", item.WeaponType, item.ItemName)
					}
					if item.WearName != "" {
						item.FullItemName = fmt.Sprintf("%v | (%v)", item.FullItemName, item.WearName)
					}
				}
			}
		}
	}

	paintKitInfo, ok := d.paintKits[item.PaintIndex]

	if !ok || item.PaintWear == 0 {
		item.MaxFloat = 0
		item.MinFloat = 0
	} else {
		if minFloat, err := strconv.ParseFloat(paintKitInfo.wearRemapMin, 32); err == nil {
			item.MinFloat = minFloat
		}
		if maxFloat, err := strconv.ParseFloat(paintKitInfo.wearRemapMax, 32); err == nil {
			item.MaxFloat = maxFloat
		}

		// weapon names must be gotten from the paintKitInfo
		item.ItemName = paintKitInfo.itemName
		item.FullItemName = fmt.Sprintf("%v | %v (%v)", item.WeaponType, item.ItemName, item.WearName)

		// 4 is the default qualitie for items, if it's 4 it does not get displayed
		if item.Quality != 4 {
			item.FullItemName = fmt.Sprintf("%v %v", item.QualityName, item.FullItemName)
		}
	}

	d.populateRarityName(&item)

	return item
}

func (d *detailer) populateRarityName(item *types.Item) {
	rarityKey, ok := d.raritieDefenitions[item.Rarity]
	if !ok {
		d.log.Debug("Raritie %v not found", rarityKey)
	}
	if item.FloatValue != 0 && item.DefIndex != 1355 {
		if rarityName, ok := d.englishTokens[rarityKey.weapon]; ok {
			item.RarityName = rarityName
		} else {
			d.log.Debug("Raritie %v for weapon english name not found.", rarityKey)
		}
	} else {
		if item.DefIndex == 1355 || item.DefIndex == 1209 || item.DefIndex == 1348 {
			if rarityName, ok := d.englishTokens[rarityKey.regular]; ok {
				item.RarityName = rarityName
			} else {
				d.log.Debug("Raritie %v for default english name not found.", rarityKey)
			}
		} else {
			if rarityName, ok := d.englishTokens[rarityKey.character]; ok {
				item.RarityName = rarityName
			} else {
				d.log.Debug("Raritie %v for character english name not found.", rarityKey)
			}
		}
	}
}

func (d *detailer) parseSticker(list []*csProto.CEconItemPreviewDataBlock_Sticker) []types.Modification {
	return nil
}

func (d *detailer) parseKeychians(list []*csProto.CEconItemPreviewDataBlock_Sticker) []types.Modification {
	return nil
}

func getWearName(floatVal float32) string {
	if floatVal >= 0.45 {
		return "Battle-Scarred"
	}
	if floatVal >= 0.37 {
		return "Well-Worn"
	}
	if floatVal >= 0.15 {
		return "Field-Tested"
	}
	if floatVal >= 0.07 {
		return "Minimal Wear"
	}
	if floatVal > 0 {
		return "Factory New"
	}

	return ""
}

func getSteamFloat(val *uint32) float32 {
	if val == nil {
		return 0
	}
	float32Val := math.Float32frombits(*val)
	return float32(float32Val)
}

func uint32OrZero(val *uint32) uint32 {
	if val == nil {
		return 0
	}
	return *val
}

func uint64OrZero(val *uint64) uint64 {
	if val == nil {
		return 0
	}
	return *val
}

func stringOrNull(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
