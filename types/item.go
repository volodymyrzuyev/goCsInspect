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
	S              uint64  //
	A              uint64  //
	D              uint64  //
	M              uint64  //
	FloatValue     float32 //
	MinFloat       float64 //
	MaxFloat       float64 //
	ItemName       string  //
	QualityName    string
	WeaponType     string //
	RarityName     string //
	WearName       string //
	FullItemName   string //
}

type Modification struct {
	Slot     uint32
	Rotation float32
	OffsetX  float32
	OffsetY  float32
	OffsetZ  float32
	CodeName string
	Material string
	Name     string
	Pattern  uint32
}
