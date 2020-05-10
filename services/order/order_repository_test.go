package main

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool
var conn *pgxpool.Conn
var repository *orderRepository

const testDbURL = "postgres://postgres:xchange@localhost/order_test?sslmode=disable&pool_max_conns=10"

func cleanup() {
	conn.Exec(context.Background(), `DELETE FROM "sell_order";`)
}

func TestMain(m *testing.M) {
	pool, err := pgxpool.Connect(context.Background(), testDbURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}

	conn, err = pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire connection")
	}

	repository = &orderRepository{pool}

	m.Run()

	defer pool.Close()
	defer conn.Release()
}

func TestCreate(t *testing.T) {
	t.Cleanup(cleanup)
	dto := &orderDTO{
		amount:   100,
		price:    1000,
		assetID:  "a7beabee-bec2-4fc6-b86d-659b3b617562",
		sellerID: "b7beabee-bec2-4fc6-b86d-659b3b617562",
	}
	ctx := context.Background()

	initialCount := repository.Count(ctx)
	result := repository.Create(ctx, dto)
	finalCount := repository.Count(ctx)

	if result.id == "" {
		t.Error("Order should have an id")
	}

	if initialCount+1 != finalCount {
		t.Error("Expected rows count will increase after create")
	}
}
