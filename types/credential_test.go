package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

func TestValidate(t *testing.T) {
	validCreds := func() types.Credentials {
		return types.Credentials{
			Username:      "Test",
			Password:      "Test",
			TwoFactorCode: "Test",
			SharedSecret:  "SGVsbG9Xb3JsZA==",
		}
	}

	t.Run("Valid Credentials", func(t *testing.T) {
		creds := validCreds()

		err := creds.Validate()

		assert.Equal(t, nil, err, "Should not have an error")
	})

	t.Run("No Username", func(t *testing.T) {
		creds := validCreds()
		creds.Username = ""

		err := creds.Validate()

		assert.Equal(t, types.InvalidCredential, err, "Expecting an error")
	})

	t.Run("No Password", func(t *testing.T) {
		creds := validCreds()
		creds.Password = ""

		err := creds.Validate()

		assert.Equal(t, types.InvalidCredential, err, "Expecting an error")
	})

	t.Run("No 2FA", func(t *testing.T) {
		creds := validCreds()
		creds.TwoFactorCode = ""

		err := creds.Validate()

		assert.Equal(t, nil, err, "Should not have an error")
	})

	t.Run("No SharedSecret", func(t *testing.T) {
		creds := validCreds()
		creds.SharedSecret = ""

		err := creds.Validate()

		assert.Equal(t, nil, err, "Should not have an error")
	})

	t.Run("No SharedSecret or 2FA", func(t *testing.T) {
		creds := validCreds()
		creds.TwoFactorCode = ""
		creds.SharedSecret = ""

		err := creds.Validate()

		assert.Equal(t, types.InvalidCredential, err, "Expecting an error")
	})
}

func TestGet2FA(t *testing.T) {
	validCreds := func() types.Credentials {
		return types.Credentials{
			Username:      "Test",
			Password:      "Test",
			TwoFactorCode: "Test",
			SharedSecret:  "SGVsbG9Xb3JsZA==",
		}
	}

	t.Run("2FA provided", func(t *testing.T) {
		creds := validCreds()

		_, actual := creds.Get2FC()

		assert.Equal(t, nil, actual, "Should not get error")
	})

	t.Run("2FA not provided", func(t *testing.T) {
		creds := validCreds()
		creds.TwoFactorCode = ""

		_, actual := creds.Get2FC()

		assert.Equal(t, nil, actual, "Should not get error")
	})
	t.Run("NO 2FA or SharedSecret", func(t *testing.T) {
		creds := validCreds()
		creds.TwoFactorCode = ""
		creds.SharedSecret = ""

		_, actual := creds.Get2FC()

		assert.Equal(t, types.InvalidSharedSecret, actual, "Should get an error")
	})

}
