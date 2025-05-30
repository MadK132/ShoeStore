package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"shoeshop/user-service/internal/handler"
	"shoeshop/user-service/internal/repository"
	"shoeshop/user-service/internal/service"
	pb "shoeshop/proto"
)

func main() {
	// Настройка логгера
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Инициализация репозитория (работа с БД)
	repo, err := repository.NewMongoRepository("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Инициализация сервиса (бизнес-логика)
	svc := service.NewUserService(repo)

	// Инициализация gRPC handler
	grpcHandler := handler.NewGRPCHandler(svc)

	// Создание gRPC сервера
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, grpcHandler)
	
	// Включаем reflection для удобства отладки (позволяет использовать grpcurl и другие инструменты)
	reflection.Register(server)

	// Запуск gRPC сервера
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Канал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	sig := <-sigChan
	fmt.Printf("\nReceived signal %v, initiating graceful shutdown\n", sig)

	// Graceful shutdown
	server.GracefulStop()
	log.Println("Server stopped gracefully")
} 