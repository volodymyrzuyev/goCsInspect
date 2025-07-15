package inspect

import (
	"regexp"
	"strconv"

	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
)

type Params struct {
	M uint64
	A uint64
	D uint64
	S uint64
}

func ParseInspectLink(inspectLink string) (Params, error) {
	// hate regex will probably break
	re := regexp.MustCompile(
		`steam://rungame/730/.\d+/.*csgo_econ_action_preview.+([SM])(\d+)A(\d+)D(\d+)$`,
	)

	matches := re.FindStringSubmatch(inspectLink)

	var m, a, d, s uint64
	var err error

	if len(matches) != 5 {
		return Params{}, errors.ErrInvalidInspectLink
	}
	if matches[1] == "M" {
		m, _ = strconv.ParseUint(matches[2], 10, 64)
		s = uint64(0)
	} else {
		s, _ = strconv.ParseUint(matches[2], 10, 64)
		m = uint64(0)
	}

	a, err = strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return Params{}, errors.ErrInvalidInspectLink
	}

	d, err = strconv.ParseUint(matches[4], 10, 64)
	if err != nil {
		return Params{}, errors.ErrInvalidInspectLink
	}

	params := Params{
		M: m,
		A: a,
		D: d,
		S: s,
	}

	return params, nil
}

func (i Params) Validate() (err error) {
	err = nil
	if i.D == 0 || i.A == 0 || (i.M == 0 && i.S == 0) {
		err = errors.ErrInvalidParameters
	}
	return
}

func (i Params) GenerateGcRequestProto() (*csProto.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest, error) {
	if err := i.Validate(); err != nil {
		return nil, err
	}

	requestProto := csProto.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest{
		ParamS: &i.S,
		ParamA: &i.A,
		ParamD: &i.D,
		ParamM: &i.M,
	}

	return &requestProto, nil
}
