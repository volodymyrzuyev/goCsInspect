package client

import (
	"bytes"
	"testing"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

type mock struct{}

func (r mock) RegisterRequest(l uint64) *chan types.Response {
	ch := make(chan types.Response)
	return &ch
}
func (r mock) HandleGCPacket(*gamecoordinator.GCPacket) {}

func validCredentials() types.Credentials {
	return types.Credentials{
		Username:      "Test",
		Password:      "Test",
		TwoFactorCode: "Test",
		SharedSecret:  "SGVsbG9Xb3JsZA==",
	}
}

func TestLogIn(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewLogger(&buf)

	t.Run("Valid", func(t *testing.T) {
		var wantErr error

		eventLoopRunner = func(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
			login <- true
		}
		cli := NewClient(mock{}, log)
		err := cli.LogIn(validCredentials(), mock{})

		assert.Equal(t, wantErr, err, "Everything provided, should not fail")

		eventLoopRunner = runEventLoop
	})

	t.Run("Timeout", func(t *testing.T) {
		wantErr := UnableToLogin

		eventLoopRunner = func(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
			return
		}
		cli := NewClient(mock{}, log)
		err := cli.LogIn(validCredentials(), mock{})

		assert.Equal(t, wantErr, err, "Chan never returns, should exit")

		eventLoopRunner = runEventLoop

	})
}

func TestLogout(t *testing.T) {
	var buffer bytes.Buffer
	log := logger.NewLogger(&buffer)

	t.Run("Valid", func(t *testing.T) {
		ch := make(chan bool)
		eventLoopRunner = func(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
			// tell the system we logged in
			login <- true
			select {
			case <-exit:
				ch <- true
			}
		}
		cli := client{log: log, disconected: true}
		cli.LogIn(validCredentials(), mock{})

		cli.LogOut()
		var ret bool
		select {
		case ret = <-ch:
		}

		assert.True(t, ret, "Event runner should get true")
		assert.True(t, cli.disconected, "Client should get disconected")
		assert.False(t, cli.Avaliable(), "Client should not be avaliable")

	})

}
