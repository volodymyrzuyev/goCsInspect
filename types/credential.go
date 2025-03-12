package types

import (
	"errors"

	twoFA "github.com/bbqtd/go-steam-authenticator"
)

type Credentials struct {
	Username      string
	Password      string
	TwoFactorCode string
	SharedSecret  string
}

var InvalidCredential = errors.New("Invalid credentials, username and password and (2FC or SharedSecret) must be provided")
var InvalidSharedSecret = errors.New("Provided shared secret is invalid")

var generateCode = twoFA.GenerateAuthCode

func (c Credentials) Validate() (err error) {
	if c.Username == "" || c.Password == "" || (c.TwoFactorCode == "" && c.SharedSecret == "") {
		return InvalidCredential
	}

	return
}

func (c Credentials) Get2FC() (code string, err error) {
	code = c.TwoFactorCode

	if code == "" {
		code, err = generateCode(c.SharedSecret, nil)
		if err != nil {
			err = InvalidSharedSecret
		}
	}

	return
}
