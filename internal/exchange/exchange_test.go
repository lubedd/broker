package exchange

import (
	"broker/broker/internal/domain"
	"broker/broker/internal/queues"
	"testing"
)

func getExchange() *Exchange {
	return &Exchange{
		queues: queues.NewQueuesService(),
		signal: make(chan uint8),
	}
}

func TestExchange(t *testing.T) {
	ex := getExchange()
	go func(ex *Exchange) {
		<- ex.signal
	}(ex)

	ex.AddMessage(domain.Message{MessageText: "message text", RoutingKey: "router key"})
	if len(ex.Messages)!=1{
		t.Error("Exchange message not add")
	}

	ex.divisionMessages()
	if len(ex.Messages)!=0{
		t.Error("Exchange division error")
	}
}
