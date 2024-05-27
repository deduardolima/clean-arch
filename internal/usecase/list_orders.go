package usecase

import (
	"github.com/deduardolima/clean-arch/internal/entity"
	"github.com/deduardolima/clean-arch/internal/event"
	"github.com/deduardolima/clean-arch/pkg/events"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrdersListed    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	EventBundle *event.EventBundle,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
		OrdersListed:    EventBundle.OrdersListedEvent,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrdersUseCase) Execute() ([]entity.Order, error) {
	orders, err := c.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	c.OrdersListed.SetPayload(orders)
	c.EventDispatcher.Dispatch(c.OrdersListed)

	return orders, nil
}
