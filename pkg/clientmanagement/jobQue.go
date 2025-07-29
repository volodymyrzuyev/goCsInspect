package clientmanagement

import (
	"context"
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
	ctx          context.Context
}

type jobQue struct {
	mu  sync.Mutex
	que *queue.Queue

	clientQue *clientRing
	l         *slog.Logger
}

func newJobQue(c *clientRing) *jobQue {
	q := &jobQue{
		que:       queue.New(),
		clientQue: c,
		l:         slog.Default().WithGroup("ClientManagment.JobQue"),
	}

	go q.runQue()

	return q
}

func (j *jobQue) registerJob(
	proto *protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest,
	ctx context.Context,
) <-chan response {
	newJob := job{
		requestProto: proto,
		responseCh:   make(chan response),
		ctx:          ctx,
	}

	j.mu.Lock()
	j.que.Enqueue(newJob)
	j.mu.Unlock()
	j.l.Debug("job added", "item_id", proto.GetParamA())

	return newJob.responseCh
}

const queIdleSleepTime = 75 * time.Millisecond

func (j *jobQue) runQue() {
	for {
		j.mu.Lock()
		if j.que.Len() == 0 {
			j.mu.Unlock()
			time.Sleep(queIdleSleepTime)
			continue
		}

		job := j.que.Dequeue().(job)
		j.mu.Unlock()
		j.l.Debug("processing job", "item_id", job.requestProto.GetParamA())
		select {
		case <-job.ctx.Done():
			close(job.responseCh)
			j.l.Debug("job expired", "item_id", job.requestProto.GetParamA())
			continue

		default:
			j.clientQue.runJob(job)
		}
	}
}
