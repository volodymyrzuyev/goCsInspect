package types

import (
	"errors"
	"regexp"
	"strconv"
)

type InspectParameters struct {
	M uint64
	A uint64
	D uint64
	S uint64
}

var (
	InvalidParameters  = errors.New("Parameters A and D and (M or S) must be provided")
	InvalidInspectLink = errors.New("Was not able to parse inspectLink")
)

func (i InspectParameters) Validate() (err error) {
	err = nil
	if i.D == 0 || i.A == 0 || (i.M == 0 && i.S == 0) {
		err = InvalidParameters
	}
	return
}

func ParseInspectLink(inspectLink string) (InspectParameters, error) {
	// hate regex will probably break
	re := regexp.MustCompile(`steam://rungame/730/.\d+/.*csgo_econ_action_preview%20([SM])(\d+)A(\d+)D(\d+)$`)

	matches := re.FindStringSubmatch(inspectLink)

	var m, a, d, s uint64
	var err error

	if len(matches) != 5 {
		return InspectParameters{}, InvalidInspectLink
	}
	if matches[1] == "M" {
		m, err = strconv.ParseUint(matches[2], 10, 64)
		s = uint64(0)
	} else {
		s, err = strconv.ParseUint(matches[2], 10, 64)
		m = uint64(0)
	}

	a, err = strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return InspectParameters{}, InvalidInspectLink
	}

	d, err = strconv.ParseUint(matches[4], 10, 64)
	if err != nil {
		return InspectParameters{}, InvalidInspectLink
	}

	params := InspectParameters{
		M: m,
		A: a,
		D: d,
		S: s,
	}

	return params, nil
}
