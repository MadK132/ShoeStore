package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shoeshop/email-service/internal/service"
	pb "shoeshop/proto"
)

type GRPCHandler struct {
	pb.UnimplementedEmailServiceServer
	emailService service.EmailService
}

func NewGRPCHandler(emailService service.EmailService) *GRPCHandler {
	return &GRPCHandler{
		emailService: emailService,
	}
}

func (h *GRPCHandler) SendRegistrationConfirmation(ctx context.Context, req *pb.UserEmailRequest) (*pb.EmailResponse, error) {
	if err := h.emailService.SendRegistrationConfirmation(ctx, req.GetUser()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send registration confirmation: %v", err)
	}

	return &pb.EmailResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) SendOrderConfirmation(ctx context.Context, req *pb.OrderEmailRequest) (*pb.EmailResponse, error) {
	if err := h.emailService.SendOrderConfirmation(ctx, req.GetOrder()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send order confirmation: %v", err)
	}

	return &pb.EmailResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) SendOrderStatusUpdate(ctx context.Context, req *pb.OrderEmailRequest) (*pb.EmailResponse, error) {
	if err := h.emailService.SendOrderStatusUpdate(ctx, req.GetOrder()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send order status update: %v", err)
	}

	return &pb.EmailResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) SendPasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*pb.EmailResponse, error) {
	if err := h.emailService.SendPasswordReset(ctx, req.GetUser(), req.GetResetToken()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send password reset email: %v", err)
	}

	return &pb.EmailResponse{
		Success: true,
	}, nil
} 