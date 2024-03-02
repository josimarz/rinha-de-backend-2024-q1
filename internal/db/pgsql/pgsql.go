package pgsql

import (
	"database/sql"

	"github.com/josimarz/rinha-de-backend-2024-q1/internal/entity"
)

type PostgreDatabaseGateway struct {
	db *sql.DB
}

func NewPostgreDatabaseGateway(db *sql.DB) *PostgreDatabaseGateway {
	return &PostgreDatabaseGateway{db}
}

func (g *PostgreDatabaseGateway) GetAccount(id uint8) (*entity.Account, error) {
	stmt, err := g.db.Prepare(`select "limit", "balance" from "account" where "id" = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	account := &entity.Account{
		Id: id,
	}
	if err := stmt.QueryRow(id).Scan(&account.Limit, &account.Balance); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return account, nil
}

func (g *PostgreDatabaseGateway) UpdateBalance(a *entity.Account) error {
	stmt, err := g.db.Prepare(`update "account" set "balance" = $1 where "id" = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Balance, a.Id)
	return err
}

func (g *PostgreDatabaseGateway) ListTransactions(id uint8) ([]*entity.Transaction, error) {
	stmt, err := g.db.Prepare(`select "amount", "kind", "description", "timestamp" from "transaction" where "accountId" = $1 order by "timestamp" desc limit 10`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions := []*entity.Transaction{}
	for rows.Next() {
		t := &entity.Transaction{}
		if err := rows.Scan(&t.Amount, &t.Kind, &t.Desc, &t.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (g *PostgreDatabaseGateway) SaveTransaction(t *entity.Transaction) error {
	stmt, err := g.db.Prepare(`insert into "transaction" ("accountId", "amount", "kind", "description", "timestamp") values ($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.AccountId, t.Amount, t.Kind, t.Desc, t.Timestamp)
	return err
}
