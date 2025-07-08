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
	"github.com/volodymyrzuyev/goCsInspect/config"
)

type Detailer struct {
	allItems *csgo.Csgo
}

func NewDetailer() *Detailer {
	languageData, err := parser.Parse(config.EnglishFile)
	if err != nil {
		panic(err)
	}

	itemData, err := parser.Parse(config.GameItems)
	if err != nil {
		panic(err)
	}

	allItems, err := csgo.New(languageData, itemData)
	if err != nil {
		panic(err)
	}

	return &Detailer{
		allItems: allItems,
	}
}

func (d *Detailer) detailModificationsStickers(protos []*protobuf.CEconItemPreviewDataBlock_Sticker) ([]types.Modification, error) {
	var mods []types.Modification

	for _, proto := range protos {
		mod := getBaseMod(proto)

		stickerSubtype, ok := d.allItems.AllStickerItems[int(mod.StickerId)]
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

func (d *Detailer) detailModificationsChains(protos []*protobuf.CEconItemPreviewDataBlock_Sticker) ([]types.Modification, error) {
	var mods []types.Modification

	for _, proto := range protos {
		mod := getBaseMod(proto)

		chain, ok := d.allItems.Keychains[int(mod.StickerId)]
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
		Slot:      proto.GetSlot(),
		StickerId: proto.GetStickerId(),
		Wear:      proto.GetWear(),
		Rotation:  proto.GetRotation(),
		OffsetX:   proto.GetOffsetX(),
		OffsetY:   proto.GetOffsetY(),
		OffsetZ:   proto.GetOffsetZ(),
		Pattern:   proto.GetPattern(),
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

func (d *Detailer) InspectItems(proto *protobuf.CEconItemPreviewDataBlock) (*types.Item, error) {
	item := &types.Item{
		ItemID:         proto.GetItemid(),
		DefIndex:       proto.GetDefindex(),
		PaintIndex:     proto.GetPaintindex(),
		Rarity:         proto.GetRarity(),
		Quality:        proto.GetQuality(),
		PaintWear:      proto.GetPaintwear(),
		PaintSeed:      proto.GetPaintseed(),
		KillEaterValue: proto.GetKilleatervalue(),
		CustomName:     proto.GetCustomname(),
		Origin:         proto.GetOrigin(),
		FloatValue:     float64(math.Float32frombits(proto.GetPaintwear())),
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

	rarity, ok := d.allItems.Rarities[int(item.Rarity)]
	if !ok {
		slog.Error("Rarity not found", "item_id", item.ItemID, "rarity_index", item.Rarity)
		return nil, errors.ErrUnknownRarity
	}

	if defIndex, ok := d.allItems.DefIndecies[int(item.DefIndex)]; ok {
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
				item.ItemName = d.allItems.Stickerkits[int(item.Stickers[0].StickerId)].Name

			case sprayDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Spraykits[int(item.Stickers[0].StickerId)].Name

			case pathDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Patchkits[int(item.Stickers[0].StickerId)].Name

			case chainDefIndex:
				item.FullItemName = item.Keychains[0].Name
				item.ItemName = d.allItems.Keychains[int(item.Keychains[0].StickerId)].Name

			case musicDefIndex:
				music, ok := d.allItems.Musickits[int(proto.GetMusicindex())]
				if !ok {
					return nil, errors.ErrUnknownMusicIndex
				}
				item.ItemName = music.Name
				item.FullItemName = fmt.Sprintf("%s | %s", item.WeaponType, item.ItemName)

			default:
				slog.Debug("Unexpected def_index, details won't be populated",
					"def_index", item.DefIndex)

			}

		case *csgo.Character:
			item.WeaponType = "Agent"
			item.RarityName = rarity.CharacterRarityName
			item.ItemName = itemType.Name
			item.FullItemName = itemType.Name

		default:
			slog.Debug("Unexpected def_index, details won't be populated",
				"def_index", item.DefIndex)

		}
	} else {
		slog.Error("Unknown def_index", "def_index", item.DefIndex)
		return nil, errors.ErrUnknownDefIndex
	}
	if item.PaintIndex != defaultPaintKitIndex {
		paintKit, ok := d.allItems.Paintkits[int(item.PaintIndex)]
		if !ok {
			return nil, errors.ErrUnknownPaintIndex
		}

		item.MaxFloat = paintKit.MaxFloat.String()
		item.MinFloat = paintKit.MinFloat.String()
		item.ItemName = paintKit.Name
		item.FullItemName = fmt.Sprintf("%s | %s (%s)", item.WeaponType, item.ItemName, item.WearName)
	}

	quality, ok := d.allItems.Qualities[int(item.Quality)]
	if !ok {
		slog.Error("Quality not found",
			"item_id", item.ItemID, "quality_index", item.Quality)
		return nil, errors.ErrUnknownRarity
	}

	item.QualityName = quality.Name
	if item.Quality != defaultItemQuality {
		item.FullItemName = fmt.Sprintf("%s %s", item.QualityName, item.FullItemName)
	}

	return item, nil
}
