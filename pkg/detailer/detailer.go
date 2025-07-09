package detailer

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/go-csgo-item-parser/csgo"
	"github.com/volodymyrzuyev/go-csgo-item-parser/parser"

	"github.com/volodymyrzuyev/goCsInspect/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/common/types"
)

type Detailer interface {
	DetailProto(proto *protobuf.CEconItemPreviewDataBlock) (*types.Item, error)
}

type detailer struct {
	allItems *csgo.Csgo
}

func NewDetailer(langugeFile, gameItems string) Detailer {
	languageData, err := parser.Parse(langugeFile)
	if err != nil {
		panic(err)
	}

	itemData, err := parser.Parse(gameItems)
	if err != nil {
		panic(err)
	}

	allItems, err := csgo.New(languageData, itemData)
	if err != nil {
		panic(err)
	}

	return &detailer{
		allItems: allItems,
	}
}

func (d *detailer) detailModificationsStickers(protos []*protobuf.CEconItemPreviewDataBlock_Sticker) ([]types.Modification, error) {
	var mods []types.Modification

	for _, proto := range protos {
		mod := getBaseMod(proto)

		stickerSubtype, ok := d.allItems.AllStickerItems[int(mod.Proto.GetStickerId())]
		if !ok {
			return nil, errors.ErrUnknownStickerModifier
		}
		switch s := stickerSubtype.(type) {
		case *csgo.Stickerkit:
			mod.CodeName = s.Id
			mod.Material = s.Variant
			mod.Name = fmt.Sprintf("Sticker | %s", s.Name)

		case *csgo.Spraykit:
			mod.CodeName = s.Id
			mod.Name = fmt.Sprintf("Sealed Graffiti | %s", s.Name)

		case *csgo.Patchkit:
			mod.CodeName = s.Id
			mod.Name = fmt.Sprintf("Patch | %s", s.Name)
		default:
			return nil, errors.ErrUnknownStickerModifier
		}

		mods = append(mods, mod)
	}

	return mods, nil
}

func (d *detailer) detailModificationsChains(protos []*protobuf.CEconItemPreviewDataBlock_Sticker) ([]types.Modification, error) {
	var mods []types.Modification

	for _, proto := range protos {
		mod := getBaseMod(proto)

		chain, ok := d.allItems.Keychains[int(mod.Proto.GetStickerId())]
		if !ok {
			return nil, errors.ErrUnknownStickerModifier
		}
		mod.CodeName = chain.Id
		mod.Name = fmt.Sprintf("Charm | %s", chain.Name)

		mods = append(mods, mod)
	}

	return mods, nil
}

func getBaseMod(proto *protobuf.CEconItemPreviewDataBlock_Sticker) types.Modification {
	return types.Modification{
		Proto: proto,
	}
}

func getWearName(float float64) string {

	if float >= 0.45 {
		return "Battle-Scarred"
	}
	if float >= 0.37 {
		return "Well-Worn"
	}
	if float >= 0.15 {
		return "Field-Tested"
	}
	if float >= 0.07 {
		return "Minimal Wear"
	}
	if float > 0 {
		return "Factory New"
	}

	return ""
}

const (
	minDefaultFloat = "0.06"
	maxDefaultFloat = "0.8"

	stickerDefIndex = 1209
	pathDefIndex    = 4609
	sprayDefIndex   = 1348
	chainDefIndex   = 1355
	musicDefIndex   = 1314

	defaultItemQuality   = 4
	defaultPaintKitIndex = 0
)

func (d *detailer) DetailProto(proto *protobuf.CEconItemPreviewDataBlock) (*types.Item, error) {
	item := &types.Item{
		Proto:      proto,
		FloatValue: float64(math.Float32frombits(proto.GetPaintwear())),
	}

	var err error
	item.Stickers, err = d.detailModificationsStickers(proto.Stickers)
	if err != nil {
		slog.Error("Unknown sticker", "proto", proto)
		return nil, err
	}

	item.Keychains, err = d.detailModificationsChains(proto.Keychains)
	if err != nil {
		slog.Error("Unknown sticker", "proto", proto)
		return nil, err
	}

	item.WearName = getWearName(item.FloatValue)

	rarity, ok := d.allItems.Rarities[int(item.Proto.GetRarity())]
	if !ok {
		slog.Error("Rarity not found", "item_id", item.Proto.GetItemid(), "rarity_index", item.Proto.GetRarity())
		return nil, errors.ErrUnknownRarity
	}

	if defIndex, ok := d.allItems.DefIndecies[int(item.Proto.GetDefindex())]; ok {
		switch itemType := defIndex.(type) {
		case *csgo.Weapon:
			item.MaxFloat = maxDefaultFloat
			item.MinFloat = minDefaultFloat
			item.WeaponType = itemType.Name
			item.RarityName = rarity.WeaponRarityName
			item.FullItemName = fmt.Sprintf("%v | (%v)", item.WeaponType, item.WearName)

		case *csgo.Gloves:
			item.MaxFloat = maxDefaultFloat
			item.MinFloat = minDefaultFloat
			item.WeaponType = itemType.Name
			item.RarityName = rarity.WeaponRarityName
			item.FullItemName = fmt.Sprintf("%v | (%v)", item.WeaponType, item.WearName)

		case *csgo.Tool:
			item.WeaponType = itemType.Name
			item.RarityName = rarity.GeneralRarityName

			switch itemType.Index {
			case stickerDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Stickerkits[int(item.Stickers[0].Proto.GetStickerId())].Name

			case sprayDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Spraykits[int(item.Stickers[0].Proto.GetStickerId())].Name

			case pathDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Patchkits[int(item.Stickers[0].Proto.GetStickerId())].Name

			case chainDefIndex:
				item.FullItemName = item.Keychains[0].Name
				item.ItemName = d.allItems.Keychains[int(item.Keychains[0].Proto.GetStickerId())].Name

			case musicDefIndex:
				music, ok := d.allItems.Musickits[int(proto.GetMusicindex())]
				if !ok {
					return nil, errors.ErrUnknownMusicIndex
				}
				item.ItemName = music.Name
				item.FullItemName = fmt.Sprintf("%s | %s", item.WeaponType, item.ItemName)

			default:
				slog.Debug("Unexpected def_index, details won't be populated",
					"def_index", item.Proto.GetDefindex())

			}

		case *csgo.Character:
			item.WeaponType = "Agent"
			item.RarityName = rarity.CharacterRarityName
			item.ItemName = itemType.Name
			item.FullItemName = itemType.Name

		default:
			slog.Debug("Unexpected def_index, details won't be populated",
				"def_index", item.Proto.GetDefindex())

		}
	} else {
		slog.Error("Unknown def_index", "def_index", item.Proto.GetDefindex())
		return nil, errors.ErrUnknownDefIndex
	}
	if item.Proto.GetPaintindex() != defaultPaintKitIndex {
		paintKit, ok := d.allItems.Paintkits[int(item.Proto.GetPaintindex())]
		if !ok {
			return nil, errors.ErrUnknownPaintIndex
		}

		item.MaxFloat = paintKit.MaxFloat.String()
		item.MinFloat = paintKit.MinFloat.String()
		item.ItemName = paintKit.Name
		item.FullItemName = fmt.Sprintf("%s | %s (%s)", item.WeaponType, item.ItemName, item.WearName)
	}

	quality, ok := d.allItems.Qualities[int(item.Proto.GetQuality())]
	if !ok {
		slog.Error("Quality not found",
			"item_id", item.Proto.GetItemid(), "quality_index", item.Proto.GetQuality())
		return nil, errors.ErrUnknownRarity
	}

	item.QualityName = quality.Name
	if item.Proto.GetQuality() != defaultItemQuality {
		item.FullItemName = fmt.Sprintf("%s %s", item.QualityName, item.FullItemName)
	}

	return item, nil
}
