package exchange

import (
	"broker/broker/internal/domain"
	"broker/broker/internal/queues"
	"sync"
)

type Exchange struct {
	mu       sync.Mutex
	Messages []domain.Message
	queues   queues.Queues
	signal   chan uint8
}

func NewExchangeProcess(queues queues.Queues) *Exchange {
	ex := &Exchange{
		Messages: []domain.Message{},
		queues:   queues,
		signal:   make(chan uint8),
	}
	go ex.newObserver()

	return ex
}

func (e *Exchange) AddMessage(message domain.Message) {
	e.mu.Lock()
	e.Messages = append(e.Messages, message)
	e.signal <- domain.SigNewMessage
	e.mu.Unlock()
}
