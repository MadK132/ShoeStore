package model

import (
	"fmt"
	"strconv"
	"time"
	pb "shoeshop/proto"
)

type Product struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Price       float64   `bson:"price"`
	Category    string    `bson:"category"`
	Brand       string    `bson:"brand"`
	Sizes       []int     `bson:"sizes"`
	Colors      []string  `bson:"colors"`
	Images      []string  `bson:"images"`
	Stock       int32     `bson:"stock"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// ToProto конвертирует доменную модель в protobuf модель
func (p *Product) ToProto() *pb.Product {
	// Конвертируем []int в []string для размеров
	sizes := make([]string, len(p.Sizes))
	for i, size := range p.Sizes {
		sizes[i] = fmt.Sprintf("%d", size)
	}

	return &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Brand:       p.Brand,
		Sizes:       sizes,
		Colors:      p.Colors,
		Images:      p.Images,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

// FromProto конвертирует protobuf модель в доменную модель
func FromProto(pbProduct *pb.Product) (*Product, error) {
	createdAt, err := time.Parse(time.RFC3339, pbProduct.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	updatedAt, err := time.Parse(time.RFC3339, pbProduct.UpdatedAt)
	if err != nil {
		updatedAt = time.Now()
	}

	// Конвертируем []string в []int для размеров
	sizes := make([]int, len(pbProduct.Sizes))
	for i, size := range pbProduct.Sizes {
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			return nil, fmt.Errorf("invalid size format: %v", err)
		}
		sizes[i] = sizeInt
	}

	return &Product{
		ID:          pbProduct.Id,
		Name:        pbProduct.Name,
		Description: pbProduct.Description,
		Price:       pbProduct.Price,
		Category:    pbProduct.Category,
		Brand:       pbProduct.Brand,
		Sizes:       sizes,
		Colors:      pbProduct.Colors,
		Images:      pbProduct.Images,
		Stock:       pbProduct.Stock,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
} 