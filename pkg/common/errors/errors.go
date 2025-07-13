package errors

import "errors"

// Client
var (
	ErrClientUnableToConnect = errors.New("client unable to connect")
	ErrClientTimeout         = errors.New("client timeout when fetching skin")
	ErrClientUnavailable     = errors.New("client is unavailable now")
)

// clientManager
var (
	ErrNoAvailableClients   = errors.New("no available clients")
	ErrInvalidManagerConfig = errors.New("detailer and storage are needed")
)

// Credentials
var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInsufficientCredentials = errors.Join(
		ErrInvalidCredentials,
		errors.New("username and password and (2FC or SharedSecret) must be provided"),
	)
	ErrInvalidSharedSecret = errors.Join(
		ErrInvalidCredentials,
		errors.New("provided SharedSecret is invalid"),
	)
)

// inspectParams
var (
	ErrInvalidParameters  = errors.New("parameters A and D and (M or S) must be provided")
	ErrInvalidInspectLink = errors.New("was not able to parse inspectLink")
)

// detailer
var (
	ErrUnknownProtoValue = errors.New("item proto has unknown properties")
	ErrUnknownRarity     = errors.Join(ErrUnknownProtoValue, errors.New("unknown rarity"))
	ErrUnknownQuality    = errors.Join(ErrUnknownProtoValue, errors.New("unknown rarity"))
	ErrUnknownDefIndex   = errors.Join(ErrUnknownProtoValue, errors.New("unknown defIndex"))
	ErrUnknownPaintIndex = errors.Join(ErrUnknownProtoValue, errors.New("unknown paintIndex"))
	ErrUnknownMusicIndex = errors.Join(
		ErrUnknownProtoValue,
		errors.New("unknown musicKit index"),
	)
	ErrUnknownStickerModifier  = errors.Join(ErrUnknownProtoValue, errors.New("unknown sticker"))
	ErrUnknownKeychainModifier = errors.Join(ErrUnknownProtoValue, errors.New("unknown sticker"))
)

// storage

var (
	ErrDB = errors.New("database error")

	ErrDBInsert       = errors.Join(ErrDB, errors.New("err inserting"))
	ErrInsertItem     = errors.Join(ErrDBInsert, errors.New("err inserting item"))
	ErrInsertSticker  = errors.Join(ErrDBInsert, errors.New("err inserting sticker"))
	ErrInsertKeychain = errors.Join(ErrDBInsert, errors.New("err inserting keychain"))

	ErrDBFetch       = errors.Join(ErrDB, errors.New("err fetching"))
	ErrFetchItem     = errors.Join(ErrDBFetch, errors.New("err fetching item"))
	ErrFetchSticker  = errors.Join(ErrDBFetch, errors.New("err fetching sticker"))
	ErrFetchKeychain = errors.Join(ErrDBFetch, errors.New("err fetching keychain"))
)
