package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RepositoryInterface interface {
	GetBalance(ctx context.Context, id string) (decimal.Decimal, error)
	UpdateBalance(ctx context.Context, id string, operation string, amount decimal.Decimal) (decimal.Decimal, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBalance(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error) {
	return s.repo.GetBalance(ctx, walletID.String())
}

func (s *Service) UpdateBalance(ctx context.Context, walletID uuid.UUID, operation string, amount decimal.Decimal) error {
	_, err := s.repo.UpdateBalance(ctx, walletID.String(), operation, amount)
	return err
}
