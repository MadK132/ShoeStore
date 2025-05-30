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

	"shoeshop/email-service/internal/handler"
	"shoeshop/email-service/internal/service"
	pb "shoeshop/proto"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Проверяем наличие необходимых переменных окружения
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("FROM_EMAIL")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || fromEmail == "" {
		log.Fatal("Missing required environment variables")
	}

	// Инициализация email сервиса
	emailSvc := service.NewEmailService(service.EmailConfig{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: smtpUser,
		Password: smtpPass,
		FromEmail: fromEmail,
	})

	// Инициализация gRPC handler
	grpcHandler := handler.NewGRPCHandler(emailSvc)

	// Создание gRPC сервера
	server := grpc.NewServer()
	pb.RegisterEmailServiceServer(server, grpcHandler)
	
	// Включаем reflection для отладки
	reflection.Register(server)

	// Запуск gRPC сервера
	port := ":50054"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Канал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting Email service on port %s", port)
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