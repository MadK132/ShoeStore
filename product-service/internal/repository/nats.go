package repository

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"shoeshop/product-service/internal/model"
)

type EventPublisher interface {
	PublishProductCreated(product *model.Product) error
	PublishProductUpdated(product *model.Product) error
	PublishProductDeleted(productID string) error
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

func (n *natsClient) PublishProductCreated(product *model.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return n.conn.Publish("product.created", data)
}

func (n *natsClient) PublishProductUpdated(product *model.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return n.conn.Publish("product.updated", data)
}

func (n *natsClient) PublishProductDeleted(productID string) error {
	return n.conn.Publish("product.deleted", []byte(productID))
}

func (n *natsClient) Close() error {
	n.conn.Close()
	return nil
} 