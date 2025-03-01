package requests

import (
	"errors"
	"sync"

	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
)

type RequestHandler interface {
	AddRequest(A int) *sync.WaitGroup
	FinishRequest(A int) error
}

type request struct {
	paramA int
	wg     *sync.WaitGroup
}

type requestHandler struct {
	inAirRequests map[int]request
	mu            sync.Mutex
}

func (r *requestHandler) AddRequest(A int) *sync.WaitGroup {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.inAirRequests[A]; exists {
		logger.ERROR.Printf("Tried to add request %v, but it allready exists", A)
		return nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	r.inAirRequests[A] = request{
		paramA: A,
		wg:     wg,
	}

	return wg
}

func (r *requestHandler) FinishRequest(A int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	req, exists := r.inAirRequests[A]
	if !exists {
		logger.ERROR.Printf("Tried to finish %v, but it does not exist", A)
		return errors.New("request not found")
	}

	req.wg.Done()

	delete(r.inAirRequests, A)

	return nil
}

func NewRequestHandler() RequestHandler {
	inAir := make(map[int]request)
	return &requestHandler{inAirRequests: inAir}
}

