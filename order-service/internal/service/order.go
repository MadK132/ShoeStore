package service

import (
	"context"
	"fmt"

	"shoeshop/order-service/internal/model"
	"shoeshop/order-service/internal/repository"
	pb "shoeshop/proto"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	ListOrders(ctx context.Context, userID string) ([]*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error
}

type orderService struct {
	repo          repository.OrderRepository
	publisher     repository.EventPublisher
	productClient pb.ProductServiceClient
	userClient    pb.UserServiceClient
	emailClient   pb.EmailServiceClient
}

func NewOrderService(
	repo repository.OrderRepository,
	publisher repository.EventPublisher,
	productClient pb.ProductServiceClient,
	userClient pb.UserServiceClient,
	emailClient pb.EmailServiceClient,
) OrderService {
	return &orderService{
		repo:          repo,
		publisher:     publisher,
		productClient: productClient,
		userClient:    userClient,
		emailClient:   emailClient,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Проверяем существование пользователя
	_, err := s.userClient.GetUser(ctx, &pb.GetUserRequest{Id: order.UserID})
	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	// Проверяем наличие товаров и их цены
	var totalAmount float64
	for _, item := range order.Items {
		product, err := s.productClient.GetProduct(ctx, &pb.GetProductRequest{Id: item.ProductID})
		if err != nil {
			return nil, fmt.Errorf("failed to get product %s: %w", item.ProductID, err)
		}

		if product.Product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}

		item.Price = product.Product.Price
		totalAmount += item.Price * float64(item.Quantity)
	}

	order.TotalAmount = totalAmount
	order.Status = model.OrderStatusPending

	// Создаем заказ
	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Публикуем событие
	if err := s.publisher.PublishOrderCreated(createdOrder); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("failed to publish order created event: %v\n", err)
	}

	// Отправляем email о создании заказа
	_, err = s.emailClient.SendOrderConfirmation(ctx, &pb.OrderEmailRequest{
		Order: createdOrder.ToProto(),
	})
	if err != nil {
		fmt.Printf("failed to send order confirmation email: %v\n", err)
	}

	return createdOrder, nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	updatedOrder, err := s.repo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	if err := s.publisher.PublishOrderUpdated(updatedOrder); err != nil {
		fmt.Printf("failed to publish order updated event: %v\n", err)
	}

	return updatedOrder, nil
}

func (s *orderService) ListOrders(ctx context.Context, userID string) ([]*model.Order, error) {
	return s.repo.List(ctx, userID)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if err := s.publisher.PublishOrderStatusChanged(id, status); err != nil {
		fmt.Printf("failed to publish order status changed event: %v\n", err)
	}

	// Отправляем email о изменении статуса заказа
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order for email notification: %w", err)
	}

	_, err = s.emailClient.SendOrderStatusUpdate(ctx, &pb.OrderEmailRequest{
		Order: order.ToProto(),
	})
	if err != nil {
		fmt.Printf("failed to send order status update email: %v\n", err)
	}

	return nil
} 