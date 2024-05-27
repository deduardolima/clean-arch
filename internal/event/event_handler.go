package event

import (
	"github.com/deduardolima/clean-arch/pkg/events"
)

type EventBundle struct {
	OrderCreatedEvent events.EventInterface
	OrdersListedEvent events.EventInterface
}

func NewEventBundle() *EventBundle {
	return &EventBundle{
		OrderCreatedEvent: NewOrderCreated(),
		OrdersListedEvent: NewOrdersListed(),
	}
}
