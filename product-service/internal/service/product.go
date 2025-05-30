package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"shoeshop/product-service/internal/model"
	"shoeshop/product-service/internal/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, filter map[string]interface{}) ([]*model.Product, error)
	SearchProducts(ctx context.Context, query string) ([]*model.Product, error)
}

type productService struct {
	repo      repository.ProductRepository
	cache     repository.Cache
	publisher repository.EventPublisher
}

func NewProductService(repo repository.ProductRepository, cache repository.Cache, publisher repository.EventPublisher) ProductService {
	return &productService{
		repo:      repo,
		cache:     cache,
		publisher: publisher,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	// Создаем продукт в БД
	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Кэшируем продукт
	if s.cache != nil {
		if err := s.cache.Set(ctx, createdProduct.ID, createdProduct); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			fmt.Printf("failed to cache product: %v\n", err)
		}
	}

	// Публикуем событие
	if err := s.publisher.PublishProductCreated(createdProduct); err != nil {
		fmt.Printf("failed to publish product created event: %v\n", err)
	}

	return createdProduct, nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	// Пробуем получить из кэша
	if s.cache != nil {
		if product, err := s.cache.Get(ctx, id); err == nil && product != nil {
			return product, nil
		}
	}

	// Получаем из БД
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Кэшируем
	if s.cache != nil {
		if err := s.cache.Set(ctx, id, product); err != nil {
			fmt.Printf("failed to cache product: %v\n", err)
		}
	}

	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	// Обновляем в БД
	updatedProduct, err := s.repo.Update(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Обновляем кэш
	if s.cache != nil {
		if err := s.cache.Set(ctx, product.ID, updatedProduct); err != nil {
			fmt.Printf("failed to update product in cache: %v\n", err)
		}
	}

	// Публикуем событие
	if err := s.publisher.PublishProductUpdated(updatedProduct); err != nil {
		fmt.Printf("failed to publish product updated event: %v\n", err)
	}

	return updatedProduct, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	// Удаляем из БД
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Удаляем из кэша
	if s.cache != nil {
		if err := s.cache.Delete(ctx, id); err != nil {
			fmt.Printf("failed to delete product from cache: %v\n", err)
		}
	}

	// Публикуем событие
	if err := s.publisher.PublishProductDeleted(id); err != nil {
		fmt.Printf("failed to publish product deleted event: %v\n", err)
	}

	return nil
}

func (s *productService) ListProducts(ctx context.Context, filter map[string]interface{}) ([]*model.Product, error) {
	// Конвертируем фильтр в BSON
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	return s.repo.List(ctx, bsonFilter)
}

func (s *productService) SearchProducts(ctx context.Context, query string) ([]*model.Product, error) {
	return s.repo.SearchProducts(ctx, query)
} 