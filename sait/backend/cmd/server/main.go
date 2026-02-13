package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"bazarpo-backend/internal/handler"
	"bazarpo-backend/internal/model"
	"bazarpo-backend/internal/repo"
	"bazarpo-backend/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	env := model.LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(env.Database)

	repositories := repo.New(db)
	authService := service.NewAuthService(repositories, env)
	carService := service.NewCarService(repositories)
	partService := service.NewPartService(repositories)
	orderService := service.NewOrderService(repositories)
	adminPartService := service.NewAdminPartService(repositories)

	if err := authService.EnsureAdmin(context.Background()); err != nil {
		log.Println("ensureAdmin:", err)
	}

	h := handler.New(authService, carService, partService, orderService, adminPartService)

	addr := ":" + env.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      h.Routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Printf("bazarPO Go API on http://localhost%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
