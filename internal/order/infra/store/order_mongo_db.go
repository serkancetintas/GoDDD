package store

import (
	"context"
	"go-practice/internal/common/aggregate"
	order "go-practice/internal/order/domain"
	"go-practice/internal/order/infra"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const collectionName = "orders"

type orderBson struct {
	ID			string			`bson:"id"`
	OrderItems	[]orderItemBson	`bson:"orderItems"`
	CustomerID	string			`bson:"customerId"`
	CreatedTime	time.Time		`bson:"createdTime"`
	Status		int				`bson:"status"`
	Version 	string			`bson:"version"`
}

type orderItemBson struct {
	ProductID	string		`bson:"productId"`
	ProductName	string		`bson:"productName"`
	ItemCount	int			`bson:"itemCount"`
	Price 		float64		`bson:"price"`
}

func FromBson(o *orderBson) *order.Order {
	var orderItems []order.OrderItem

	for _, oi := range o.OrderItems {
		orderItems = append(orderItems, FromOrderItemBson(&oi))
	}

	ord, _ := order.NewOrder(order.ID(o.ID),
		orderItems,
		order.Status(o.Status),
		aggregate.Version(o.Version),
		o.CustomerID)

	ord.Clear()

	return ord
}

func FromOrderItemBson(o *orderItemBson) order.OrderItem {
	return order.OrderItem{
		ProductID: o.ProductID,
		ProductName: o.ProductName,
		ItemCount: o.ItemCount,
		Price: o.Price,
	}
}

type OrderMongoRepository struct {
	mStore *infra.MongoStore
}

func NewOrderMongoRepository (mongoStore *infra.MongoStore) *OrderMongoRepository {
	return &OrderMongoRepository{ mStore: mongoStore }
}

func (r *OrderMongoRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	var result []*orderBson

	if err := r.mStore.FindAll(ctx, collectionName, bson.M{}, &result); err != nil {
		return nil, err
	}

	var orders []*order.Order

	for _, o := range result {
		orders = append(orders, FromBson(o))
	}

	return orders, nil
}

func (r *OrderMongoRepository) Get(ctx context.Context, id string) (*order.Order, error) {
	var (
		bsonResult = &orderBson{}
		query      = bson.M{"id": id}
	)

	if err := r.mStore.FindOne(ctx, collectionName, query, nil, bsonResult); err != nil {
		return nil, err
	}

	return FromBson(bsonResult), nil
}

func (r *OrderMongoRepository) Update(ctx context.Context, o *order.Order) error {
	query := bson.M{"id": o.ID(), "version": o.Version()}
	update := bson.M{"$set": bson.M{"status": o.Status(), "version": aggregate.NewVersion().String()}}

	return r.mStore.Update(ctx, collectionName, query, update)
}

func (r *OrderMongoRepository) Create(ctx context.Context, o *order.Order) error {
	bOrder := fromOrder(o)
	if bOrder.Version == "" {
		bOrder.Version = aggregate.NewVersion().String()
	}
	return r.mStore.Store(ctx, collectionName, bOrder)
}

func fromOrder(o *order.Order) *orderBson {
	var orderItemBsons []orderItemBson

	for _,oi := range *o.OrderItems() {
		orderItemBsons = append(orderItemBsons, fromOrderItem(&oi))
	}
	return &orderBson{
		ID:         o.ID(),
		Status:     int(o.Status()),
		CustomerID: o.CustomerId(),
		OrderItems: orderItemBsons,
		Version:    o.Version(),
		CreatedTime: o.CreatedTime(),
	}
}

func fromOrderItem(oi *order.OrderItem) orderItemBson {
	return orderItemBson{
		ProductID: oi.ProductID,
		ProductName: oi.ProductName,
		ItemCount: oi.ItemCount,
		Price: oi.Price,
	}
}




