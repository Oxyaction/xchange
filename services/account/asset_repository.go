package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type assetRepository struct {
	pool *pgxpool.Pool
}

type Asset struct {
	Id   string
	Name string
}

type AssetBalance struct {
	AssetId   string
	AccountId string
	Balance   int
}

var InsufficientFunds error = errors.New("Insufficient funds")

func (r *assetRepository) Create(ctx context.Context) *Asset {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT COUNT(*) FROM "asset"`)
	var maxId int
	row.Scan(&maxId)

	assetName := fmt.Sprintf("Asset %d", maxId+1)

	row = conn.QueryRow(ctx, `INSERT INTO "asset" (name) VALUES ($1) RETURNING id`, assetName)
	asset := new(Asset)
	err = row.Scan(&asset.Id)
	asset.Name = assetName
	if err != nil {
		log.Fatalf("Unable to INSERT: %v", err)
	}
	return asset
}

func (r *assetRepository) GetBalance(ctx context.Context, assetId string, accountId string) *AssetBalance {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT "balance" FROM "account_asset" WHERE "account_id" = $1 AND "asset_id" = $2`,
		accountId, assetId)
	var balance int
	err = row.Scan(&balance)

	if err != nil {
		return nil
	}

	return &AssetBalance{
		AssetId:   assetId,
		AccountId: accountId,
		Balance:   balance,
	}
}

func (r *assetRepository) ChangeBalance(ctx context.Context, balance *AssetBalance) (*AssetBalance, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Release()

	currentBalance := r.GetBalance(ctx, balance.AssetId, balance.AccountId)
	if currentBalance == nil {
		row := conn.QueryRow(ctx, `INSERT INTO "account_asset" ("asset_id", "account_id", "balance") VALUES ($1, $2, $3)`,
			balance.AssetId, balance.AccountId, 0)
		if err := row.Scan(); err != nil {
			return nil, errors.New("Can not create balance record")
		}
		currentBalance = &AssetBalance{
			AccountId: balance.AccountId,
			AssetId:   balance.AssetId,
			Balance:   0,
		}
	}

	if (currentBalance.Balance + balance.Balance) < 0 {
		return nil, InsufficientFunds
	}

	row := conn.QueryRow(
		ctx,
		`UPDATE "account_asset" SET "balance" = "balance" + $1 WHERE "account_id" = $2 AND "asset_id" = $3 RETURNING "balance"`,
		balance.Balance,
		balance.AccountId,
		balance.AssetId,
	)

	var balanceValue int
	if err := row.Scan(&balanceValue); err != nil {
		return nil, errors.New("Can not update balance")
	}

	currentBalance.Balance = balanceValue
	return currentBalance, nil
}
