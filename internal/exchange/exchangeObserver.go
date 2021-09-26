package exchange

import (
	"broker/broker/internal/domain"
)

func (e *Exchange) newObserver() {
	for signal := range e.signal {
		switch signal {
		case domain.SigNewMessage:
			e.divisionMessages()
		}
	}
}
