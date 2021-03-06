package api

import (
	"context"
	"go-practice/internal/order/application"

	"github.com/labstack/echo/v4"
)

type Config struct {
	Server struct {
		Port    string `json:"port"`
		Timeout int    `json:"timeout"`
	} `json:"server"`
	MongoDB struct {
		URL      string `json:"url"`
		Database string `json:"database"`
		Timeout  int `json:"timeout"`
	} `json:"mongoDb"`
	Context struct {
		Timeout int `json:"timeout"`
	} `json:"context"`
}

type Server struct {
	cfg                    Config
	echo                   *echo.Echo
	orderService 			application.OrderService
	orderHandler			OrderHandler
}

func NewServer(cfg Config,
	e *echo.Echo,
	o application.OrderService,
	h OrderHandler,
	) *Server {
	server := Server{
		cfg:                    cfg,
		echo:                   e,
		orderService: 			o,
		orderHandler: 			h,
	}

	server.useHealth()
	server.useRoutes()
	server.useMiddlewares()

	return &server
}

func (s *Server) Start() error {
	port := s.cfg.Server.Port
	return s.echo.Start(port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (s *Server) Fatal(err error) { s.echo.Logger.Fatal(err) }

func (s *Server) Config() Config { return s.cfg }
