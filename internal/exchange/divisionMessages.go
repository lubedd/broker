package exchange

func (e *Exchange) divisionMessages() {
	e.mu.Lock()
	for _, message := range e.Messages {
		e.queues.AddMessage(message.RoutingKey, message.MessageText)
	}
	e.Messages = nil
	e.mu.Unlock()
}
