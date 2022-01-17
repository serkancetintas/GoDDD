package api

type OrderDTO struct {
	ID 		    string 		    `json:"id"`
	CustomerID 	string 			`json:"customerId"`
	OrderItems  []OrderItemDTO 	`json:"orderItems"`
}

type OrderItemDTO struct {
	ProductID	string	`json:"productId"`
	ProductName	string	`json:"productName"`
	ItemCount	int		`json:"itemCount"`
	Price		float64	`json:"price"`
}
