package globalTypes

type InspectParams struct {
	ParamM int64
	ParamA int64
	ParamD int64
	ParamS int64
}

func (i InspectParams) Validate() bool {
	if i.ParamA == 0 || i.ParamD == 0 {
		return false
	}

	if i.ParamM == 0 && i.ParamS == 0 {
		return false
	}

	return true
}
