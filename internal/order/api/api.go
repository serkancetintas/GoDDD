package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go-practice/internal/common/aggregate"
	"go-practice/internal/order/application"
	order "go-practice/internal/order/domain"
	"net/http"
)

type OrderAPI struct {
	OrderService application.OrderService
}

func NewOrderAPI(o application.OrderService) OrderAPI {
	return OrderAPI{o}
}

func (o OrderAPI) getOrders(c echo.Context)  error {
	return handleR(c, http.StatusOK, func(ctx context.Context) (interface{}, error) {
		orders, err := o.OrderService.All(ctx)

		return orders, err
	})
}

func (o OrderAPI) getOrder(c echo.Context) error {
	id := c.Param("id")

	return handleR(c, http.StatusOK, func(ctx context.Context) (interface{}, error) {
		order, err := o.OrderService.FindByID(ctx, id)

		return order, err
	})
}

func (o OrderAPI) create(c echo.Context) error {
	orderDTO := new(OrderDTO)
	if err := c.Bind(orderDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var orderItems []order.OrderItem
	for _, oi :=  range orderDTO.OrderItems {
		orderItems = append(orderItems,
			order.OrderItem{ ProductID: oi.ProductID,
				             ProductName: oi.ProductName,
							 ItemCount: oi.ItemCount,
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
		})
}

