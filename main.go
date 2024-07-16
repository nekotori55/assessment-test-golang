package main

import (
	acc "assessment-test/account"
	"fmt"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"

	// Предоставляет удобное логирование с датой и временем
	"github.com/gofiber/fiber/v2/log"
	// "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Вполне можно использовать middleware для логирования встроенный в fiber
	// но тогда придется придумать как вытаскивать id из запроса в конкретных методах
	// app.Use(logger.New(logger.Config{
	// 	Format: "${time} ${method} ${status} ${path}\n",
	// }))

	app.Post("/accounts", CreateAccountHandler)

	app.Post("/accounts/:id/deposit/:amount", DepositHandler)

	app.Post("/accounts/:id/withdraw/:amount", WithdrawHandler)

	app.Get("/accounts/:id/balance", GetBalanceHandler)

	app.Listen(":3000")
}

var accounts = make(map[string]acc.BankAccount)
var accountsMutex sync.RWMutex

func findAccountByID(id string) (acc.BankAccount, bool) {
	accountsMutex.RLock()
	defer accountsMutex.RUnlock()
	acct, found := accounts[id]
	return acct, found
}

func WithdrawHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	amountParam := c.Params("amount")
	amount, err := strconv.ParseFloat(amountParam, 64)
	if (idParam == "") || (amountParam == "") || err != nil {
		log.Error("Withdraw. Bad parameters: [id: " + idParam + "] [amount: " + "]")
		return fiber.ErrBadRequest
	}

	account, found := findAccountByID(idParam)
	if !found {
		log.Error("id not found: " + idParam)
		return fiber.NewError(fiber.StatusNotFound, "account with this id was not found")
	}

	// Горутина использована чисто демонстрационно (для задания)
	// Т.к. по умолчанию фреймворк fiber запускает обработчики как горутины,
	// а логика не требует дополнительных горутин и, следовательно, каналов
	//
	// так же горутины можно было бы перенести в сами методы операций с аккаунтом,
	// в данном случае разницы никакой
	errChannel := make(chan error)
	go func() {
		errChannel <- account.Withdraw(amount)
	}()
	err = <-errChannel

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

	account, found := findAccountByID(idParam)
	if !found {
		log.Error("id not found: " + idParam)
		return fiber.NewError(fiber.StatusNotFound, "Error! Account with this id was not found")
	}

	errChannel := make(chan error)
	go func() {
		errChannel <- account.Deposit(amount)
	}()
	err = <-errChannel

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

	account, found := findAccountByID(idParam)
	if !found {
		log.Error("id not found: " + idParam)
		return fiber.NewError(fiber.StatusNotFound, "Error! Account with this id was not found")
	}

	result := make(chan float64)
	go func() {
		result <- account.GetBalance()
	}()
	amount := <-result

	amountStr := fmt.Sprint(amount)
	log.Info("[id: " + idParam + "] requested their balance, which is: " + amountStr)
	return c.SendString("[" + idParam + "]'s balance is: " + amountStr)
}

func CreateAccountHandler(c *fiber.Ctx) error {
	accountsMutex.Lock()
	defer accountsMutex.Unlock()

	newAccount := acc.NewAccount()
	id := newAccount.GetID()
	accounts[id] = &newAccount

	log.Info("Account was created with id = " + id)
	return c.SendString("Successfully created an account with id: " + id)
}
