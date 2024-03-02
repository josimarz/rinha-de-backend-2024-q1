package entity

import "time"

type Transaction struct {
	AccountId uint8
	Amount    int64
	Kind      string
	Desc      string
	Timestamp time.Time
}

func NewTransaction(accountId uint8, amount int64, kind, desc string) *Transaction {
	return &Transaction{
		AccountId: accountId,
		Amount:    amount,
		Kind:      kind,
		Desc:      desc,
		Timestamp: time.Now(),
	}
}
