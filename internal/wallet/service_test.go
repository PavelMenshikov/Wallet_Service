package wallet

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	getBalanceFunc    func(ctx context.Context, id string) (decimal.Decimal, error)
	updateBalanceFunc func(ctx context.Context, id string, operation string, amount decimal.Decimal) (decimal.Decimal, error)
}

func (m *mockRepository) GetBalance(ctx context.Context, id string) (decimal.Decimal, error) {
	return m.getBalanceFunc(ctx, id)
}

func (m *mockRepository) UpdateBalance(ctx context.Context, id string, operation string, amount decimal.Decimal) (decimal.Decimal, error) {
	return m.updateBalanceFunc(ctx, id, operation, amount)
}

func TestService_GetBalance(t *testing.T) {
	mockRepo := &mockRepository{
		getBalanceFunc: func(ctx context.Context, id string) (decimal.Decimal, error) {
			return decimal.NewFromFloat(100.0), nil
		},
	}

	service := &Service{repo: mockRepo}
	id := uuid.New()
	balance, err := service.GetBalance(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(100.0), balance)
}

func TestService_UpdateBalance(t *testing.T) {
	mockRepo := &mockRepository{
		updateBalanceFunc: func(ctx context.Context, id string, operation string, amount decimal.Decimal) (decimal.Decimal, error) {
			assert.Equal(t, "DEPOSIT", operation)
			assert.Equal(t, decimal.NewFromFloat(50.0), amount)
			return decimal.NewFromFloat(150.0), nil
		},
	}

	service := &Service{repo: mockRepo}
	id := uuid.New()
	err := service.UpdateBalance(context.Background(), id, "DEPOSIT", decimal.NewFromFloat(50.0))

	assert.NoError(t, err)
}
