package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/entity"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/gateway"
)

type AccountNotFoundError struct{}

func (e *AccountNotFoundError) Error() string {
	return "account not found"
}

type BalanceOutput struct {
	Total     int64     `json:"total"`
	Timestamp time.Time `json:"data_extrato"`
	Limit     int64     `json:"limite"`
}

type TransactionOutput struct {
	Amount    int64     `json:"valor"`
	Kind      string    `json:"tipo"`
	Desc      string    `json:"descricao"`
	Timestamp time.Time `json:"realizada_em"`
}

type GetBankStatementOutput struct {
	Balance      *BalanceOutput       `json:"saldo"`
	Transactions []*TransactionOutput `json:"ultimas_transacoes"`
}

type GetBankStatementUseCase struct {
	dbGateway gateway.DatabaseGateway
}

func NewGetBankStatementUseCase(dbGateway gateway.DatabaseGateway) *GetBankStatementUseCase {
	return &GetBankStatementUseCase{dbGateway}
}

func (uc *GetBankStatementUseCase) Execute(id uint8) (*GetBankStatementOutput, error) {
	account, err := uc.dbGateway.GetAccount(id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, &AccountNotFoundError{}
	}
	transactions, err := uc.dbGateway.ListTransactions(id)
	if err != nil {
		return nil, err
	}
	o := &GetBankStatementOutput{
		Balance: &BalanceOutput{
			Total:     account.Balance,
			Timestamp: time.Now(),
			Limit:     account.Limit,
		},
		Transactions: []*TransactionOutput{},
	}
	for _, t := range transactions {
		o.Transactions = append(o.Transactions, &TransactionOutput{
			Amount:    t.Amount,
			Kind:      t.Kind,
			Desc:      t.Desc,
			Timestamp: t.Timestamp,
		})
	}
	return o, nil
}

type DoTransactionInput struct {
	AccountId uint8
	Amount    int64
	Kind      string
	Desc      string
}

type DoTransactionOuput struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

type DoTransactionUseCase struct {
	dbGateway gateway.DatabaseGateway
	rs        *redsync.Redsync
}

func NewDoTransactionUseCase(dbGateway gateway.DatabaseGateway, rs *redsync.Redsync) *DoTransactionUseCase {
	return &DoTransactionUseCase{dbGateway, rs}
}

func (uc *DoTransactionUseCase) Execute(input *DoTransactionInput) (*DoTransactionOuput, error) {
	name := fmt.Sprintf("%v", input.AccountId)
	mu := uc.rs.NewMutex(name)
	if err := mu.Lock(); err != nil {
		return nil, err
	}
	defer func() {
		if ok, err := mu.Unlock(); !ok || err != nil {
			log.Fatal(err)
		}
	}()
	account, err := uc.dbGateway.GetAccount(input.AccountId)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, &AccountNotFoundError{}
	}
	if input.Kind == "c" {
		account.Deposit(input.Amount)
	}
	if input.Kind == "d" {
		if err := account.Withdraw(input.Amount); err != nil {
			return nil, err
		}
	}
	t := entity.NewTransaction(input.AccountId, input.Amount, input.Kind, input.Desc)
	if err := uc.dbGateway.UpdateBalance(account); err != nil {
		return nil, err
	}
	if err := uc.dbGateway.SaveTransaction(t); err != nil {
		return nil, err
	}
	o := &DoTransactionOuput{
		Limit:   account.Limit,
		Balance: account.Balance,
	}
	return o, nil
}
