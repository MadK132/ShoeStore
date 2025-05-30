package service

import (
	"context"
	"time"
	"golang.org/x/crypto/bcrypt"
	"strings"

	"shoeshop/user-service/internal/model"
	"shoeshop/user-service/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]*model.User, error)
	Register(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// Добавляем дату регистрации
	user.RegistrationDate = time.Now().Format(time.RFC3339)
	
	// Создаем пользователя в БД
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// Получаем текущего пользователя из БД
	currentUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	// Сохраняем неизменяемые поля
	user.ID = currentUser.ID
	user.CreatedAt = currentUser.CreatedAt
	user.PasswordHash = currentUser.PasswordHash
	user.IsAdmin = currentUser.IsAdmin
	user.RegistrationDate = currentUser.RegistrationDate
	user.OrderIDs = currentUser.OrderIDs

	// Обновляем время изменения
	user.UpdatedAt = time.Now()

	// Проверяем и форматируем адрес
	if user.ShippingAddress != "" && !strings.HasPrefix(user.ShippingAddress, "г.") {
		user.ShippingAddress = "г." + user.ShippingAddress
	}

	// Проверяем и форматируем телефон
	if user.Phone != "" && !strings.HasPrefix(user.Phone, "+") {
		user.Phone = "+" + user.Phone
	}

	// Обновляем пользователя в БД
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	// Установка дополнительных полей
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.RegistrationDate = time.Now().Format(time.RFC3339)
	user.IsAdmin = false
	user.Balance = 845 // Начальный баланс
	user.OrderIDs = []string{} // Инициализация пустым слайсом

	// Создание пользователя
	return s.repo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (*model.User, error) {
	// Получение пользователя по email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email)
} 