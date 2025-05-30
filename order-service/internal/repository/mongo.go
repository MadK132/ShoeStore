package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shoeshop/order-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	GetByID(ctx context.Context, id string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) (*model.Order, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error
}

type mongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRepository(uri string) (OrderRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database("shoeshop").Collection("orders")
	
	// Создаем индексы
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return nil, err
	}

	return &mongoRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *mongoRepository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}

	order.ID = result.InsertedID.(string)
	return order, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *mongoRepository) Update(ctx context.Context, order *model.Order) (*model.Order, error) {
	order.UpdatedAt = time.Now()

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": order.ID},
		bson.M{"$set": order},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedOrder model.Order
	if err := result.Decode(&updatedOrder); err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}

func (r *mongoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoRepository) List(ctx context.Context, userID string) ([]*model.Order, error) {
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*model.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *mongoRepository) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
} 