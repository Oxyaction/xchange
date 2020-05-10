package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type orderRepository struct {
	pool *pgxpool.Pool
}

type orderDTO struct {
	assetID  string
	amount   int
	price    int
	sellerID string
}

type Order struct {
	orderDTO
	id        string
	createdAt time.Time
}

func (r *orderRepository) Create(ctx context.Context, dto *orderDTO) *Order {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO
		"sell_order" (asset_id, amount, price, seller_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, asset_id, amount, price, seller_id, created_at`,
		dto.assetID, dto.amount, dto.price, dto.sellerID)

	order := new(Order)
	err = row.Scan(&order.id, &order.assetID, &order.amount, &order.price, &order.sellerID, &order.createdAt)

	if err != nil {
		log.Fatalf("Unable to INSERT: %v", err)
	}
	return order
}

func (r *orderRepository) Count(ctx context.Context) int {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT COUNT(*) FROM "sell_order"`)

	var count int
	err = row.Scan(&count)

	if err != nil {
		log.Fatalf("Unable to INSERT: %v", err)
	}

	return count
}
