package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shoeshop/user-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.User, error)
}

type mongoRepository struct {
	client            *mongo.Client
	collection        *mongo.Collection
	countersCollection *mongo.Collection
}

func NewMongoRepository(uri string) (UserRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database("shoeshop")
	collection := db.Collection("users")
	countersCollection := db.Collection("counters")

	// Создание уникального индекса для email
	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	// Инициализация счетчика пользователей
	_, err = countersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": "userId"},
		bson.M{"$setOnInsert": bson.M{"value": 0}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return nil, err
	}

	return &mongoRepository{
		client:            client,
		collection:        collection,
		countersCollection: countersCollection,
	}, nil
}

func (r *mongoRepository) getNextUserID(ctx context.Context) (string, error) {
	result := r.countersCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "userId"},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var counter struct {
		Value int `bson:"value"`
	}
	if err := result.Decode(&counter); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", counter.Value), nil
}

func (r *mongoRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	// Получение следующего ID
	nextID, err := r.getNextUserID(ctx)
	if err != nil {
		return nil, err
	}
	user.ID = nextID

	// Установка временных меток
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mongoRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mongoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoRepository) List(ctx context.Context) ([]*model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
} 