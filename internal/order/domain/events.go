package order

type (
	OrderCreated struct {
		id         string
		orderItems []OrderItem
		customerId string
	}
	PaidEvent struct {
		id string
	}
	CancelledEvent struct {
		id string
	}
)
