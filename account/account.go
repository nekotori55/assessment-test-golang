package account

import (
	"errors"
)

type account struct {
	id         string
	balance    float64
	operations chan Operation
}

type Operation struct {
	action string
	amount float64
	result chan float64
	err    chan error
}

func (a *account) ProcessOperations() {
	for op := range a.operations {
		if op.amount <= 0 {
			op.err <- errors.New("Amount must be greater than zero")
			continue
		}

		switch op.action {
		case "deposit":
			a.balance += op.amount
			op.err <- nil
		case "withdraw":
			if a.balance >= op.amount {
				a.balance -= op.amount
				op.err <- nil
			} else {
				op.err <- errors.New("Insufficient funds")
			}
		case "balance":
			op.result <- a.balance
			op.err <- nil
		}
	}
}

func NewAccount() account {
	a := account{
		id:         GetNewID(),
		operations: make(chan Operation),
	}
	// Запуск горутины по обработке операций
	go a.ProcessOperations()
	return a
}

func (a *account) GetID() string {
	// id - immutable,
	// поэтому потоковая безопасность не нужна
	return a.id
}

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

func (a *account) Deposit(amount float64) error {
	err := make(chan error)
	a.operations <- Operation{action: "deposit", amount: amount, err: err}
	return <-err
}

func (a *account) Withdraw(amount float64) error {
	err := make(chan error)
	a.operations <- Operation{action: "withdraw", amount: amount, err: err}
	return <-err
}

func (a *account) GetBalance() float64 {
	result := make(chan float64)
	a.operations <- Operation{action: "balance", result: result}
	return <-result
}
