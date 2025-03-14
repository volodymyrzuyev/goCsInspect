package types

import "time"

var (
	TimeOutDuration = time.Second * 5
	RequestCooldown = time.Second * 2
)

const (
	CsGameID               = 730
	InspectRequestProtoID  = 9156
	InspectResponseProtoID = 9157
)
