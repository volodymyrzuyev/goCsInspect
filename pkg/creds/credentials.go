package creds

import (
	"github.com/Philipp15b/go-steam/v3"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"

	twoFA "github.com/bbqtd/go-steam-authenticator"
)

type Account struct {
	Username      string
	Password      string
	TwoFactorCode string
	SharedSecret  string
}

func (c Account) Validate() (err error) {
	if c.Username == "" || c.Password == "" || (c.TwoFactorCode == "" && c.SharedSecret == "") {
		return errors.ErrInsufficientCredentials
	}

	return
}

func (c Account) Get2FC() (code string, err error) {
	code = c.TwoFactorCode

	if code == "" {
		code, err = twoFA.GenerateAuthCode(c.SharedSecret, nil)
		if err != nil {
			err = errors.ErrInvalidSharedSecret
		}
	}

	return
}

func (c Account) GenerateLogOnDetails() (steam.LogOnDetails, error) {
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
