package queues

import (
	"broker/broker/internal/domain"
	pb "broker/broker/internal/proto"
	"sync"
)

type QueuesService struct {
	mu     sync.Mutex
	Queues map[string]*domain.Queue
}

func NewQueuesService() *QueuesService {
	return &QueuesService{
		Queues: make(map[string]*domain.Queue),
	}
}

func (q *QueuesService) newQueue(queueName string) {
	q.Queues[queueName] = &domain.Queue{
		Name:      queueName,
		Signal:    make(chan uint8),
		Consumers: make(map[uint64]*domain.Consumer),
	}

	go q.newObserver(q.Queues[queueName])
	go q.sendMessagesByTimer(q.Queues[queueName])
}

func (q *QueuesService) AddMessage(queueName, message string) {
	q.mu.Lock()
	if _, exist := q.Queues[queueName]; !exist {
		q.newQueue(queueName)
	}
	q.Queues[queueName].Messages = append(q.Queues[queueName].Messages, message)

	q.mu.Unlock()
}

func (q *QueuesService) AddConsumer(queueName string, stream pb.Broker_ConsumerChatServer) uint64 {
	q.mu.Lock()
	if _, exist := q.Queues[queueName]; !exist {
		q.newQueue(queueName)
	}

	id := uint64(len(q.Queues[queueName].Consumers)) + 1
	q.Queues[queueName].Consumers[id] = &domain.Consumer{
		Id:     id,
		Stream: stream,
		Signal: make(chan uint8),
	}
	q.mu.Unlock()
	return id
}

func (q *QueuesService) ConsumerTaskAccepted(queueName string, id uint64) {
	q.mu.Lock()
	if queue, exist := q.Queues[queueName]; exist {
		if consumer, exist := queue.Consumers[id]; exist {
			if consumer.WaitSignal {
				consumer.Signal <- domain.SigConsumerAcceptTask
			}
		}
	}
	q.mu.Unlock()
}