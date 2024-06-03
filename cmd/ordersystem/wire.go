//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"goexpert-clean-architecture/internal/usecase/order"

	"github.com/google/wire"
	"goexpert-clean-architecture/internal/entity"
	"goexpert-clean-architecture/internal/event"
	"goexpert-clean-architecture/internal/infra/database"
	"goexpert-clean-architecture/internal/infra/rest"
	"goexpert-clean-architecture/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *order.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		order.NewCreateOrderUseCase,
	)
	return &order.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *sql.DB) *order.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		order.NewListOrderUseCase,
	)
	return &order.ListOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *rest.OrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		rest.NewWebOrderHandler,
	)
	return &rest.OrderHandler{}
}
