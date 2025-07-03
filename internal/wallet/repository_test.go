package wallet

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), "postgres://wallet:wallet@localhost:5432/wallet?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	_, err = db.Exec(context.Background(), "DELETE FROM wallets")
	if err != nil {
		t.Fatalf("Failed to clean DB: %v", err)
	}

	return db
}

func TestRepository_GetBalance(t *testing.T) {
	db := setupDB(t)
	repo := NewRepository(db)

	_, err := db.Exec(context.Background(), "INSERT INTO wallets (id, balance) VALUES ($1, $2)", "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8", decimal.NewFromFloat(100.0))
	assert.NoError(t, err)

	balance, err := repo.GetBalance(context.Background(), "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8")
	assert.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(100.0), balance)
}

func TestRepository_UpdateBalance_DEPOSIT(t *testing.T) {
	db := setupDB(t)
	repo := NewRepository(db)

	tx, _ := db.Begin(context.Background())
	_, err := tx.Exec(context.Background(), "INSERT INTO wallets (id, balance) VALUES ($1, $2)", "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8", decimal.NewFromFloat(100.0))
	assert.NoError(t, err)
	tx.Commit(context.Background())

	err = repo.UpdateBalance(context.Background(), "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8", "DEPOSIT", decimal.NewFromFloat(50.0))
	assert.NoError(t, err)

	balance, err := repo.GetBalance(context.Background(), "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8")
	assert.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(150.0), balance)
}

func TestRepository_UpdateBalance_WITHDRAW(t *testing.T) {
	db := setupDB(t)
	repo := NewRepository(db)

	tx, _ := db.Begin(context.Background())
	_, err := tx.Exec(context.Background(), "INSERT INTO wallets (id, balance) VALUES ($1, $2)", "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8", decimal.NewFromFloat(100.0))
	assert.NoError(t, err)
	tx.Commit(context.Background())

	err = repo.UpdateBalance(context.Background(), "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8", "WITHDRAW", decimal.NewFromFloat(30.0))
	assert.NoError(t, err)

	balance, err := repo.GetBalance(context.Background(), "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8")
	assert.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(70.0), balance)
}
