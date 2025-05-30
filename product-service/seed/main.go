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

func main() {
	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Get database and collection
	db := client.Database("shoeshop")
	collection := db.Collection("products")

	// Drop existing collection
	if err := collection.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	// Product data variations
	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "New Balance", "Under Armour", "ASICS", "Converse", "Vans", "Skechers"}
	models := []string{"Air Max", "Superstar", "RS-X", "Classic", "Fresh Foam", "Charged", "Gel-Nimbus", "Chuck Taylor", "Old Skool", "Go Walk"}
	materials := []string{"Leather", "Mesh", "Synthetic", "Canvas", "Knit", "Suede", "Textile", "Nylon", "Cotton", "Polyester"}
	categories := []string{"Running", "Casual", "Sport", "Training", "Basketball", "Tennis", "Walking", "Skateboarding", "Lifestyle", "Gym"}
	colors := [][]string{
		{"White", "Black"},
		{"Black", "Red"},
		{"Navy", "White"},
		{"Grey", "Blue"},
		{"White", "Green"},
		{"Black", "Gold"},
		{"Red", "White"},
		{"Blue", "Grey"},
		{"Green", "Black"},
		{"Purple", "White"},
	}
	features := [][]string{
		{"Comfortable", "Durable", "Stylish"},
		{"Classic", "Comfortable", "Versatile"},
		{"Modern", "Lightweight", "Sporty"},
		{"Breathable", "Flexible", "Supportive"},
		{"Cushioned", "Stable", "Responsive"},
		{"Waterproof", "Grippy", "Protective"},
		{"Eco-friendly", "Recyclable", "Sustainable"},
		{"Slip-resistant", "Anti-odor", "Quick-dry"},
		{"Shock-absorbing", "Ventilated", "Balanced"},
		{"Memory foam", "Arch support", "Heel cushioning"},
	}

	// Shipping variations
	shippingMethods := []string{"Standard", "Express", "Next Day"}
	shippingPrices := []float64{5.99, 12.99, 19.99}
	deliveryTimes := []string{"3-5 business days", "1-2 business days", "Next business day"}

	// Create 50 products
	var products []interface{}
	for i := 1; i <= 50; i++ {
		brandIndex := rand.Intn(len(brands))
		modelIndex := rand.Intn(len(models))
		name := fmt.Sprintf("%s %s %d", brands[brandIndex], models[modelIndex], i)
		
		// Create shipping info with multiple delivery options
		var shippingOptions []bson.D
		for j := 0; j < len(shippingMethods); j++ {
			shippingOption := bson.D{
				{Key: "method", Value: shippingMethods[j]},
				{Key: "price", Value: shippingPrices[j]},
				{Key: "estimated_delivery", Value: deliveryTimes[j]},
				{Key: "tracking_available", Value: true},
				{Key: "international_shipping", Value: j > 0}, // Only Express and Next Day for international
			}
			shippingOptions = append(shippingOptions, shippingOption)
		}

		product := bson.D{
			{Key: "_id", Value: strconv.Itoa(i)},
			{Key: "name", Value: name},
			{Key: "brand", Value: brands[brandIndex]},
			{Key: "size", Value: []int{36, 37, 38, 39, 40, 41, 42, 43, 44}},
			{Key: "price", Value: 50 + rand.Float64()*100},
			{Key: "material", Value: materials[rand.Intn(len(materials))]},
			{Key: "color", Value: colors[rand.Intn(len(colors))]},
			{Key: "release_date", Value: time.Now().AddDate(0, rand.Intn(12), rand.Intn(28)).Format("2006-01-02")},
			{Key: "discount", Value: float64(rand.Intn(30)) / 100},
			{Key: "stock_quantity", Value: 10 + rand.Intn(91)},
			{Key: "rating", Value: 3.5 + rand.Float64()*1.5},
			{Key: "sku", Value: fmt.Sprintf("SKU%04d", i)},
			{Key: "weight", Value: fmt.Sprintf("%dg", 150+rand.Intn(101))},
			{Key: "features", Value: features[rand.Intn(len(features))]},
			{Key: "category", Value: categories[rand.Intn(len(categories))]},
			{Key: "shipping_info", Value: bson.D{
				{Key: "dimensions", Value: bson.D{
					{Key: "length", Value: "30cm"},
					{Key: "width", Value: "20cm"},
					{Key: "height", Value: "15cm"},
				}},
				{Key: "shipping_options", Value: shippingOptions},
				{Key: "free_shipping_eligible", Value: rand.Float64() < 0.3}, // 30% chance for free shipping
				{Key: "shipping_restrictions", Value: []string{}},
				{Key: "handling_time", Value: "1 business day"},
			}},
		}
		products = append(products, product)
	}

	// Insert products
	for _, product := range products {
		_, err := collection.InsertOne(ctx, product)
		if err != nil {
			log.Printf("Error inserting product: %v", err)
			continue
		}
		log.Printf("Successfully inserted product: %v", product.(bson.D).Map()["name"])
	}

	log.Println("Seed completed successfully")
} 