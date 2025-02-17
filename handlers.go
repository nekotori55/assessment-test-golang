package main

import (
	accs "assessment-test/account_manager"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func WithdrawHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	amountParam := c.Params("amount")
	amount, err := strconv.ParseFloat(amountParam, 64)
	if (idParam == "") || (amountParam == "") || err != nil {
		log.Error("Withdraw. Bad parameters: [id: " + idParam + "] [amount: " + "]")
		return fiber.ErrBadRequest
	}

	account, err := findAccountByID(idParam)
	if err != nil {
		log.Error("id not found: " + idParam)
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	err = account.Withdraw(amount)
	if err != nil {
		log.Error("[id: " + idParam + "] " + err.Error())
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	log.Info("[id: " + idParam + "] successfully withdrew " + amountParam)
	return c.SendStatus(200)
}

func DepositHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	amountParam := c.Params("amount")
	amount, err := strconv.ParseFloat(amountParam, 64)
	if (idParam == "") || (amountParam == "") || err != nil {
		log.Error("Deposit. Bad parameters: [id: " + idParam + "] [amount: " + amountParam + "]")
		return fiber.ErrBadRequest
	}

	account, err := findAccountByID(idParam)
	if err != nil {
		log.Error("id not found: " + idParam)
		return fiber.NewError(fiber.StatusNotFound, "Error! Account with this id was not found")
	}

	err = account.Deposit(amount)

	if err != nil {
		log.Error("[id: " + idParam + "] " + err.Error())
		return fiber.NewError(fiber.StatusForbidden, "Error! "+err.Error())
	}

	log.Info("[id: " + idParam + "] successfully deposited " + amountParam)
	return c.SendStatus(200)
}

func GetBalanceHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return fiber.ErrBadRequest
	}

	account, err := findAccountByID(idParam)
	if err != nil {
		log.Error("id not found: " + idParam)
		return fiber.NewError(fiber.StatusNotFound, "Error! Account with this id was not found")
	}

	amount := account.GetBalance()

	amountStr := fmt.Sprint(amount)
	log.Info("[id: " + idParam + "] requested their balance, which is: " + amountStr)
	return c.SendString("[" + idParam + "]'s balance is: " + amountStr)
}

func CreateAccountHandler(c *fiber.Ctx) error {
	resChannel := make(chan interface{})
	accountStorageOperations <- accs.Operation{Action: "add", Result: resChannel}
	id := (<-resChannel).(string)

	log.Info("Account was created with id = " + id)
	return c.SendString("Successfully created an account with id: " + id)
}
