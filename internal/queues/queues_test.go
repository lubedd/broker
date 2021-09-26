package queues

import (
	"broker/broker/internal/domain"
	mock_proto "broker/broker/internal/proto/mocks"
	"github.com/golang/mock/gomock"
	"log"
	"testing"
	"time"
)

func getQueue() *QueuesService {
	return &QueuesService{
		Queues: make(map[string]*domain.Queue),
	}
}

func TestNewQueue(t *testing.T) {
	log.Println("TestNewQueue")
	q := getQueue()
	queueName := "Queue1"

	q.newQueue(queueName)
	if _, exist := q.Queues[queueName]; !exist {
		t.Error("Queue not add")
	}
	log.Println("OK")
}

func TestQueue(t *testing.T) {
	log.Println("TestQueue")

	q := getQueue()
	ctrl := gomock.NewController(t)
	Broker_ConsumerChatServer := mock_proto.NewMockBroker_ConsumerChatServer(ctrl)
	queueName := "Queue1"
	q.Queues[queueName] = &domain.Queue{
		Name:      queueName,
		Signal:    make(chan uint8),
		Consumers: make(map[uint64]*domain.Consumer),
	}

	q.AddMessage(queueName, "message")
	if len(q.Queues[queueName].Messages) != 1 {
		t.Error("Queue message not add")
	}

	idConsumer := q.AddConsumer(queueName, Broker_ConsumerChatServer)
	if len(q.Queues[queueName].Consumers) != 1 {
		t.Error("Queue message not add")
	}

	acceptSihnal := make(chan uint8)
	go func(q *QueuesService, queueName string, idConsumer uint64) {
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-q.Queues[queueName].Consumers[idConsumer].Signal:
			acceptSihnal <- 1
			return
		case <-timer.C:
			t.Error("Queue consumer not accepted task")
		}
		acceptSihnal <- 1
	}(q, queueName, idConsumer)

	q.Queues[queueName].Consumers[idConsumer].WaitSignal = true
	q.ConsumerTaskAccepted(queueName, idConsumer)
	<-acceptSihnal
	q.sendMessage(q.Queues[queueName])
	log.Println("OK")
}
