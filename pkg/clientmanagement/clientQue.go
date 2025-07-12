package clientmanagement

import (
	"log/slog"
	"sync"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/volodymyrzuyev/goCsInspect/internal/client"
)

type clientQue struct {
	mu             sync.Mutex
	que            *queue.Queue
	clientCooldown time.Duration
}

func newClientQue(c time.Duration) *clientQue {
	return &clientQue{que: queue.New(), clientCooldown: c}
}

func (c *clientQue) runJob(j job) {
	len := c.len()

	for range len {
		c.mu.Lock()
		cli := c.que.Dequeue().(client.InspectClient)
		c.que.Enqueue(cli)
		c.mu.Unlock()

		if cli.IsAvailable() {
			slog.Debug("Reqiesting response proto", "item_id", j.requestProto.GetParamA(), "client", cli.Username())
			resp, err := cli.InspectItem(j.requestProto)
			j.responseCh <- response{responseProto: resp, err: err}
			return
		}
	}

	slog.Debug("no clients to resolve the job, sleeping", "item_id", j.requestProto.GetParamA())
	time.Sleep(c.clientCooldown)
	c.runJob(j)
}

func (c *clientQue) addClient(cli client.InspectClient) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.que.Enqueue(cli)
}

func (c *clientQue) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.que.Len()
}
