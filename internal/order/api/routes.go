package api

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go-practice/internal/common/aggregate"
	order "go-practice/internal/order/domain"
	"net/http"
)

const orderBaseURL string = "/orders"
const version string = "v1"

func (s *Server) useRoutes() {
	v1 := s.echo.Group("/api/" + version)
	v1.GET(orderBaseURL, func(c echo.Context) error {
		return handleR(c, http.StatusOK, func(ctx context.Context) (interface{}, error) {
			orders, err := s.orderService.All(ctx)

			return orders, err
		})
	})

	v1.GET(orderBaseURL+"/:id", func(c echo.Context) error {
		id := c.Param("id")

		return handleR(c, http.StatusOK, func(ctx context.Context) (interface{}, error) {
			order, err := s.orderService.FindByID(ctx, id)

			return order, err
		})
	})

	v1.POST(orderBaseURL, func(c echo.Context) error {
		var orderDTO OrderDTO
		if err := c.Bind(&orderDTO); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var orderItems []order.OrderItem
		fmt.Println("order dto", orderDTO)
		for _, oi :=  range orderDTO.OrderItems {
			orderItems = append(orderItems,
				order.OrderItem{ ProductID: oi.ProductID,
					ProductName: oi.ProductName,
					ItemCount: oi.ItemCount,
					Price: oi.Price,
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
				return s.orderService.Create(ctx, ordr)
			})
	})

	//v1.POST(orderBaseURL, s.orderCommandController.create)
}

func (s *Server) useHealth() {
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})
}
