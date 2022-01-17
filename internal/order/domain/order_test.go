package order_test

import (
	"github.com/stretchr/testify/assert"
	"go-practice/internal/common/aggregate"
	"go-practice/internal/order/domain"
	"testing"
)

func TestNewOrder(t *testing.T) {
	o := fakeOrder()

	assert.NotNil(t, o)
}

func TestPayOrder(t *testing.T) {
	o := fakeOrder()

	o.Pay()

	assert.Equal(t, order.Paid, o.Status())
}

func TestCancelOrder(t *testing.T) {
	o := fakeOrder()
	o.Cancel()

	assert.Equal(t, order.Canceled, o.Status())
}

func fakeOrder () *order.Order {
	orderItems := []order.OrderItem{
		{
			"1",
			"Iphone 11",
			1,
			10000,
		},
	}

	o, _ := order.NewOrder("123", orderItems, order.Submitted, aggregate.NewVersion(), "456" )

	return o
}