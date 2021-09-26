package app

import (
	handlersGrpc "broker/broker/internal/delivery/grpc"
	"broker/broker/internal/exchange"
	pb "broker/broker/internal/proto"
	"broker/broker/internal/queues"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

func Run() {
	serverAddress := flag.String("address", "127.0.0.1:5300", "a string")
	flag.Parse()

	newQueues := queues.NewQueues()
	newExchange := exchange.NewExchange(newQueues)
	handlers := handlersGrpc.NewHandlers(newExchange, newQueues)

	listener, err := net.Listen("tcp", *serverAddress)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterBrokerServer(grpcServer, &handlers)

	fmt.Println("Server listen address:", *serverAddress)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
