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

type responseHandler struct{}

func (r responseHandler) RegisterRequest(l uint64) *chan types.Response {
	return handler(l)
}

func (r responseHandler) HandleGCPacket(*gamecoordinator.GCPacket) {}

var handler = run

func run(l uint64) *chan types.Response {
	ch := make(chan types.Response)
	return &ch
}

func TestLogIn(t *testing.T) {
	var buf string
	log := logger.NewLogger(bytes.NewBufferString(buf))

	workingCreds := func() types.Credentials {
		return types.Credentials{
			Username:      "Test",
			Password:      "Test",
			TwoFactorCode: "Test",
			SharedSecret:  "SGVsbG9Xb3JsZA==",
		}
	}

	t.Run("Valid", func(t *testing.T) {
		eventLoopRunner = func(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
			login <- true
		}
		cli := NewClient(responseHandler{}, log)

		err := cli.LogIn(workingCreds(), responseHandler{})

		assert.Equal(t, nil, err, "Everything provided, should not fail")

		eventLoopRunner = runEventLoop
	})

	t.Run("Timeout", func(t *testing.T) {
		eventLoopRunner = func(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
			return
		}

		cli := NewClient(responseHandler{}, log)

		err := cli.LogIn(workingCreds(), responseHandler{})

		assert.Equal(t, UnableToLogin, err, "Chan never returns, should exit")

		eventLoopRunner = runEventLoop

	})
}

func TestLogout(t *testing.T) {
	var buffer bytes.Buffer
	log := logger.NewLogger(&buffer)
	workingCreds := func() types.Credentials {
		return types.Credentials{
			Username:      "Test",
			Password:      "Test",
			TwoFactorCode: "Test",
			SharedSecret:  "SGVsbG9Xb3JsZA==",
		}
	}

	t.Run("Valid", func(t *testing.T) {
		ch := make(chan bool)
		eventLoopRunner = func(curClient *steam.Client, logInInfo steam.LogOnDetails, login chan bool, log logger.Logger, exit chan bool) {
			login <- true
			select {
			case <-exit:
				ch <- true
			}
		}

		cli := client{log: log, disconected: true}

		cli.LogIn(workingCreds(), responseHandler{})

		cli.LogOut()

		var ret bool
		select {
		case ret = <-ch:
		}

		assert.Equal(t, true, ret, "Event runner should get true")
		assert.Equal(t, true, cli.disconected, "Client should get disconected")
		assert.Equal(t, false, cli.Avaliable(), "Client should not be avaliable")

	})

}
