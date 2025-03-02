package globalTypes

import "time"

type Item struct {
	AccountID          int
	ItemID             int
	DefIndex           int
	PaintIndex         int
	Rarity             int
	Quality            int
	Paintwear          int
	Paintseed          int
	Killeaterscoretype int
	Killeatervalue     int
	Customname         string
	Stickers           []Modifier
	Inventory          int
	Origin             int
	Questid            int
	Dropreason         int
	Musicindex         int
	Entindex           int
	Petindex           int
	Keychains          []Modifier
	ParamD             int
	ParamM             int
	ParamS             int
	FloatValue         float64
	MaxFloat           float64
	MinFloat           float64
	WeaponType         string
	ItemName           string
	RarityName         string
	QualityName        string
	OriginName         string
	WearName           string
	MarketHashName     string
	LastModified       time.Time
}

type Modifier struct {
	Slot           int
	ModifierId     int
	Wear           float64
	Scale          float64
	Rotation       float64
	TintId         int
	OffsetX        float64
	OffsetY        float64
	OffsetZ        float64
	Pattern        int
	MarketHashName string
}
