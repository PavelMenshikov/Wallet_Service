package main

import (
	"context"
	"log"
	"os"

	"wallet-service/internal/wallet"

	_ "wallet-service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type repoAdapter struct {
	*wallet.Repository
}

func (r *repoAdapter) UpdateBalance(ctx context.Context, id string, operation string, amount decimal.Decimal) (decimal.Decimal, error) {
	err := r.Repository.UpdateBalance(ctx, id, operation, amount)
	if err != nil {
		return decimal.Zero, err
	}

	return r.Repository.GetBalance(ctx, id)
}

// @title Wallet Service API
// @version 1.0
// @description API для управления балансом кошельков
// @host localhost:8080
// @BasePath /

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	log.Printf("Connecting to DB: %s", dbURL)

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	repo := &repoAdapter{wallet.NewRepository(pool)}
	service := wallet.NewService(repo)

	app := fiber.New()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Post("/api/v1/wallet", wallet.UpdateBalanceHandler(service))
	app.Get("/api/v1/wallets/:id", wallet.GetBalanceHandler(service))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Wallet API is running. Go to /swagger/index.html for docs.")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
