package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const orderBaseURL string = "/orders"
const version string = "v1"

func (s *Server) useRoutes() {
	v1 := s.echo.Group("/api/" + version)

	v1.GET(orderBaseURL, s.orderHandler.getOrders)
	v1.GET(orderBaseURL+"/:id", s.orderHandler.getOrder)
	v1.POST(orderBaseURL, s.orderHandler.create)
}

func (s *Server) useHealth() {
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})
}
