package handler

import (
	"context"
	"time"

	"shoeshop/user-service/internal/service"
	"shoeshop/user-service/internal/model"
	pb "shoeshop/proto"
)

type GRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewGRPCHandler(userService service.UserService) *GRPCHandler {
	return &GRPCHandler{
		userService: userService,
	}
}

func (h *GRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	user := &model.User{
		Username:        req.Username,
		Email:          req.Email,
		PasswordHash:   req.Password,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		ShippingAddress: req.ShippingAddress,
	}
	
	createdUser, err := h.userService.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return &pb.UserResponse{
		User: createdUser.ToProto(),
	}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// В реальном приложении здесь бы генерировался JWT токен
	return &pb.LoginResponse{
		Token: "dummy-token",
		User:  user.ToProto(),
	}, nil
}

func (h *GRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	return &pb.UserResponse{
		User: user.ToProto(),
	}, nil
}

func (h *GRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := model.FromProto(req.User)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt = time.Now()
	
	updatedUser, err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return &pb.UserResponse{
		User: updatedUser.ToProto(),
	}, nil
}

func (h *GRPCHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{Success: false}, err
	}
	
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (h *GRPCHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// TODO: Implement password reset logic
	return &pb.ResetPasswordResponse{Success: true}, nil
}

func (h *GRPCHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	
	return &pb.UserResponse{
		User: user.ToProto(),
	}, nil
}

// Добавьте остальные методы gRPC здесь (UpdateUser, DeleteUser, etc.) 