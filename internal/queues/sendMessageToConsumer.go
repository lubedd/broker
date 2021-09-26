package queues

import (
	"broker/broker/internal/domain"
	pb "broker/broker/internal/proto"
	"errors"
	"time"
)

func (q *QueuesService) sendMessage(queue *domain.Queue) {
	if len(queue.Consumers) == 0 {
		return
	}
	var err error
	for i, message := range queue.Messages {
		for _, consumer := range queue.Consumers {
			err = SendMessagesToConsumer(consumer, domain.Message{RoutingKey: queue.Name, MessageText: message})
			if err != nil {
				delete(queue.Consumers, consumer.Id)
			} else {
				break
			}
		}

		if err == nil {
			queue.Messages = queue.Messages[i+1:]
		}
	}
}

func SendMessagesToConsumer(consumer *domain.Consumer, message domain.Message) error {
	err := consumer.Stream.SendMsg(&pb.Consumer{RoutingKey: message.RoutingKey, Message: message.MessageText, Id: consumer.Id})
	if err != nil {
		return err
	}

	consumer.WaitSignal = true
	t := time.NewTimer(10 * time.Second)
	select {
	case <-consumer.Signal:
		consumer.WaitSignal = false
		break
	case <-t.C:
		consumer.WaitSignal = false
		return errors.New("consumer die")
	}
	return nil
}
