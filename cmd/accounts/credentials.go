package accounts

import (
	"fmt"

	twoFA "github.com/bbqtd/go-steam-authenticator"
)

type Credentials struct {
	Username      string
	Password      string
	TwoFactorCode string
	SharedSecret  string
}

func (c Credentials) validate() error {
	if c.Username == "" {
		return UsernameNotProvided
	}

	if c.Password == "" {
		return PasswordNotProvided
	}

	if c.TwoFactorCode == "" && c.SharedSecret == "" {
		c.TwoFactorCode = c.request2FA()
	}

	return nil
}

func (c Credentials) get2FA() (string, error) {
	if c.TwoFactorCode != "" {
		return c.TwoFactorCode, nil
	}

	return twoFA.GenerateAuthCode(c.SharedSecret, nil)
}

func (c Credentials) request2FA() string {
	var ret string
	fmt.Printf("%v No 2FA provided. Enter 2FA code: ", c.Username)
	fmt.Scan(&ret)
	return ret
}
