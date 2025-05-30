package repository

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"shoeshop/order-service/internal/model"
)

type EventPublisher interface {
	PublishOrderCreated(order *model.Order) error
	PublishOrderUpdated(order *model.Order) error
	PublishOrderStatusChanged(orderID string, status model.OrderStatus) error
	Close() error
}

type natsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (EventPublisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &natsClient{
		conn: nc,
	}, nil
}

func (n *natsClient) PublishOrderCreated(order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return n.conn.Publish("order.created", data)
}

func (n *natsClient) PublishOrderUpdated(order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return n.conn.Publish("order.updated", data)
}

func (n *natsClient) PublishOrderStatusChanged(orderID string, status model.OrderStatus) error {
	data, err := json.Marshal(map[string]interface{}{
		"order_id": orderID,
		"status":   status,
	})
	if err != nil {
		return err
	}
	return n.conn.Publish("order.status_changed", data)
}

func (n *natsClient) Close() error {
	n.conn.Close()
	return nil
} 