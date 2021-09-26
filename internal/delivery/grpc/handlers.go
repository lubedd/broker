package grpc

import (
	"broker/broker/internal/domain"
	"broker/broker/internal/exchange"
	"broker/broker/internal/queues"
	pb "broker/broker/internal/proto"
	"context"
)

type Handlers struct {
	exchange exchange.ExchangeService
	queues   queues.Queues
}

func NewHandlers(exchange exchange.ExchangeService, queues queues.Queues) Handlers {
	return Handlers{exchange, queues}
}

func (h *Handlers) AddMessage(c context.Context, request *pb.RequestProducer) (response *pb.ResponseProducer, err error) {
	message := domain.Message{
		MessageText: request.MessageText,
		RoutingKey:  request.RoutingKey,
	}

	h.exchange.AddMessage(message)
	return &pb.ResponseProducer{Message: domain.ProducerTaskAccepted}, nil
}

func (h *Handlers) ConsumerChat(stream pb.Broker_ConsumerChatServer) (err error) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		switch resp.Message {
		case domain.ConsumersOpenConnection:
			idConsumer := h.queues.AddConsumer(resp.RoutingKey, stream)
			if err := stream.Send(&pb.Consumer{RoutingKey: resp.RoutingKey, Message: resp.Message, Id: idConsumer}); err != nil {
				break
			}
		case domain.ConsumersTaskAccepted:
			h.queues.ConsumerTaskAccepted(resp.RoutingKey, resp.Id)
		}
	}

	return nil
}
