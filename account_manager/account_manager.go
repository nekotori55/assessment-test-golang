package account_manager

import (
	acc "assessment-test/account"
	"errors"
)

type Operation struct {
	Action string
	Id     string
	Result chan interface{}
	Err    chan error
}

var accounts = make(map[string]acc.BankAccount)

func ProcessOperations(operations <-chan Operation) {
	for op := range operations {
		switch op.Action {
		case "get":
			res, found := accounts[op.Id]
			if found {
				op.Result <- res
				op.Err <- nil
			} else {
				op.Result <- nil
				op.Err <- errors.New("Account with id = " + op.Id + " not found")
			}
		case "add":
			newAccount := acc.NewAccount()
			id := newAccount.GetID()
			accounts[id] = &newAccount
			op.Result <- id
		}
	}
}
