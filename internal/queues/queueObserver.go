package queues

import (
	"broker/broker/internal/domain"
	"time"
)

func (q QueuesService) newObserver(queue *domain.Queue) {
	for signal := range queue.Signal {
		switch signal {
		case domain.SigSendMessages:
			q.sendMessage(queue)
			queue.Timer.Reset(5 * time.Second)
		}
	}
}

func (q QueuesService) sendMessagesByTimer(queue *domain.Queue) {
	queue.Timer = time.NewTimer(time.Second)
	for {
		<-queue.Timer.C
		queue.Signal <- domain.SigSendMessages
	}
}
