package globalTypes

type Skin struct {
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
	Stickers           []Stickers
	Inventory          int
	Origin             int
	Questid            int
	Dropreason         int
	Musicindex         int
	Entindex           int
	Petindex           int
	Keychains          []Keychains
}

type Keychains struct {
	Slot      int
	StickerId int
	Wear      float64
	Scale     float64
	Rotation  float64
	TintId    int
	OffsetX   float64
	OffsetY   float64
	OffsetZ   float64
	Pattern   int
}

type Stickers struct {
	Slot      int
	StickerId int
	Wear      float64
	Scale     float64
	Rotation  float64
	TintId    int
	OffsetX   float64
	OffsetY   float64
	OffsetZ   float64
	Pattern   int
}
