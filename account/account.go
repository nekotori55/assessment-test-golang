package account

import "errors"

type account struct {
	id      string
	balance float64
}

func NewAccount() account {
	return account{
		id: GetNewID(),
	}
}

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

func (a *account) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("amount must be positive")
	}
	a.balance += amount
	return nil
}

func (a *account) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("amount must be positive")
	}
	if (a.balance - amount) < 0 {
		return errors.New("insufficient funds")
	}
	a.balance -= amount
	return nil
}

func (a *account) GetBalance() float64 {
	return a.balance
}
