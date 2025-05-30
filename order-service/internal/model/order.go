package model

import (
	"time"
	pb "shoeshop/proto"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusShipped  OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

type OrderItem struct {
	ProductID string  `bson:"product_id"`
	Quantity  int32   `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type Order struct {
	ID              string      `bson:"_id,omitempty"`
	UserID          string      `bson:"user_id"`
	Items           []OrderItem `bson:"items"`
	TotalAmount     float64     `bson:"total_amount"`
	Status          OrderStatus `bson:"status"`
	ShippingAddress string      `bson:"shipping_address"`
	CreatedAt       time.Time   `bson:"created_at"`
	UpdatedAt       time.Time   `bson:"updated_at"`
	PaymentMethod   string      `bson:"payment_method"`
	PaymentID       string      `bson:"payment_id,omitempty"`
}

// ToProto конвертирует доменную модель в protobuf модель
func (o *Order) ToProto() *pb.Order {
	items := make([]*pb.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:    item.Price,
		}
	}

	return &pb.Order{
		Id:              o.ID,
		UserId:          o.UserID,
		Items:           items,
		TotalAmount:     o.TotalAmount,
		Status:          string(o.Status),
		ShippingAddress: o.ShippingAddress,
		CreatedAt:       o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       o.UpdatedAt.Format(time.RFC3339),
		PaymentMethod:   o.PaymentMethod,
		PaymentId:       o.PaymentID,
	}
}

// FromProto конвертирует protobuf модель в доменную модель
func FromProto(pbOrder *pb.Order) (*Order, error) {
	items := make([]OrderItem, len(pbOrder.Items))
	for i, item := range pbOrder.Items {
		items[i] = OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	createdAt, err := time.Parse(time.RFC3339, pbOrder.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	updatedAt, err := time.Parse(time.RFC3339, pbOrder.UpdatedAt)
	if err != nil {
		updatedAt = time.Now()
	}

	return &Order{
		ID:              pbOrder.Id,
		UserID:          pbOrder.UserId,
		Items:           items,
		TotalAmount:     pbOrder.TotalAmount,
		Status:          OrderStatus(pbOrder.Status),
		ShippingAddress: pbOrder.ShippingAddress,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		PaymentMethod:   pbOrder.PaymentMethod,
		PaymentID:       pbOrder.PaymentId,
	}, nil
} 