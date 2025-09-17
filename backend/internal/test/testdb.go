package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// TestDBConfig holds configuration for test database
type TestDBConfig struct {
	DatabaseURL string
	MaxConns    int32
}

// GetTestDBConfig returns test database configuration
func GetTestDBConfig() TestDBConfig {
	databaseURL := "postgres://postgres:postgres@localhost:5432/briefbot_test?sslmode=disable"
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		databaseURL = url
	}

	return TestDBConfig{
		DatabaseURL: databaseURL,
		MaxConns:    5,
	}
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	config := GetTestDBConfig()

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	require.NoError(t, err, "Failed to create database pool")

	// Test connection
	err = pool.Ping(context.Background())
	require.NoError(t, err, "Failed to ping database")

	// Clean up function
	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// CleanupTestDB cleans up test data
func CleanupTestDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	// Clean up test data but keep schema
	tables := []string{"items", "users"}

	for _, table := range tables {
		_, err := pool.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// WithTestTransaction executes a function within a database transaction that rolls back
func WithTestTransaction(t *testing.T, pool *pgxpool.Pool, fn func(tx pgx.Tx)) {
	t.Helper()

	tx, err := pool.Begin(context.Background())
	require.NoError(t, err, "Failed to begin transaction")

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
			panic(r)
		}
		tx.Rollback(context.Background())
	}()

	fn(tx)
}
