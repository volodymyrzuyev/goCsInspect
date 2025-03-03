package accounts

import "fmt"

var (
	PasswordNotProvided = fmt.Errorf("Account password not provided")
	UsernameNotProvided = fmt.Errorf("Account username not provided")
	UnableToLogIn       = fmt.Errorf("Unable to login")
	NonValidParams      = fmt.Errorf("Params are not valid")
	NoAvaliableAccounts = fmt.Errorf("No avaliable accounts")
	TimeOut             = fmt.Errorf("Request timeout")
	InternalError       = fmt.Errorf("Internal Error")
	staleItem           = fmt.Errorf("Stale DB item")
)
