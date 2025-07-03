package wallet

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBalance(ctx context.Context, id string) (decimal.Decimal, error) {
	var balance decimal.Decimal
	err := r.db.QueryRow(ctx, "SELECT balance FROM wallets WHERE id = $1", id).Scan(&balance)
	return balance, err
}

func (r *Repository) UpdateBalance(ctx context.Context, id string, operation string, amount decimal.Decimal) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	var query string
	if operation == "DEPOSIT" {
		query = "UPDATE wallets SET balance = balance + $1 WHERE id = $2"
	} else if operation == "WITHDRAW" {
		query = "UPDATE wallets SET balance = balance - $1 WHERE id = $2"
	} else {
		return fmt.Errorf("invalid operation type: %s", operation)
	}

	_, err = tx.Exec(ctx, query, amount, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
