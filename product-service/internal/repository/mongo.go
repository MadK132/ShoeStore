package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shoeshop/product-service/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	GetByID(ctx context.Context, id string) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) (*model.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter bson.M) ([]*model.Product, error)
	SearchProducts(ctx context.Context, query string) ([]*model.Product, error)
}

type mongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRepository(uri string) (ProductRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database("shoeshop").Collection("products")
	
	// Создаем индексы
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "brand", Value: 1}},
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

func (r *mongoRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	product.ID = result.InsertedID.(string)
	return product, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *mongoRepository) Update(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.UpdatedAt = time.Now()

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": product.ID},
		bson.M{"$set": product},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedProduct model.Product
	if err := result.Decode(&updatedProduct); err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *mongoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoRepository) List(ctx context.Context, filter bson.M) ([]*model.Product, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *mongoRepository) SearchProducts(ctx context.Context, query string) ([]*model.Product, error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
} 