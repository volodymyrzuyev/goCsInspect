package clientmanagement

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/golang-collections/collections/queue"
)

type response struct {
	responseProto *protobuf.CEconItemPreviewDataBlock
	err           error
}

type job struct {
	requestProto *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest
	responseCh   chan response
}

type jobQue struct {
	mu        sync.Mutex
	que       *queue.Queue
	clientQue *clientQue
}

func newJobQue(c *clientQue) *jobQue {
	q := &jobQue{que: queue.New(), clientQue: c}

	go q.runQue()

	return q
}

func (q *jobQue) registerJob(proto *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest) <-chan response {
	newJob := job{
		requestProto: proto,
		responseCh:   make(chan response),
	}

	q.mu.Lock()
	q.que.Enqueue(newJob)
	q.mu.Unlock()
	slog.Debug("job added", "item_id", proto.GetParamA())

	return newJob.responseCh
}

const queIdleSleepTime = 75 * time.Millisecond

func (q *jobQue) runQue() {
	for {
		q.mu.Lock()
		if q.que.Len() > 0 {
			j := q.que.Dequeue().(job)
			q.mu.Unlock()
			slog.Debug("processing job", "item_id", j.requestProto.GetParamA())
			q.clientQue.runJob(j)

		} else {
			q.mu.Unlock()
			time.Sleep(queIdleSleepTime)
		}
	}
}
