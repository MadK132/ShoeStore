package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	pb "shoeshop/proto"
)

type EmailService interface {
	SendRegistrationConfirmation(ctx context.Context, user *pb.User) error
	SendOrderConfirmation(ctx context.Context, order *pb.Order) error
	SendOrderStatusUpdate(ctx context.Context, order *pb.Order) error
	SendPasswordReset(ctx context.Context, user *pb.User, resetToken string) error
}

type EmailConfig struct {
	SMTPHost  string
	SMTPPort  string
	Username  string
	Password  string
	FromEmail string
}

type emailService struct {
	config EmailConfig
}

func NewEmailService(config EmailConfig) EmailService {
	return &emailService{
		config: config,
	}
}

func (s *emailService) SendRegistrationConfirmation(ctx context.Context, user *pb.User) error {
	subject := "Welcome to ShoeShop!"
	templateData := map[string]interface{}{
		"Username": user.Username,
	}

	body := `
	<h2>Welcome to ShoeShop, {{.Username}}!</h2>
	<p>Thank you for registering with us. We're excited to have you as our customer.</p>
	<p>You can now browse our collection and make purchases.</p>
	<p>Best regards,<br>ShoeShop Team</p>
	`

	return s.sendEmail(user.Email, subject, body, templateData)
}

func (s *emailService) SendOrderConfirmation(ctx context.Context, order *pb.Order) error {
	subject := fmt.Sprintf("Order Confirmation #%s", order.Id)
	templateData := map[string]interface{}{
		"OrderID":     order.Id,
		"TotalAmount": order.TotalAmount,
		"Items":       order.Items,
	}

	body := `
	<h2>Order Confirmation</h2>
	<p>Thank you for your order #{{.OrderID}}!</p>
	<h3>Order Details:</h3>
	<ul>
	{{range .Items}}
		<li>{{.ProductId}} - Quantity: {{.Quantity}} - Price: ${{.Price}}</li>
	{{end}}
	</ul>
	<p><strong>Total Amount: ${{.TotalAmount}}</strong></p>
	<p>We'll notify you when your order ships.</p>
	<p>Best regards,<br>ShoeShop Team</p>
	`

	return s.sendEmail(order.ShippingAddress, subject, body, templateData)
}

func (s *emailService) SendOrderStatusUpdate(ctx context.Context, order *pb.Order) error {
	subject := fmt.Sprintf("Order Status Update #%s", order.Id)
	templateData := map[string]interface{}{
		"OrderID": order.Id,
		"Status":  order.Status,
	}

	body := `
	<h2>Order Status Update</h2>
	<p>Your order #{{.OrderID}} status has been updated to: {{.Status}}</p>
	<p>If you have any questions, please contact our support team.</p>
	<p>Best regards,<br>ShoeShop Team</p>
	`

	return s.sendEmail(order.ShippingAddress, subject, body, templateData)
}

func (s *emailService) SendPasswordReset(ctx context.Context, user *pb.User, resetToken string) error {
	subject := "Password Reset Request"
	templateData := map[string]interface{}{
		"Username":   user.Username,
		"ResetToken": resetToken,
	}

	body := `
	<h2>Password Reset Request</h2>
	<p>Hello {{.Username}},</p>
	<p>We received a request to reset your password. Use the following token to reset your password:</p>
	<p><strong>{{.ResetToken}}</strong></p>
	<p>If you didn't request this, please ignore this email.</p>
	<p>Best regards,<br>ShoeShop Team</p>
	`

	return s.sendEmail(user.Email, subject, body, templateData)
}

func (s *emailService) sendEmail(to, subject, bodyTemplate string, data map[string]interface{}) error {
	// Парсим шаблон
	tmpl, err := template.New("email").Parse(bodyTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Применяем данные к шаблону
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Формируем email сообщение
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n%s",
		to, s.config.FromEmail, subject, mime, body.String())

	// Отправляем email
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	if err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
} 