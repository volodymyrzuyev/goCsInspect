package errors

import "errors"

// Client

var (
	ErrClientUnableToConnect = errors.New("client unable to connect")
)

// Credentials

var (
	ErrInvalidCredentials   = errors.New("invalid credentials, username and password and (2FC or SharedSecret) must be provided")
	ErrInvalidSharedSecret = errors.New("provided SharedSecret is invalid")
)
