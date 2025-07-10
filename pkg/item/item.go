package item

import "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"

type Item struct {
	Proto        *protobuf.CEconItemPreviewDataBlock
	Stickers     []Modification
	Keychains    []Modification
	FloatValue   float64 //
	MinFloat     string  //
	MaxFloat     string  //
	ItemName     string  //
	QualityName  string  //
	WeaponType   string  //
	RarityName   string  //
	WearName     string  //
	FullItemName string  //
}

type Modification struct {
	Proto *protobuf.CEconItemPreviewDataBlock_Sticker

	CodeName string
	Material string
	Name     string
}
