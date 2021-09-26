package domain

import (
	pb "broker/broker/internal/proto"
	"time"
)

type Queue struct {
	Name      string
	Messages  []string
	Consumers map[uint64]*Consumer
	Signal    chan uint8
	Timer     *time.Timer
}

type Consumer struct {
	WaitSignal bool
	Id         uint64
	Signal     chan uint8
	Stream     pb.Broker_ConsumerChatServer
}
