package requests

import (
	"errors"
	"sync"

	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
)

type RequestHandler interface {
	AddRequest(A int, errChan chan error)
	FinishRequest(A int, erR error) error
}

type request struct {
	paramA int
	ch     *chan error
}

type requestHandler struct {
	inAirRequests map[int]request
	mu            sync.Mutex
}

func (r *requestHandler) AddRequest(A int, errChan chan error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.inAirRequests[A]; exists {
		logger.ERROR.Printf("Tried to add request %v, but it allready exists", A)
		return
	}

	r.inAirRequests[A] = request{
		paramA: A,
		ch:     &errChan,
	}
}

func (r *requestHandler) FinishRequest(A int, erR error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.inAirRequests[A]
	if !exists {
		logger.ERROR.Printf("Tried to finish %v, but it does not exist", A)
		return errors.New("request not found")
	}

	*r.inAirRequests[A].ch <- erR

	delete(r.inAirRequests, A)

	return nil
}

func NewRequestHandler() RequestHandler {
	inAir := make(map[int]request)
	return &requestHandler{inAirRequests: inAir}
}
