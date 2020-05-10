package main

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool
var conn *pgxpool.Conn
var repository *accountRepository

const testDbURL = "postgres://postgres:xchange@localhost/account_test?sslmode=disable&pool_max_conns=10"

func cleanup() {
	conn.Exec(context.Background(), "DELETE FROM account_asset;")
	conn.Exec(context.Background(), "DELETE FROM account;")
	conn.Exec(context.Background(), "DELETE FROM asset;")
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

	repository = &accountRepository{pool}

	m.Run()

	defer pool.Close()
	defer conn.Release()
}

func TestCreate(t *testing.T) {
	t.Cleanup(cleanup)
	account := repository.Create(context.Background())
	if len(account.Id) != 36 {
		t.Errorf("Expected `account.Id` to be uuid, got '%s'", account.Id)
	}
	if account.Balance != 0 {
		t.Error("Initial balance should be '0'")
	}
}

func TestGet(t *testing.T) {
	t.Cleanup(cleanup)
	account := repository.Create(context.Background())
	storedAccount, err := repository.Get(context.Background(), account.Id)
	if err != nil {
		t.Error(err)
	}
	if storedAccount.Id != account.Id {
		t.Errorf("Storted account id should be equal to created id, got '%s'", storedAccount.Id)
	}
	if storedAccount.Balance != account.Balance {
		t.Errorf("Storted account balance should be equal to created balance, got '%d'", storedAccount.Balance)
	}
}

func TestGetNotFound(t *testing.T) {
	t.Cleanup(cleanup)
	storedAccount, err := repository.Get(context.Background(), "a7beabee-bec2-4fc6-b86d-659b3b617562")
	if storedAccount != nil {
		t.Error("should return nil on non-existing account")
	}
	if err == nil {
		t.Error("should return error on non-existing account")
	}
	if err.Error() != "Account not found" {
		t.Error("'Account not found' error expected")
	}
}

func TestSave(t *testing.T) {
	t.Cleanup(cleanup)
	account := repository.Create(context.Background())
	account.Balance = 150
	err := repository.Save(context.Background(), account)
	if err != nil {
		t.Error(err)
	}
	storedAccount, err := repository.Get(context.Background(), account.Id)
	if err != nil {
		t.Error(err)
	}
	if storedAccount.Balance != account.Balance {
		t.Errorf("Expected to save balance, %d != %d", storedAccount.Balance, account.Balance)
	}
}

func TestSaveNonExisting(t *testing.T) {
	t.Cleanup(cleanup)
	account := new(Account)
	account.Id = "a7beabee-bec2-4fc6-b86d-659b3b617562"
	err := repository.Save(context.Background(), account)
	if err == nil {
		t.Error("Expected to return an error on non-exiting record")
	}
	if err.Error() != "Account not found" {
		t.Errorf("'Account not found' error expected")
	}
}
