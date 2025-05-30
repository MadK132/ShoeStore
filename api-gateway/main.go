package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "shoeshop/proto"
)

type APIGateway struct {
	userClient    pb.UserServiceClient
	productClient pb.ProductServiceClient
	orderClient   pb.OrderServiceClient
	emailClient   pb.EmailServiceClient
}

func NewAPIGateway() (*APIGateway, error) {
	// Connect to User Service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to Product Service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to Order Service
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to Email Service
	emailConn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &APIGateway{
		userClient:    pb.NewUserServiceClient(userConn),
		productClient: pb.NewProductServiceClient(productConn),
		orderClient:   pb.NewOrderServiceClient(orderConn),
		emailClient:   pb.NewEmailServiceClient(emailConn),
	}, nil
}

func main() {
	gateway, err := NewAPIGateway()
	if err != nil {
		log.Fatalf("Failed to create API Gateway: %v", err)
	}

	r := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// API routes group
	api := r.Group("/api")
	{
		// User routes
		api.POST("/users/register", gateway.createUser)
		api.POST("/users/login", gateway.login)
		api.GET("/users/profile", gateway.getUser)
		api.PUT("/users/profile/update", gateway.updateUser)

		// Product routes
		api.POST("/products", gateway.createProduct)
		api.GET("/products/:id", gateway.getProduct)
		api.PUT("/products/:id", gateway.updateProduct)
		api.DELETE("/products/:id", gateway.deleteProduct)
		api.GET("/products", gateway.listProducts)

		// Order routes
		api.POST("/orders", gateway.createOrder)
		api.GET("/orders/:id", gateway.getOrder)
		api.GET("/orders/user/:userId", gateway.listUserOrders)
		api.PUT("/orders/:id/status", gateway.updateOrderStatus)
	}

	log.Fatal(r.Run(":8080"))
}

// User handlers
func (g *APIGateway) createUser(c *gin.Context) {
	var req struct {
		Username        string `json:"username"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		Phone          string `json:"phone"`
		ShippingAddress string `json:"shippingAddress"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert email to lowercase
	req.Email = strings.ToLower(req.Email)

	// Validate phone number format
	if !strings.HasPrefix(req.Phone, "+") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must start with '+'"})
		return
	}

	// Validate shipping address format
	if !strings.HasPrefix(req.ShippingAddress, "г.") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping address must start with 'г.'"})
		return
	}

	pbReq := &pb.RegisterRequest{
		Username:        req.Username,
		Email:          req.Email,
		Password:       req.Password,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		ShippingAddress: req.ShippingAddress,
	}

	resp, err := g.userClient.Register(context.Background(), pbReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) getUser(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}

	// Convert email to lowercase
	email = strings.ToLower(email)

	// Get user by email
	user, err := g.userClient.GetUserByEmail(context.Background(), &pb.GetUserByEmailRequest{Email: email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (g *APIGateway) updateUser(c *gin.Context) {
	// Получаем данные из запроса
	var req struct {
		User *pb.User `json:"user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.User == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User data is required"})
		return
	}

	// Получаем текущего пользователя
	currentUser, err := g.userClient.GetUserByEmail(context.Background(), &pb.GetUserByEmailRequest{
		Email: req.User.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user: " + err.Error()})
		return
	}

	// Сохраняем неизменяемые поля
	req.User.Id = currentUser.User.Id
	req.User.CreatedAt = currentUser.User.CreatedAt
	req.User.PasswordHash = currentUser.User.PasswordHash
	req.User.IsAdmin = currentUser.User.IsAdmin
	req.User.OrderIds = currentUser.User.OrderIds
	req.User.Balance = currentUser.User.Balance

	// Обновляем время изменения
	req.User.UpdatedAt = time.Now().Format(time.RFC3339)

	// Отправляем запрос на обновление
	resp, err := g.userClient.UpdateUser(context.Background(), &pb.UpdateUserRequest{
		User: req.User,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) login(c *gin.Context) {
	var req pb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert email to lowercase
	req.Email = strings.ToLower(req.Email)

	resp, err := g.userClient.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Product handlers
func (g *APIGateway) createProduct(c *gin.Context) {
	var product pb.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := g.productClient.CreateProduct(context.Background(), &pb.CreateProductRequest{Product: &product})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) getProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := g.productClient.GetProduct(context.Background(), &pb.GetProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) updateProduct(c *gin.Context) {
	var product pb.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := g.productClient.UpdateProduct(context.Background(), &pb.UpdateProductRequest{Product: &product})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) deleteProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := g.productClient.DeleteProduct(context.Background(), &pb.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) listProducts(c *gin.Context) {
	category := c.Query("category")
	brand := c.Query("brand")

	resp, err := g.productClient.ListProducts(context.Background(), &pb.ListProductsRequest{
		Category: category,
		Brand:    brand,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Order handlers
func (g *APIGateway) createOrder(c *gin.Context) {
	var order pb.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := g.orderClient.CreateOrder(context.Background(), &pb.CreateOrderRequest{Order: &order})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send order confirmation email
	_, err = g.userClient.GetUser(context.Background(), &pb.GetUserRequest{Id: order.UserId})
	if err != nil {
		log.Printf("Failed to get user for order confirmation: %v", err)
	} else {
		err = g.sendOrderEmail(resp.Order)
		if err != nil {
			log.Printf("Failed to send order confirmation email: %v", err)
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) getOrder(c *gin.Context) {
	id := c.Param("id")
	resp, err := g.orderClient.GetOrder(context.Background(), &pb.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) listUserOrders(c *gin.Context) {
	userID := c.Param("userId")
	
	resp, err := g.orderClient.ListOrders(context.Background(), &pb.ListOrdersRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *APIGateway) updateOrderStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID := c.Param("id")
	_, err := g.orderClient.UpdateOrderStatus(context.Background(), &pb.UpdateOrderStatusRequest{
		Id:     orderID,
		Status: req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func (g *APIGateway) sendOrderEmail(order *pb.Order) error {
	_, err := g.emailClient.SendOrderConfirmation(context.Background(), &pb.OrderEmailRequest{
		Order: order,
	})
	return err
} 