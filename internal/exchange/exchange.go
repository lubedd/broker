package exchange

import (
	"broker/broker/internal/domain"
	"broker/broker/internal/queues"
)

type ExchangeService interface {
	AddMessage(domain.Message)
}

func NewExchange(queues queues.Queues) ExchangeService {
	return NewExchangeProcess(queues)
}
