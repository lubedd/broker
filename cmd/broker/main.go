package main

import "broker/broker/internal/app"

// TODO: Парсер json в задачи
// TODO: Exchange type (Funout, Topic)
// TODO: Durable queue, exchange
// TODO: Dead letters pool
// TODO: Web интерфейс с профилированием брокера

func main() {
	app.Run()
}
