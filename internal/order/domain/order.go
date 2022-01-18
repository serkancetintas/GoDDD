package order

import (
	aggregate "go-practice/internal/common/aggregate"
	"time"
)

type ID string

type Order struct {
	aggregate.Root
	id          ID
	orderItems  []OrderItem
	customerID	string
	createdTime time.Time
	status      Status
	version     aggregate.Version
}

func NewOrder(id ID, orderItems []OrderItem,
		status Status, version aggregate.Version,
		customerID string) (*Order, error) {
	o := Order{
		id: id,
		orderItems: orderItems,
		createdTime: time.Now(),
		version: version,
		status: status,
		customerID: customerID,
	}

	if err := valid(&o); err != nil {
		return nil, err
	}

	o.AddEvent(CreatedEvent{id: string(id)})

	return &o, nil
}

func (o *Order) Pay() {
	o.status = Paid
	o.AddEvent(PaidEvent{id: string(o.id)})
}

func (o *Order) Cancel() {
	o.status = Canceled
	o.AddEvent(CancelledEvent{id: string(o.id)})
}

func (o *Order) Status() Status {  return o.status }

func valid(o *Order) error {
	if o.id == "" || o.orderItems == nil || o.customerID == "" {
		return ErrInvalidValue
	}

	return nil
}

func (o *Order) ID() string { return string(o.id) }

func (o *Order) Version() string { return o.version.String() }

func (o *Order) CustomerId() string { return o.customerID }

func (o *Order) CreatedTime() time.Time { return o.createdTime }

func (o *Order) OrderItems() *[]OrderItem {
	return &o.orderItems
}