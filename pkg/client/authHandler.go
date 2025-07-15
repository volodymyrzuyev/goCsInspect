package client

import (
	"errors"

	"github.com/Philipp15b/go-steam/v3"
)

type auth struct {
	client        *steam.Client
	details       *steam.LogOnDetails
	logonComplete chan<- error
}

func newAuth(client *steam.Client, details *steam.LogOnDetails, logonComplete chan<- error) *auth {
	return &auth{
		client:        client,
		details:       details,
		logonComplete: logonComplete,
	}
}

func (a *auth) logOn() {
	a.client.Auth.LogOn(a.details)
}

func (a *auth) HandleEvent(event any) {
	switch e := event.(type) {
	case *steam.ConnectedEvent:
		a.logOn()
	case *steam.LoggedOnEvent:
		a.logonComplete <- nil
	case *steam.LogOnFailedEvent:
		a.logonComplete <- errors.New(e.Result.String())
	}
}
