package main

import (
	acc "assessment-test/account"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/accounts", CreateAccountHandler)

	app.Post("/accounts/:id/deposit/:amount", DepositHandler)

	app.Post("/accounts/:id/withdraw/:amount", WithdrawHandler)

	app.Get("/accounts/:id/balance", GetBalanceHandler)

	app.Listen(":3000")
}

var accounts = make(map[string]acc.BankAccount)

func WithdrawHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	amountParam := c.Params("amount")
	amount, err := strconv.ParseFloat(amountParam, 64)
	if (idParam == "") || (amountParam == "") || err != nil {
		return fiber.ErrBadRequest
	}

	account, found := accounts[idParam]
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "account with this id was not found")
	}

	// Горутина использована чисто демонстрационно (для задания)
	// Т.к. по умолчанию фреймворк fiber запускает обработчики как горутины,
	// а логика не требует дополнительных горутин и, следовательно, каналов
	errChannel := make(chan error)
	go func() {
		errChannel <- account.Withdraw(amount)
	}()
	err = <-errChannel

	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	return c.SendStatus(200)
}

func DepositHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	amountParam := c.Params("amount")
	amount, err := strconv.ParseFloat(amountParam, 64)
	if (idParam == "") || (amountParam == "") || err != nil {
		return fiber.ErrBadRequest
	}

	account, found := accounts[idParam]
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "account with this id was not found")
	}

	errChannel := make(chan error)
	go func() {
		errChannel <- account.Deposit(amount)
	}()
	err = <-errChannel

	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}
	return c.SendStatus(200)
}

func GetBalanceHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return fiber.ErrBadRequest
	}

	account, found := accounts[idParam]
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "account with this id was not found")
	}

	result := make(chan float64)
	go func() {
		result <- account.GetBalance()
	}()
	amount := <-result
	return c.SendString("[" + idParam + "] balance is: " + fmt.Sprint(amount))
}

func CreateAccountHandler(c *fiber.Ctx) error {
	accountChannel := make(chan acc.BankAccount)
	idChannel := make(chan string)

	go func() {
		newAccount := acc.NewAccount()

		idChannel <- newAccount.GetID()
		accountChannel <- &newAccount
	}()

	id := <-idChannel
	// В данном случае использование каналов оправдано,
	// т.к. maps не являются thread-safe объектами
	// a канал ждёт "доступности" и приёмника и передатчика
	accounts[id] = <-accountChannel

	return c.SendString("Successfully created an account with id: " + id)
}
