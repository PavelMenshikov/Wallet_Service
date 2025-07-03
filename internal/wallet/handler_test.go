package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBalanceHandler(t *testing.T) {
	id := uuid.New()
	app := fiber.New()

	mockRepo := &mockRepository{
		updateBalanceFunc: func(ctx context.Context, walletID string, operation string, amount decimal.Decimal) (decimal.Decimal, error) {
			return decimal.NewFromFloat(200.0), nil
		},
		getBalanceFunc: func(ctx context.Context, id string) (decimal.Decimal, error) {
			return decimal.NewFromFloat(200.0), nil
		},
	}

	service := &Service{repo: mockRepo}
	app.Post("/api/v1/wallet", UpdateBalanceHandler(service))

	reqBody := WalletRequest{
		WalletID:      id,
		OperationType: "DEPOSIT",
		Amount:        decimal.NewFromFloat(100.0),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetBalanceHandler(t *testing.T) {
	id := uuid.New()
	app := fiber.New()

	mockRepo := &mockRepository{
		getBalanceFunc: func(ctx context.Context, id string) (decimal.Decimal, error) {
			return decimal.NewFromFloat(123.45), nil
		},
	}

	service := &Service{repo: mockRepo}
	app.Get("/api/v1/wallets/:id", GetBalanceHandler(service))

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/wallets/%s", id.String()), nil)
	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
