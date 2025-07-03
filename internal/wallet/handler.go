package wallet

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletRequest struct {
	WalletID      uuid.UUID       `json:"walletId"`
	OperationType string          `json:"operationType"`
	Amount        decimal.Decimal `json:"amount"`
}

// UpdateBalanceHandler обновляет баланс кошелька
// @Summary Обновить баланс
// @Tags Wallet
// @Accept json
// @Produce json
// @Param wallet body WalletRequest true "Запрос на обновление баланса"
// @Success 200 {object} map[string]string "новый баланс"
// @Failure 400 {object} map[string]string "ошибка запроса"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /api/v1/wallet [post]

func UpdateBalanceHandler(service *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req WalletRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid operation type"})
		}

		if req.Amount.LessThan(decimal.NewFromFloat(0.0)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Amount must be positive"})
		}

		if err := service.UpdateBalance(context.Background(), req.WalletID, req.OperationType, req.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		newBalance, err := service.GetBalance(context.Background(), req.WalletID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"balance": newBalance.String()})
	}
}

// GetBalanceHandler возвращает баланс кошелька
// @Summary Получить баланс
// @Tags Wallet
// @Produce json
// @Param id path string true "UUID кошелька"
// @Success 200 {object} map[string]string "текущий баланс"
// @Failure 400 {object} map[string]string "ошибка запроса"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /api/v1/wallets/{id} [get]

func GetBalanceHandler(service *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid wallet ID"})
		}

		balance, err := service.GetBalance(context.Background(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"balance": balance.String()})
	}
}
