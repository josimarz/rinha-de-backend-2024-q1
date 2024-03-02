package gateway

import "github.com/josimarz/rinha-de-backend-2024-q1/internal/entity"

type DatabaseGateway interface {
	GetAccount(uint8) (*entity.Account, error)
	UpdateBalance(*entity.Account) error
	ListTransactions(uint8) ([]*entity.Transaction, error)
	SaveTransaction(*entity.Transaction) error
}
