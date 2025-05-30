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

	"shoeshop/product-service/internal/handler"
	"shoeshop/product-service/internal/repository"
	"shoeshop/product-service/internal/service"
	pb "shoeshop/proto"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Инициализация MongoDB репозитория
	repo, err := repository.NewMongoRepository("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Failed to create MongoDB repository: %v", err)
	}

	// Инициализация Redis кэша
	cache, err := repository.NewRedisCache("localhost:6379")
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		// Продолжаем работу без кэша
	}

	// Инициализация NATS для событий
	natsClient, err := repository.NewNatsClient("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	// Инициализация сервиса
	svc := service.NewProductService(repo, cache, natsClient)

	// Инициализация gRPC handler
	grpcHandler := handler.NewGRPCHandler(svc)

	// Создание gRPC сервера
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, grpcHandler)
	
	// Включаем reflection для отладки
	reflection.Register(server)

	// Запуск gRPC сервера
	port := ":50052"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Канал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting Product service on port %s", port)
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
	if err := cache.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}
	if err := natsClient.Close(); err != nil {
		log.Printf("Error closing NATS connection: %v", err)
	}
	
	log.Println("Server stopped gracefully")
} 