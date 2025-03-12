package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

func TestInspectParamsValidator(t *testing.T) {
	validParams := func() types.InspectParameters {
		return types.InspectParameters{
			M: 10000,
			A: 10000,
			D: 10000,
			S: 10000,
		}
	}

	t.Run("Valid Params", func(t *testing.T) {
		params := validParams()

		err := params.Validate()

		assert.Equal(t, nil, err, "All params provided, should be no err")
	})

	t.Run("No M", func(t *testing.T) {
		params := validParams()
		params.M = 0

		err := params.Validate()

		assert.Equal(t, nil, err, "Either S or M are needed, is one is missing, there should be no err")
	})

	t.Run("No S", func(t *testing.T) {
		params := validParams()
		params.S = 0

		err := params.Validate()

		assert.Equal(t, nil, err, "Either S or M are needed, is one is missing, there should be no err")
	})

	t.Run("No S or S", func(t *testing.T) {
		params := validParams()
		params.S = 0
		params.M = 0

		err := params.Validate()

		assert.Equal(t, types.InvalidParameters, err, "S or M are required, should be an err")
	})

	t.Run("Invalid D", func(t *testing.T) {
		params := validParams()
		params.D = 0

		err := params.Validate()

		assert.Equal(t, types.InvalidParameters, err, "D is a required filed, should be an err")
	})

	t.Run("Invalid A", func(t *testing.T) {
		params := validParams()
		params.A = 0

		err := params.Validate()

		assert.Equal(t, types.InvalidParameters, err, "A is a required filed, should be an err")
	})
}

func TestParsingInspectLink(t *testing.T) {
	t.Run("Valid link with M", func(t *testing.T) {
		inspectLink := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M5634632664282712758A35234550478D9657633252597596896"
		params, err := types.ParseInspectLink(inspectLink)

		assert.Equal(t, uint64(5634632664282712758), params.M, "M params should be equal")
		assert.Equal(t, uint64(35234550478), params.A, "A params should be equal")
		assert.Equal(t, uint64(9657633252597596896), params.D, "D params should be equal")
		assert.Equal(t, uint64(0), params.S, "S params should be equal")
		assert.Equal(t, nil, err, "Should be no error")
	})

	t.Run("Valid link with S", func(t *testing.T) {
		inspectLink := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20S5634632664282712758A35234550478D9657633252597596896"
		params, err := types.ParseInspectLink(inspectLink)

		assert.Equal(t, uint64(0), params.M, "M params should be equal")
		assert.Equal(t, uint64(35234550478), params.A, "A params should be equal")
		assert.Equal(t, uint64(9657633252597596896), params.D, "D params should be equal")
		assert.Equal(t, uint64(5634632664282712758), params.S, "S params should be equal")
		assert.Equal(t, nil, err, "Should be no error")
	})

	t.Run("Messed up M link", func(t *testing.T) {
		inspectLink := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M5634632sdasd664282712758A35234550478D9657633252597596896"
		params, err := types.ParseInspectLink(inspectLink)

		assert.Equal(t, uint64(0), params.M, "M params should be equal")
		assert.Equal(t, uint64(0), params.A, "A params should be equal")
		assert.Equal(t, uint64(0), params.D, "D params should be equal")
		assert.Equal(t, uint64(0), params.S, "S params should be equal")
		assert.Equal(t, types.InvalidInspectLink, err, "Should be no error")
	})

	t.Run("No M link", func(t *testing.T) {
		inspectLink := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%A35234550478D9657633252597596896"
		params, err := types.ParseInspectLink(inspectLink)

		assert.Equal(t, uint64(0), params.M, "M params should be equal")
		assert.Equal(t, uint64(0), params.A, "A params should be equal")
		assert.Equal(t, uint64(0), params.D, "D params should be equal")
		assert.Equal(t, uint64(0), params.S, "S params should be equal")
		assert.Equal(t, types.InvalidInspectLink, err, "Should be no error")
	})

	t.Run("M and S together link", func(t *testing.T) {
		inspectLink := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M56346S32664282712758A35234550478D9657633252597596896"
		params, err := types.ParseInspectLink(inspectLink)

		assert.Equal(t, uint64(0), params.M, "M params should be equal")
		assert.Equal(t, uint64(0), params.A, "A params should be equal")
		assert.Equal(t, uint64(0), params.D, "D params should be equal")
		assert.Equal(t, uint64(0), params.S, "S params should be equal")
		assert.Equal(t, types.InvalidInspectLink, err, "Should be an error")
	})
}
