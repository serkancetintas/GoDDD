package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go-practice/internal/common/config"
	"go-practice/internal/common/must"
	"go-practice/internal/common/rabbit"
	"go-practice/internal/common/shutdown"
	"go-practice/internal/order/api"
	"go-practice/internal/order/application"
	"go-practice/internal/order/infra"
	"go-practice/internal/order/infra/store"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	configurationManager := config.NewConfigurationManager()
	rabbitConfig := configurationManager.GetRabbitConfig()
	queuesConfig := configurationManager.GetQueuesConfig()

	rabbitClient := rabbit.NewRabbitClient(rabbitConfig, queuesConfig)
	defer rabbitClient.CloseConnection()

	rabbitClient.DeclareExchangeQueueBindings()

	cleanup, err := run(os.Stdout, rabbitClient)
	defer cleanup()

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	shutdown.Gracefully()
}

func run(w io.Writer, r *rabbit.Client) (func(), error) {
	server := buildServer(w, r)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			server.Fatal(errors.New("server could not be started"))
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(server.Config().Context.Timeout)*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Fatal(err)
		}
	}, nil
}

func buildServer(w io.Writer, r *rabbit.Client) *api.Server {
	var cfg api.Config
	readConfig(&cfg)

	mongoStore := infra.NewMongoStore(cfg.MongoDB.URL, cfg.MongoDB.Database, time.Second*time.Duration(cfg.Context.Timeout))
	repository := store.NewOrderMongoRepository(mongoStore)

	service := application.NewOrderService(repository)
	handler := api.NewOrderHandler(service, r)

	e := echo.New()
	e.Logger.SetOutput(w)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return api.NewServer(cfg, e, service, handler)
}

func readConfig(cfg *api.Config) {
	viper.SetConfigFile(`./config.json`)

	must.NotFailF(viper.ReadInConfig)
	must.NotFail(viper.Unmarshal(cfg))
}
