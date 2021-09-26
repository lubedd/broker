package queues

import pb "broker/broker/internal/proto"

type Queues interface {
	AddMessage(queueName string, message string)
	AddConsumer(queueName string, consumer pb.Broker_ConsumerChatServer) uint64
	ConsumerTaskAccepted(queueName string, id uint64)
}

func NewQueues() Queues {
	return NewQueuesService()
}
