package main

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var notFound error = errors.New("Account not found")

type Account struct {
	Id      string
	Balance int
}

type accountRepository struct {
	pool *pgxpool.Pool
}

func (r *accountRepository) Create(ctx context.Context) *Account {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		"INSERT INTO account (balance) VALUES ($1) RETURNING id",
		0)
	account := new(Account)
	err = row.Scan(&account.Id)
	if err != nil {
		log.Fatalf("Unable to INSERT: %v", err)
	}
	return account
}

func (r *accountRepository) getConnection(ctx context.Context) *pgxpool.Conn {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	return conn
}

func (r *accountRepository) Get(ctx context.Context, id string) (*Account, error) {
	conn := r.getConnection(ctx)
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT id, balance FROM account WHERE id = $1", id)

	account := new(Account)
	err := row.Scan(&account.Id, &account.Balance)
	if err == pgx.ErrNoRows {
		return nil, notFound
	}

	if err != nil {
		log.Fatalf("Unable to SELECT: %v", err)
	}

	return account, nil
}

func (r *accountRepository) Save(ctx context.Context, account *Account) error {
	conn := r.getConnection(ctx)
	defer conn.Release()

	ct, err := conn.Exec(ctx, "UPDATE account SET balance = $1 WHERE id = $2", account.Balance, account.Id)
	if err != nil {
		log.Fatalf("Unable to UPDATE: %v\n", err)
	}

	if ct.RowsAffected() == 0 {
		return notFound
	}

	return nil
}
