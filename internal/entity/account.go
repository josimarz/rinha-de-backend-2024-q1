package entity

type InsuficientBalanceError struct{}

func (e *InsuficientBalanceError) Error() string {
	return "insuficient balance"
}

type Account struct {
	Id      uint8
	Limit   int64
	Balance int64
}

func (e *Account) Deposit(amount int64) {
	e.Balance += amount
}

func (e *Account) Withdraw(amount int64) error {
	e.Balance -= amount
	if e.Limit+e.Balance < 0 {
		return &InsuficientBalanceError{}
	}
	return nil
}
