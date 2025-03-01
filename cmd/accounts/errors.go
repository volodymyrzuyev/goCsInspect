package accounts

import "fmt"

var (
	PasswordNotProvided    = fmt.Errorf("Account password not provided")
	UsernameNotProvided    = fmt.Errorf("Account username not provided")
	No2FAmechanistProvided = fmt.Errorf("2FA code, or ")
	UnableToLogIn          = fmt.Errorf("Unable to login")
	NonValidParams         = fmt.Errorf("Params are not valid")
	NoAvaliableAccounts    = fmt.Errorf("No avaliable accounts")
	LoginTimeOut           = fmt.Errorf("Login timeout")
)
