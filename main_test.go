package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func waitForServer(url string, timeout time.Duration) error {
	client := http.Client{
		Timeout: timeout,
	}
	for start := time.Now(); time.Since(start) < timeout; {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode != http.StatusNotFound {
			resp.Body.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server did not start within timeout")
}

func TestMainFunction(t *testing.T) {

	dbURL := "postgres://wallet:wallet@localhost:5432/wallet?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("Пропущен тест: PostgreSQL не запущен или миграции не применены: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(context.Background(), "DELETE FROM wallets")
	assert.NoError(t, err)

	go main()

	serverURL := "http://localhost:8080/api/v1/wallets/00000000-0000-0000-0000-000000000000"
	err = waitForServer(serverURL, 10*time.Second)
	assert.NoError(t, err)

	t.Run("POST /api/v1/wallet", func(t *testing.T) {

		walletID := "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8"
		payload := fmt.Sprintf(`{
            "walletId": "%s",
            "operationType": "DEPOSIT",
            "amount": 100.0
        }`, walletID)

		resp, err := http.Post("http://localhost:8080/api/v1/wallet", "application/json", strings.NewReader(payload))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"balance":100.0`)
	})

	t.Run("GET /api/v1/wallets/{id}", func(t *testing.T) {
		walletID := "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8"
		url := fmt.Sprintf("http://localhost:8080/api/v1/wallets/%s", walletID)

		resp, err := http.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"balance":100.0`)
	})
}
