package account

import "errors"

type Account struct {
	id      string
	balance float64
}

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

func (a Account) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("amount must be positive")
	}
	a.balance += amount
	return nil
}

func (a Account) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("amount must be positive")
	}
	if (a.balance - amount) < 0 {
		return errors.New("insufficient funds")
	}
	a.balance -= amount
	return nil
}

func (a Account) GetBalance() float64 {
	return a.balance
}
