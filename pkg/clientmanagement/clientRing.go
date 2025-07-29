package clientmanagement

import (
	"container/ring"
	"log/slog"
	"sync"
	"time"

	"github.com/volodymyrzuyev/goCsInspect/pkg/client"
)

type clientRing struct {
	mu   sync.Mutex
	ring *ring.Ring

	clientCooldown time.Duration
	l              *slog.Logger
}

func newClientQue(c time.Duration) *clientRing {
	return &clientRing{
		clientCooldown: c,
		l:              slog.Default().WithGroup("ClientManagment.ClientQue"),
	}
}

func (r *clientRing) runJob(j job) {
	len := r.len()

	for range len {
		r.mu.Lock()
		cli := r.ring.Value.(client.Client)
		r.ring = r.ring.Next()
		r.mu.Unlock()

		if cli.IsAvailable() {
			r.l.Debug(
				"requesting preview block",
				"item_id",
				j.requestProto.GetParamA(),
				"client",
				cli.Username(),
			)
			resp, err := cli.InspectItem(j.ctx, j.requestProto)
			j.responseCh <- response{responseProto: resp, err: err}
			close(j.responseCh)
			return
		}
	}

	r.l.Debug("no clients to resolve the job, sleeping", "item_id", j.requestProto.GetParamA())
	time.Sleep(r.clientCooldown)
	r.runJob(j)
}

func (r *clientRing) addClient(cli client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newClient := ring.New(1)
	newClient.Value = cli

	if r.ring == nil {
		r.ring = newClient
		return
	}

	r.ring.Link(newClient)
}

func (r *clientRing) len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ring.Len()
}
