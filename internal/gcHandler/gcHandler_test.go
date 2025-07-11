package gcHandler

import (
	"io"
	"log/slog"
	"testing"
	"time"

	csProto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
)

func TestStoreResponse(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	itemId := uint64(1337)
	dummyProto := &csProto.CEconItemPreviewDataBlock{Itemid: &itemId}

	gcHandlerI := NewGcHandler(10 * time.Second)
	gc := gcHandlerI.(*gcHandler)

	cleanUp := func(itemId uint64) {
		gc.mu.Lock()

		delete(gc.pendingResponses, itemId)
		delete(gc.responses, itemId)

		gc.mu.Unlock()
	}

	t.Run("Response Received Before Data Requested", func(t *testing.T) {
		gc.storeResponse(dummyProto)

		assert.Equal(t, dummyProto, gc.responses[itemId], "should be the same")

		cleanUp(itemId)
	})

	t.Run("Response Received After Data Requested", func(t *testing.T) {
		ch := make(chan *csProto.CEconItemPreviewDataBlock)
		gc.mu.Lock()
		gc.pendingResponses[itemId] = ch
		gc.mu.Unlock()

		go gc.storeResponse(dummyProto)

		select {
		case resp := <-ch:
			assert.Equal(t, dummyProto, resp, "should be the same")
		case <-time.After(time.Millisecond):
			t.Fatal("did not get response after a millisecond")
		}

		cleanUp(itemId)
	})
}

func TestGetResponse(t *testing.T) {
	itemId := uint64(1337)
	dummyProto := &csProto.CEconItemPreviewDataBlock{Itemid: &itemId}

	gcHandlerI := NewGcHandler(1 * time.Second)
	gc := gcHandlerI.(*gcHandler)

	t.Run("Response Received Before Data Requested", func(t *testing.T) {
		gc.storeResponse(dummyProto)

		gc.responses[itemId] = dummyProto

		response, err := gc.GetResponse(itemId)

		_, ok := gc.responses[itemId]

		assert.Equal(t, dummyProto, response, "should be the same")
		assert.Nil(t, err, "there should be no error")
		assert.False(t, ok, "the map entry for itemId should not exist")
	})

	t.Run("Response Received After Data Requested", func(t *testing.T) {
		go func() {
			time.Sleep(1 * time.Millisecond)
			gc.mu.Lock()
			gc.pendingResponses[itemId] <- dummyProto
			gc.mu.Unlock()
		}()

		response, err := gc.GetResponse(itemId)

		assert.Equal(t, dummyProto, response, "should be the same")
		assert.Nil(t, err, "there should be no error")
		delete(gc.pendingResponses, itemId)
	})

	t.Run("No Response", func(t *testing.T) {
		response, err := gc.GetResponse(itemId)

		assert.Nil(t, response, "response should be nil")
		assert.Equal(t, errors.ErrClientTimeout, err, "there should be an error")
	})
}
