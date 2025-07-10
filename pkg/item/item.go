package item

import (
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
)

type Item struct {
	Accountid          uint32
	Itemid             uint64
	Defindex           uint32
	Paintindex         uint32
	Rarity             uint32
	Quality            uint32
	Paintwear          uint32
	Paintseed          uint32
	Killeaterscoretype uint32
	Killeatervalue     uint32
	Customname         string
	Inventory          uint32
	Origin             uint32
	Questid            uint32
	Dropreason         uint32
	Musicindex         uint32
	Entindex           int32
	Petindex           uint32

	Stickers     []*Modification
	Keychains    []*Modification
	FloatValue   float64
	MinFloat     string
	MaxFloat     string
	ItemName     string
	QualityName  string
	WeaponType   string
	RarityName   string
	WearName     string
	FullItemName string
}

func (i *Item) PopulateProto(proto *protobuf.CEconItemPreviewDataBlock) {
	i.Accountid = proto.GetAccountid()
	i.Itemid = proto.GetItemid()
	i.Defindex = proto.GetDefindex()
	i.Paintindex = proto.GetPaintindex()
	i.Rarity = proto.GetRarity()
	i.Quality = proto.GetQuality()
	i.Paintwear = proto.GetPaintwear()
	i.Paintseed = proto.GetPaintseed()
	i.Killeaterscoretype = proto.GetKilleaterscoretype()
	i.Killeatervalue = proto.GetKilleatervalue()
	i.Customname = proto.GetCustomname()
	i.Inventory = proto.GetInventory()
	i.Origin = proto.GetOrigin()
	i.Questid = proto.GetQuestid()
	i.Dropreason = proto.GetDropreason()
	i.Musicindex = proto.GetMusicindex()
	i.Entindex = proto.GetEntindex()
	i.Petindex = proto.GetPetindex()
	if len(i.Keychains) == 0 {
		i.Keychains = ParseProtoMods(proto.GetKeychains())
	} else {
		for i, chain := range i.Keychains {
			chain.PopulateProto(proto.GetKeychains()[i])
		}
	}
	if len(i.Stickers) == 0 {
		i.Stickers = ParseProtoMods(proto.GetStickers())
	} else {
		for i, stickers := range i.Stickers {
			stickers.PopulateProto(proto.GetStickers()[i])
		}
	}
}

type Modification struct {
	Slot          uint32
	StickerId     uint32
	Wear          float32
	Scale         float32
	Rotation      float32
	TintId        uint32
	OffsetX       float32
	OffsetY       float32
	OffsetZ       float32
	Pattern       uint32
	HighlightReel uint32

	CodeName string
	Material string
	Name     string
}

func ParseProtoMods(protos []*protobuf.CEconItemPreviewDataBlock_Sticker) []*Modification {
	mods := make([]*Modification, len(protos))

	for i, proto := range protos {
		mods[i] = &Modification{
			Slot:          proto.GetSlot(),
			StickerId:     proto.GetStickerId(),
			Wear:          proto.GetWear(),
			Scale:         proto.GetScale(),
			Rotation:      proto.GetRotation(),
			TintId:        proto.GetTintId(),
			OffsetX:       proto.GetOffsetX(),
			OffsetY:       proto.GetOffsetY(),
			OffsetZ:       proto.GetOffsetZ(),
			Pattern:       proto.GetPattern(),
			HighlightReel: proto.GetHighlightReel(),
		}
	}

	return mods
}

func (m *Modification) PopulateProto(proto *protobuf.CEconItemPreviewDataBlock_Sticker) {
	m.Slot = proto.GetSlot()
	m.StickerId = proto.GetStickerId()
	m.Wear = proto.GetWear()
	m.Scale = proto.GetScale()
	m.Rotation = proto.GetRotation()
	m.TintId = proto.GetTintId()
	m.OffsetX = proto.GetOffsetX()
	m.OffsetY = proto.GetOffsetY()
	m.OffsetZ = proto.GetOffsetZ()
	m.Pattern = proto.GetPattern()
	m.HighlightReel = proto.GetHighlightReel()
}
