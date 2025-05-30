package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShippingInfo struct {
	Dimensions struct {
		Length float64 `bson:"length"`
		Width  float64 `bson:"width"`
		Height float64 `bson:"height"`
	} `bson:"dimensions"`
	Weight          float64 `bson:"weight"`
	ShippingMethods []struct {
		Method      string  `bson:"method"`
		Price       float64 `bson:"price"`
		DeliveryTime string `bson:"delivery_time"`
	} `bson:"shipping_methods"`
	FreeShipping     bool `bson:"free_shipping"`
	International    bool `bson:"international"`
	TrackingAvailable bool `bson:"tracking_available"`
}

type Product struct {
	ID            string       `bson:"_id"`
	NumericID     int         `bson:"id"`
	Name          string      `bson:"name"`
	Brand         string      `bson:"brand"`
	Sizes         []int       `bson:"sizes"`
	Price         float64     `bson:"price"`
	Material      string      `bson:"material"`
	Colors        []string    `bson:"colors"`
	ReleaseDate   string      `bson:"release_date"`
	Discount      int         `bson:"discount"`
	StockQuantity int         `bson:"stock_quantity"`
	Rating        float64     `bson:"rating"`
	SKU           string      `bson:"sku"`
	Weight        string      `bson:"weight"`
	Features      []string    `bson:"features"`
	Category      string      `bson:"category"`
	ShippingInfo  ShippingInfo `bson:"shipping_info"`
}

func getNextSequence(client *mongo.Client, sequenceName string) (int, error) {
	collection := client.Database("shoeshop").Collection("counters")
	
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		ID  string `bson:"_id"`
		Seq int    `bson:"seq"`
	}

	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}

func main() {
	// Подключение к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Получение коллекции
	collection := client.Database("shoeshop").Collection("products")

	// Очистка существующих данных
	collection.DeleteMany(context.Background(), bson.M{})
	client.Database("shoeshop").Collection("counters").DeleteMany(context.Background(), bson.M{})

	// Данные для генерации продуктов
	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "New Balance"}
	categories := []string{"Running", "Basketball", "Training", "Lifestyle"}
	materials := []string{"Leather", "Synthetic", "Canvas", "Mesh"}
	colors := [][]string{
		{"Black", "White"},
		{"Red", "Blue"},
		{"Grey", "White"},
		{"Navy", "Red"},
		{"Green", "Black"},
	}
	features := [][]string{
		{"Breathable", "Lightweight", "Durable"},
		{"Water-resistant", "Shock-absorbing", "Anti-slip"},
		{"Cushioned", "Flexible", "Supportive"},
		{"Quick-drying", "Ventilated", "Stable"},
	}

	// Генерация 50 продуктов
	for i := 1; i <= 50; i++ {
		// Получаем следующий ID
		numericID, err := getNextSequence(client, "productid")
		if err != nil {
			log.Printf("Error getting next sequence: %v", err)
			continue
		}

		brand := brands[rand.Intn(len(brands))]
		category := categories[rand.Intn(len(categories))]
		material := materials[rand.Intn(len(materials))]
		productColors := colors[rand.Intn(len(colors))]
		productFeatures := features[rand.Intn(len(features))]
		
		// Генерация случайных размеров (от 36 до 46)
		var sizes []int
		minSize := 36 + rand.Intn(3)
		maxSize := 44 + rand.Intn(3)
		for size := minSize; size <= maxSize; size++ {
			sizes = append(sizes, size)
		}

		// Генерация информации о доставке
		shippingInfo := ShippingInfo{
			Dimensions: struct {
				Length float64 `bson:"length"`
				Width  float64 `bson:"width"`
				Height float64 `bson:"height"`
			}{
				Length: 30 + rand.Float64()*10,
				Width:  15 + rand.Float64()*5,
				Height: 10 + rand.Float64()*5,
			},
			Weight:            0.5 + rand.Float64()*1.0,
			FreeShipping:      rand.Float64() > 0.5,
			International:     rand.Float64() > 0.3,
			TrackingAvailable: rand.Float64() > 0.2,
		}

		// Добавляем методы доставки
		shippingInfo.ShippingMethods = []struct {
			Method       string  `bson:"method"`
			Price       float64 `bson:"price"`
			DeliveryTime string `bson:"delivery_time"`
		}{
			{Method: "Standard", Price: 5.99, DeliveryTime: "3-5 days"},
			{Method: "Express", Price: 15.99, DeliveryTime: "1-2 days"},
			{Method: "Next Day", Price: 25.99, DeliveryTime: "Next day"},
		}

		product := Product{
			ID:            strconv.Itoa(numericID),
			NumericID:     numericID,
			Name:          fmt.Sprintf("%s Air Max %d", brand, i),
			Brand:         brand,
			Category:      category,
			Price:         50.0 + float64(rand.Intn(101)),
			Sizes:         sizes,
			Material:      material,
			Colors:        productColors,
			ReleaseDate:   time.Now().Format("2006-01-02"),
			Discount:      rand.Intn(31),
			StockQuantity: 10 + rand.Intn(91),
			Rating:        3.5 + rand.Float64()*1.5,
			SKU:           fmt.Sprintf("SKU%d", 1000+i),
			Weight:        fmt.Sprintf("%dg", 200+rand.Intn(300)),
			Features:      productFeatures,
			ShippingInfo:  shippingInfo,
		}

		_, err = collection.InsertOne(context.Background(), product)
		if err != nil {
			log.Printf("Error inserting product %s: %v", product.Name, err)
		} else {
			log.Printf("Successfully inserted product: %s (ID: %d)", product.Name, product.NumericID)
		}
	}

	log.Printf("Successfully inserted 50 products")
} 