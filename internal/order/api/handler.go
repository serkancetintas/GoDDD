package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go-practice/internal/common/aggregate"
	"go-practice/internal/common/rabbit"
	"go-practice/internal/order/application"
	order "go-practice/internal/order/domain"
	"net/http"
)

type OrderHandler struct {
	OrderService application.OrderService
	EventBus     *rabbit.Client
}

func NewOrderHandler(o application.OrderService, r *rabbit.Client) OrderHandler {
	return OrderHandler{o, r}
}

func (o OrderHandler) getOrders(c echo.Context) error {
	return handleR(c, http.StatusOK, func(ctx context.Context) (interface{}, error) {
		orders, err := o.OrderService.All(ctx)

		return orders, err
	})
}

func (o OrderHandler) getOrder(c echo.Context) error {
	id := c.Param("id")

	return handleR(c, http.StatusOK, func(ctx context.Context) (interface{}, error) {
		order, err := o.OrderService.FindByID(ctx, id)

		return order, err
	})
}

func (o OrderHandler) create(c echo.Context) error {
	var orderDTO OrderDTO
	if err := c.Bind(&orderDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var orderItems []order.OrderItem
	for _, oi := range orderDTO.OrderItems {
		orderItems = append(orderItems,
			order.OrderItem{ProductID: oi.ProductID,
				ProductName: oi.ProductName,
				ItemCount:   oi.ItemCount,
				Price:       oi.Price,
			})
	}

	ordr, err := order.NewOrder(order.ID(orderDTO.ID),
		orderItems,
		order.Submitted,
		aggregate.NewVersion(),
		orderDTO.CustomerID)

	if err != nil {
		return errors.Wrap(err, "create order")
	}

	return handle(c,
		http.StatusCreated,
		func(ctx context.Context) error {
			return o.OrderService.Create(ctx, ordr)
		}, ordr.Root, o.EventBus)
}
