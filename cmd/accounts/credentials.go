package accounts

import twoFA "github.com/bbqtd/go-steam-authenticator"

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
		return No2FAmechanistProvided
	}

	return nil
}

func (c Credentials) get2FA() (string, error) {
	if c.TwoFactorCode != "" {
		return c.TwoFactorCode, nil
	}

	return twoFA.GenerateAuthCode(c.SharedSecret, nil)
}
