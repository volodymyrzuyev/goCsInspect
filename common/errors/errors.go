package errors

import "errors"

// Client

var (
	ErrClientUnableToConnect = errors.New("client unable to connect")
	ErrClientTimeout         = errors.New("client timeout when fetching skin")
)

// Credentials

var (
	ErrInvalidCredentials  = errors.New("invalid credentials, username and password and (2FC or SharedSecret) must be provided")
	ErrInvalidSharedSecret = errors.New("provided SharedSecret is invalid")
)

// inspectParams

var (
	ErrInvalidParameters  = errors.New("parameters A and D and (M or S) must be provided")
	ErrInvalidInspectLink = errors.New("pas not able to parse inspectLink")
)
