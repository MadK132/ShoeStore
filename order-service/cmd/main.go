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

	"shoeshop/order-service/internal/handler"
	"shoeshop/order-service/internal/repository"
	"shoeshop/order-service/internal/service"
	pb "shoeshop/proto"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Инициализация MongoDB репозитория
	repo, err := repository.NewMongoRepository("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Failed to create MongoDB repository: %v", err)
	}

	// Инициализация NATS для событий
	natsClient, err := repository.NewNatsClient("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	// Подключение к Product Service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()
	productClient := pb.NewProductServiceClient(productConn)

	// Подключение к User Service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := pb.NewUserServiceClient(userConn)

	// Подключение к Email Service
	emailConn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to email service: %v", err)
	}
	defer emailConn.Close()
	emailClient := pb.NewEmailServiceClient(emailConn)

	// Инициализация сервиса
	svc := service.NewOrderService(repo, natsClient, productClient, userClient, emailClient)

	// Инициализация gRPC handler
	grpcHandler := handler.NewGRPCHandler(svc)

	// Создание gRPC сервера
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, grpcHandler)
	
	// Включаем reflection для отладки
	reflection.Register(server)

	// Запуск gRPC сервера
	port := ":50053"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Канал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting Order service on port %s", port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	sig := <-sigChan
	fmt.Printf("\nReceived signal %v, initiating graceful shutdown\n", sig)

	// Graceful shutdown
	server.GracefulStop()
	
	// Закрываем соединения
	if err := natsClient.Close(); err != nil {
		log.Printf("Error closing NATS connection: %v", err)
	}
	
	log.Println("Server stopped gracefully")
} 