//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/deduardolima/clean-arch/internal/entity"
	"github.com/deduardolima/clean-arch/internal/event"
	"github.com/deduardolima/clean-arch/internal/infra/database"
	"github.com/deduardolima/clean-arch/internal/infra/web"
	"github.com/deduardolima/clean-arch/internal/usecase"
	"github.com/deduardolima/clean-arch/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setEventBundle = wire.NewSet(
	event.NewEventBundle,
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setEventBundle,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setEventBundle,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setEventBundle,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
