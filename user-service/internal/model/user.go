package model

import (
	"time"
	pb "shoeshop/proto"
)

// User представляет доменную модель пользователя
type User struct {
	ID               string    `bson:"_id,omitempty"`
	Username         string    `bson:"username"`
	Email            string    `bson:"email"`
	PasswordHash     string    `bson:"password_hash"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
	FirstName        string    `bson:"first_name"`
	LastName         string    `bson:"last_name"`
	ShippingAddress  string    `bson:"shipping_address"`
	Phone            string    `bson:"phone"`
	RegistrationDate string    `bson:"registration_date"`
	OrderIDs         []string  `bson:"order_ids"`
	IsAdmin          bool      `bson:"is_admin"`
	Balance          float64   `bson:"balance"`
}

// ToProto конвертирует доменную модель в protobuf модель
func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:               u.ID,
		Username:         u.Username,
		Email:           u.Email,
		PasswordHash:     u.PasswordHash,
		CreatedAt:       u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       u.UpdatedAt.Format(time.RFC3339),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		ShippingAddress: u.ShippingAddress,
		Phone:           u.Phone,
		RegistrationDate: u.RegistrationDate,
		OrderIds:        u.OrderIDs,
		IsAdmin:         u.IsAdmin,
		Balance:         u.Balance,
	}
}

// FromProto конвертирует protobuf модель в доменную модель
func FromProto(pbUser *pb.User) (*User, error) {
	createdAt, err := time.Parse(time.RFC3339, pbUser.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	updatedAt, err := time.Parse(time.RFC3339, pbUser.UpdatedAt)
	if err != nil {
		updatedAt = time.Now()
	}

	return &User{
		ID:           pbUser.Id,
		Username:     pbUser.Username,
		Email:        pbUser.Email,
		PasswordHash: pbUser.PasswordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		FirstName:    pbUser.FirstName,
		LastName:     pbUser.LastName,
		ShippingAddress: pbUser.ShippingAddress,
		Phone:        pbUser.Phone,
		RegistrationDate: pbUser.RegistrationDate,
		OrderIDs:     pbUser.OrderIds,
		IsAdmin:      pbUser.IsAdmin,
		Balance:      pbUser.Balance,
	}, nil
} 