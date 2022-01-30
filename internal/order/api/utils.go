package api

import (
	"context"
	"errors"
	"go-practice/internal/common/aggregate"
	"go-practice/internal/common/config"
	"go-practice/internal/common/rabbit"
	"reflect"

	"github.com/labstack/echo/v4"
)

var ErrInvalidRequest = errors.New("invalid Request params")

func handle(c echo.Context,
	statusCode int,
	fn func(ctx context.Context) error,
	d aggregate.Root,
	b *rabbit.Client) error {
	if err := fn(c.Request().Context()); err != nil {
		return err
	}

	for _, event := range d.Events() {
		msg, _ := rabbit.Serialize(event)
		b.Publish(config.Message{
			Queue: reflect.TypeOf(event).String(),
			Body: config.MessageBody{
				Data: msg,
			},
		})
	}

	return c.JSON(statusCode, "")
}

func handleR(c echo.Context,
	statusCode int,
	fn func(ctx context.Context) (interface{}, error)) error {
	result, err := fn(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(statusCode, result)
}

func CallMethod(i interface{}, methodName string) interface{} {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call([]reflect.Value{})[0].Interface()
	}

	// return or panic, method not found of either type
	return ""
}
