package detailer

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/go-csgo-item-parser/csgo"
	"github.com/volodymyrzuyev/go-csgo-item-parser/parser"

	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/item"
)

type Detailer interface {
	// Populates extra details for items
	DetailProto(proto *protobuf.CEconItemPreviewDataBlock) (*item.Item, error)
	// Updates data used to detail a protobuf
	UpdateItems(newItems *csgo.Csgo)
}

type detailer struct {
	allItems *csgo.Csgo
	l        *slog.Logger
}

func NewDetailerGameFiles(langugeFile, gameItems string) (Detailer, error) {
	ln := slog.Default().WithGroup("Detailer")

	languageData, err := parser.Parse(langugeFile)
	if err != nil {
		ln.Error("could not parser language file", "error", err)
		return nil, err
	}

	itemData, err := parser.Parse(gameItems)
	if err != nil {
		ln.Error("could not parser item file", "error", err)
		return nil, err
	}

	allItems, err := csgo.New(languageData, itemData)
	if err != nil {
		ln.Error("could not parser cs files", "error", err)
		return nil, err
	}

	return NewDetailerWithCSItems(allItems)
}

func NewDetailerWithCSItems(items *csgo.Csgo) (Detailer, error) {
	return &detailer{allItems: items, l: slog.Default().WithGroup("Detailer")}, nil
}

func (d *detailer) detailModificationsStickers(item *item.Item) error {
	for _, sticker := range item.Stickers {

		stickerSubtype, ok := d.allItems.AllStickerItems[int(sticker.StickerId)]
		if !ok {
			d.l.Error("unknown sticker id", "id", sticker.StickerId)
			return errors.ErrUnknownStickerModifier
		}

		switch s := stickerSubtype.(type) {
		case *csgo.Stickerkit:
			sticker.CodeName = s.Id
			sticker.Material = s.Variant
			sticker.Name = fmt.Sprintf("Sticker | %s", s.Name)

		case *csgo.Spraykit:
			sticker.CodeName = s.Id
			sticker.Name = fmt.Sprintf("Sealed Graffiti | %s", s.Name)

		case *csgo.Patchkit:
			sticker.CodeName = s.Id
			sticker.Name = fmt.Sprintf("Patch | %s", s.Name)
		default:
			d.l.Error("unknown sticker id", "id", sticker.StickerId)
			return errors.ErrUnknownStickerModifier
		}
	}

	return nil
}

func (d *detailer) detailModificationsChains(item *item.Item) error {
	for _, chainMod := range item.Keychains {

		chain, ok := d.allItems.Keychains[int(chainMod.StickerId)]
		if !ok {
			d.l.Debug("unknown keychain id", "id", chainMod.StickerId)
			return errors.ErrUnknownStickerModifier
		}

		chainMod.CodeName = chain.Id
		chainMod.Name = fmt.Sprintf("Charm | %s", chain.Name)
	}

	return nil
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
	statTrackQuality     = 9
)

func (d *detailer) DetailProto(proto *protobuf.CEconItemPreviewDataBlock) (*item.Item, error) {
	item := &item.Item{}
	item.PopulateProto(proto)
	item.FloatValue, _ = strconv.ParseFloat(
		fmt.Sprintf("%.15f", float64(math.Float32frombits(proto.GetPaintwear()))),
		32,
	)

	if err := d.detailModificationsStickers(item); err != nil {
		return nil, err
	}

	if err := d.detailModificationsChains(item); err != nil {
		return nil, err
	}

	item.WearName = getWearName(item.FloatValue)

	rarity, ok := d.allItems.Rarities[int(proto.GetRarity())]
	if !ok {
		d.l.Error("rarity not found", "rarity_index", proto.GetRarity())
		return nil, errors.ErrUnknownRarity
	}

	if defIndex, ok := d.allItems.DefIndecies[int(proto.GetDefindex())]; ok {
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
				item.ItemName = d.allItems.Stickerkits[int(proto.Stickers[0].GetStickerId())].Name

			case sprayDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Spraykits[int(proto.Stickers[0].GetStickerId())].Name

			case pathDefIndex:
				item.FullItemName = item.Stickers[0].Name
				item.ItemName = d.allItems.Patchkits[int(proto.Stickers[0].GetStickerId())].Name

			case chainDefIndex:
				item.FullItemName = item.Keychains[0].Name
				item.ItemName = d.allItems.Keychains[int(proto.Keychains[0].GetStickerId())].Name

			case musicDefIndex:
				music, ok := d.allItems.Musickits[int(proto.GetMusicindex())]
				if !ok {
					d.l.Error("unknown music kit", "music_index", proto.GetMusicindex())
					return nil, errors.ErrUnknownMusicIndex
				}
				item.ItemName = music.Name
				item.FullItemName = fmt.Sprintf("%s | %s", item.WeaponType, item.ItemName)

			default:
				d.l.Warn("unexpected def_index, details won't be populated", "def_index", proto.GetDefindex())
			}

		case *csgo.Collectible:
			item.WeaponType = "Pin"
			item.RarityName = rarity.GeneralRarityName
			item.ItemName = itemType.Name
			item.FullItemName = itemType.Name

		case *csgo.Character:
			item.WeaponType = "Agent"
			item.RarityName = rarity.CharacterRarityName
			item.ItemName = itemType.Name
			item.FullItemName = itemType.Name

		default:
			d.l.Warn("unexpected def_index, details won't be populated", "def_index", proto.GetDefindex())

		}
	} else {
		slog.Error("unknown def_index", "def_index", proto.GetDefindex())
		return nil, errors.ErrUnknownDefIndex
	}
	if proto.GetPaintindex() != defaultPaintKitIndex {
		paintKit, ok := d.allItems.Paintkits[int(proto.GetPaintindex())]
		if !ok {
			slog.Error("unknown paintIndex", "paintIndex", proto.GetPaintindex())
			return nil, errors.ErrUnknownPaintIndex
		}

		item.MaxFloat = paintKit.MaxFloat.String()
		item.MinFloat = paintKit.MinFloat.String()
		item.ItemName = paintKit.Name
		item.FullItemName = fmt.Sprintf(
			"%s | %s (%s)",
			item.WeaponType,
			item.ItemName,
			item.WearName,
		)
	} else {
		// since the item has a default paintIndex, and it's a knife, this means
		// it's a vanila knife, they don't show their WearName in the FullItemName
		if _, ok := d.allItems.Knives[int(proto.GetDefindex())]; ok {
			item.FullItemName = item.WeaponType
		}
	}

	quality, ok := d.allItems.Qualities[int(proto.GetQuality())]
	if !ok {
		slog.Error("unknown quality", "quality_index", proto.GetQuality())
		return nil, errors.ErrUnknownRarity
	}

	item.QualityName = quality.Name
	if proto.GetQuality() != defaultItemQuality {
		if _, ok := d.allItems.Knives[int(proto.GetDefindex())]; proto.Killeaterscoretype != nil &&
			ok {
			item.FullItemName = fmt.Sprintf(
				"%s %s %s",
				item.QualityName,
				d.allItems.Qualities[statTrackQuality].Name,
				item.FullItemName,
			)
		} else {
			item.FullItemName = fmt.Sprintf("%s %s", item.QualityName, item.FullItemName)
		}
	}

	return item, nil
}

func (d *detailer) UpdateItems(newItems *csgo.Csgo) {
	d.allItems = newItems
}
