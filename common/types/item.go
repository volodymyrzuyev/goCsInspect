package types

type Item struct {
	ItemID         uint64 //
	DefIndex       uint32 //
	PaintIndex     uint32 //
	Rarity         uint32 //
	Quality        uint32 //
	PaintWear      uint32 //
	PaintSeed      uint32 //
	KillEaterValue uint32 //
	CustomName     string //
	Origin         uint32 //
	Stickers       []Modification
	Keychains      []Modification
	FloatValue     float64 //
	MinFloat       string  //
	MaxFloat       string  //
	ItemName       string  //
	QualityName    string  //
	WeaponType     string  //
	RarityName     string  //
	WearName       string  //
	FullItemName   string  //
}

type Modification struct {
	Slot      uint32
	StickerId uint32
	Wear      float32
	Rotation  float32
	OffsetX   float32
	OffsetY   float32
	OffsetZ   float32
	Pattern   uint32

	CodeName string
	Material string
	Name     string
}
