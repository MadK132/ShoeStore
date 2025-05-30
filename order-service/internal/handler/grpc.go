package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shoeshop/order-service/internal/model"
	"shoeshop/order-service/internal/service"
	pb "shoeshop/proto"
)

type GRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	orderService service.OrderService
}

func NewGRPCHandler(orderService service.OrderService) *GRPCHandler {
	return &GRPCHandler{
		orderService: orderService,
	}
}

func (h *GRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	order, err := model.FromProto(req.GetOrder())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order data: %v", err)
	}

	createdOrder, err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &pb.OrderResponse{
		Order: createdOrder.ToProto(),
	}, nil
}

func (h *GRPCHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.orderService.GetOrder(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	return &pb.OrderResponse{
		Order: order.ToProto(),
	}, nil
}

func (h *GRPCHandler) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	order, err := model.FromProto(req.GetOrder())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order data: %v", err)
	}

	updatedOrder, err := h.orderService.UpdateOrder(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order: %v", err)
	}

	return &pb.OrderResponse{
		Order: updatedOrder.ToProto(),
	}, nil
}

func (h *GRPCHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.orderService.ListOrders(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = order.ToProto()
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
	}, nil
}

func (h *GRPCHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	err := h.orderService.UpdateOrderStatus(ctx, req.GetId(), model.OrderStatus(req.GetStatus()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &pb.UpdateOrderStatusResponse{
		Success: true,
	}, nil
} 