package main

import (
	"context"
	"fmt"
	"testing"
)

func TestAssetCreate(t *testing.T) {
	t.Cleanup(cleanup)
	asset := assetRepo.Create(context.Background())
	expected := "Asset 1"
	fmt.Println(asset.Name, expected)
	if asset.Name != expected {
		t.Errorf("Expected 'asset.Name' to be '%s', got '%s'", expected, asset.Name)
	}
}

func TestAssetGetBalanceNotExist(t *testing.T) {
	t.Cleanup(cleanup)
	asset := assetRepo.GetBalance(context.Background(), "foo", "bar")
	if asset != nil {
		t.Error("Expected balance to be nil")
	}
}

func TestAssetGetBalance(t *testing.T) {
	t.Cleanup(cleanup)

	account := accountRepo.Create(context.Background())
	asset := assetRepo.Create(context.Background())
	balanceValue := 500

	row := conn.QueryRow(context.Background(), `INSERT INTO "account_asset" ("asset_id", "account_id", "balance") VALUES ($1, $2, $3) RETURNING "balance"`,
		asset.Id, account.Id, balanceValue)

	var b int
	if err := row.Scan(&b); err != nil {
		t.Fatal(err)
	}

	balance := assetRepo.GetBalance(context.Background(), asset.Id, account.Id)
	if balance.Balance != balanceValue {
		t.Errorf("Expected 'balance.Balance' to be '%d', got '%d'", balanceValue, balance.Balance)
	}
}
