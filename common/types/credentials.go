package types

import (
	"errors"
	"github.com/Philipp15b/go-steam/v3"

	twoFA "github.com/bbqtd/go-steam-authenticator"
)

type Credentials struct {
	Username      string
	Password      string
	TwoFactorCode string
	SharedSecret  string
}

var InvalidCredential = errors.New("Invalid credentials, username and password and (2FC or SharedSecret) must be provided")
var InvalidSharedSecret = errors.New("Provided SharedSecret is invalid")

func (c Credentials) Validate() (err error) {
	if c.Username == "" || c.Password == "" || (c.TwoFactorCode == "" && c.SharedSecret == "") {
		return InvalidCredential
	}

	return
}

func (c Credentials) Get2FC() (code string, err error) {
	code = c.TwoFactorCode

	if code == "" {
		code, err = twoFA.GenerateAuthCode(c.SharedSecret, nil)
		if err != nil {
			err = InvalidSharedSecret
		}
	}

	return
}

func (c Credentials) GenerateLogOnDetails() (steam.LogOnDetails, error) {
	err := c.Validate()
	if err != nil {
		return steam.LogOnDetails{}, err
	}

	twoFA, err := c.Get2FC()
	if err != nil {
		return steam.LogOnDetails{}, err
	}

	logInInfo := steam.LogOnDetails{
		Username:      c.Username,
		Password:      c.Password,
		TwoFactorCode: twoFA,
	}

	return logInInfo, nil
}
