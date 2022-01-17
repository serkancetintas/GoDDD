package application

import (
	"context"
	order "go-practice/internal/order/domain"
	"go-practice/internal/order/infra/store"
)

type OrderService struct {
	OrderRepository *store.OrderMongoRepository
}

func NewOrderService(o *store.OrderMongoRepository) OrderService {
	return OrderService{OrderRepository: o}
}

func (o *OrderService) All(ctx context.Context) ([]*order.Order, error) {
	return o.OrderRepository.GetAll(ctx)
}

func (o *OrderService) FindByID(ctx context.Context, id string) (*order.Order, error) {
	return o.OrderRepository.Get(ctx, id)
}

func (o *OrderService) Create(ctx context.Context, order *order.Order) error {
	return o.OrderRepository.Create(ctx, order)
}

func (o *OrderService) Update(ctx context.Context, order *order.Order) error {
	return o.OrderRepository.Update(ctx, order)
}
