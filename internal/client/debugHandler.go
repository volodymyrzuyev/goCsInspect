package client

import (
	"encoding/hex"
	"log/slog"
	"reflect"

	"github.com/Philipp15b/go-steam/v3/protocol"
	"github.com/davecgh/go-spew/spew"
)

type debug struct {
	username     string
	eventLogger  *slog.Logger
	packetLogger *slog.Logger

	packetId uint64
	eventId  uint64
}

func newDebug(username string, eventLogger, packetLogger *slog.Logger) *debug {
	return &debug{
		username:     username,
		eventLogger:  eventLogger,
		packetLogger: packetLogger,

		packetId: 0,
		eventId:  0,
	}
}

func (d *debug) HandlePacket(packet *protocol.Packet) {
	d.packetId++
	text := packet.String() + "\n\n" + hex.Dump(packet.Data)
	d.packetLogger.Debug("Got packet", "username", d.username, "packet_id", d.packetId, "packet_EMsg", packet.EMsg, "data", text)
}

func (d *debug) HandleEvent(event any) {
	d.eventId++
	d.eventLogger.Debug("Got event", "username", d.username, "event_id", d.eventId, "event_name", name(event), "data", []byte(spew.Sdump(event)))
}

func name(obj any) string {
	val := reflect.ValueOf(obj)
	ind := reflect.Indirect(val)
	if ind.IsValid() {
		return ind.Type().Name()
	} else {
		return val.Type().Name()
	}
}
