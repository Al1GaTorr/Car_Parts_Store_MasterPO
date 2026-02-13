package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Env struct {
	MongoURI string
	Database string
}

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func loadEnv() Env {
	return Env{
		MongoURI: getEnv("MONGO_URI", "mongodb://127.0.0.1:27017/bazarPO"),
		Database: getEnv("MONGO_DB", "bazarPO"),
	}
}

type Car struct {
	VIN       string `json:"vin" bson:"vin"`
	Make      string `json:"make" bson:"make"`
	Model     string `json:"model" bson:"model"`
	Year      int    `json:"year" bson:"year"`
	Engine    string `json:"engine,omitempty" bson:"engine,omitempty"`
	Fuel      string `json:"fuel,omitempty" bson:"fuel,omitempty"`
	MileageKM int    `json:"mileage_km,omitempty" bson:"mileage_km,omitempty"`
}

type VehicleFit struct {
	Make     string `bson:"make" json:"make"`
	Model    string `bson:"model" json:"model"`
	YearFrom int    `bson:"year_from" json:"year_from"`
	YearTo   int    `bson:"year_to" json:"year_to"`
}

type Compatibility struct {
	Type     string       `bson:"type" json:"type"`
	VINs     []string     `bson:"vins,omitempty" json:"vins,omitempty"`
	Vehicles []VehicleFit `bson:"vehicles,omitempty" json:"vehicles,omitempty"`
}

type Part struct {
	SKU       string        `bson:"sku" json:"sku"`
	Name      string        `bson:"name" json:"name"`
	Category  string        `bson:"category" json:"category"`
	Type      string        `bson:"type,omitempty" json:"type,omitempty"`
	Brand     string        `bson:"brand,omitempty" json:"brand,omitempty"`
	PriceKZT  int           `bson:"price_kzt" json:"price_kzt"`
	Currency  string        `bson:"currency" json:"currency"`
	StockQty  int           `bson:"stock_qty" json:"stock_qty"`
	IsVisible bool          `bson:"is_visible" json:"is_visible"`
	Images    []string      `bson:"images,omitempty" json:"images,omitempty"`
	Compat    Compatibility `bson:"compatibility" json:"compatibility"`
	Recommend []string      `bson:"recommend_for_issue_codes,omitempty" json:"recommend_for_issue_codes,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	env := loadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(env.Database)

	carsCol := db.Collection("cars")
	partsCol := db.Collection("parts")
	ordersCol := db.Collection("orders")
	usersCol := db.Collection("users")

	// Clear collections (keep users)
	_, _ = carsCol.DeleteMany(ctx, bson.M{})
	_, _ = partsCol.DeleteMany(ctx, bson.M{})
	_, _ = ordersCol.DeleteMany(ctx, bson.M{})

	// Load cars
	b, err := os.ReadFile("seed_data/cars.json")
	if err != nil {
		log.Fatal("read cars.json:", err)
	}
	var cars []Car
	if err := json.Unmarshal(b, &cars); err != nil {
		log.Fatal("parse cars.json:", err)
	}
	if len(cars) == 0 {
		log.Fatal("cars.json empty")
	}
	carDocs := make([]any, 0, len(cars))
	for _, c := range cars {
		c.VIN = strings.ToUpper(strings.TrimSpace(c.VIN))
		carDocs = append(carDocs, c)
	}
	_, err = carsCol.InsertMany(ctx, carDocs)
	if err != nil {
		log.Fatal("insert cars:", err)
	}

	// Ensure admin exists if env has it (server will also ensure)
	_ = usersCol

	// Categories (must match frontend PartCategory values)
	categories := []string{
		"Автохимия", "Автоаксессуары", "Масла и жидкости", "Инструменты", "Диски и шины", "Автолампы",
		"Накидки на сидения", "Накидки с обогревом", "Автокресла и бустеры", "Предпусковые обогреватели",
		"Пуско зарядные устройства", "Провода пусковые", "Канистры", "Домкраты", "Компрессоры",
		"Огнетушители", "Наборы автомобилиста", "Аптечки", "Тросы буксировочные", "Знаки аварийной остановки",
		"Аккумуляторы", "Щетки дворников",
	}

	issueMap := map[string][]string{
		"brake_pads_worn": {"Тормозные колодки", "Комплект колодок"},
		"oil_change":      {"Моторное масло", "Масляный фильтр"},
		"battery_dead":    {"Аккумулятор"},
		"wipers_bad":      {"Щетки дворников"},
		"headlight_out":   {"Лампа"},
	}

	imagePool := []string{
		"https://images.unsplash.com/photo-1515923168762-9f9f2b67b7ea?auto=format&fit=crop&q=80&w=900",
		"https://images.unsplash.com/photo-1603386329225-868f9b1ee6ae?auto=format&fit=crop&q=80&w=900",
		"https://images.unsplash.com/photo-1517677208171-0bc6725a3e60?auto=format&fit=crop&q=80&w=900",
		"https://images.unsplash.com/photo-1558981806-ec527fa84c39?auto=format&fit=crop&q=80&w=900",
		"https://images.unsplash.com/photo-1619642751034-765dfdf7c58e?auto=format&fit=crop&q=80&w=900",
		"https://images.unsplash.com/photo-1530124566582-a618bc2615dc?auto=format&fit=crop&q=80&w=900",
	}

	// Build make/model list from cars
	type MM struct{ Make, Model string }
	uniq := map[string]MM{}
	for _, c := range cars {
		key := strings.ToLower(c.Make + "|" + c.Model)
		uniq[key] = MM{Make: c.Make, Model: c.Model}
	}
	mms := make([]MM, 0, len(uniq))
	for _, v := range uniq {
		mms = append(mms, v)
	}

	makeSKU := func(prefix string) string {
		const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		b := make([]byte, 8)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return fmt.Sprintf("%s-%s", prefix, string(b))
	}

	parts := make([]any, 0, len(categories)*6)
	for _, cat := range categories {
		for i := 0; i < 6; i++ {
			mm := mms[rand.Intn(len(mms))]
			yearFrom := 2008 + rand.Intn(8)
			yearTo := yearFrom + 6 + rand.Intn(6)

			name := fmt.Sprintf("%s #%d", cat, i+1)
			reco := []string{}
			typ := "generic"
			// assign some issue codes for realism
			if cat == "Аккумуляторы" && i%2 == 0 {
				reco = []string{"battery_dead"}
				name = "Аккумулятор 12V"
				typ = "battery"
			}
			if cat == "Щетки дворников" && i%2 == 0 {
				reco = []string{"wipers_bad"}
				name = "Щетки дворников (комплект)"
				typ = "wiper_blades"
			}
			if cat == "Масла и жидкости" && i%2 == 0 {
				reco = []string{"oil_change"}
				name = "Моторное масло 5W-30 (4л)"
				typ = "engine_oil"
			}
			if cat == "Автоаксессуары" && i%3 == 0 {
				reco = []string{"brake_pads_worn"}
				name = "Тормозные колодки (комплект)"
				typ = "brake_pads"
			}
			if cat == "Автолампы" && i%3 == 0 {
				reco = []string{"headlight_out"}
				name = "Лампа головного света"
				typ = "headlight_bulb"
			}
			// ensure reco list strings correspond to issue map if present
			_ = issueMap

			p := Part{
				SKU:       makeSKU("SKU"),
				Name:      name,
				Category:  cat,
				Type:      typ,
				Brand:     []string{"SKF", "Bosch", "Mann", "Shell", "Castrol", "Philips"}[rand.Intn(6)],
				PriceKZT:  1500 + rand.Intn(85000),
				Currency:  "KZT",
				StockQty:  1 + rand.Intn(30),
				IsVisible: true,
				Images:    []string{imagePool[rand.Intn(len(imagePool))]},
				Compat: Compatibility{
					Type: "vehicle",
					Vehicles: []VehicleFit{
						{Make: mm.Make, Model: mm.Model, YearFrom: yearFrom, YearTo: yearTo},
					},
				},
				Recommend: reco,
				CreatedAt: time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour),
			}

			parts = append(parts, p)
		}
	}

	_, err = partsCol.InsertMany(ctx, parts)
	if err != nil {
		log.Fatal("insert parts:", err)
	}

	fmt.Printf("Seed OK: cars=%d, parts=%d\n", len(cars), len(parts))
}
