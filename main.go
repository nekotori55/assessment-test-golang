package main

import (
	acc "assessment-test/account"
	accs "assessment-test/account_manager"

	"github.com/gofiber/fiber/v2"
	// Предоставляет удобное логирование с датой и временем
	// "github.com/gofiber/fiber/v2/middleware/logger"
)

var accountStorageOperations = make(chan accs.Operation)

func main() {
	app := fiber.New()

	// Вполне можно использовать middleware для логирования встроенный в fiber
	// но тогда придется придумать как вытаскивать id из запроса в конкретных методах
	// app.Use(logger.New(logger.Config{
	// 	Format: "${time} ${method} ${status} ${path}\n",
	// }))
	//

	// Запуск горутины по обработке операций с хранением аккаунтов
	go accs.ProcessOperations(accountStorageOperations)

	app.Post("/accounts", CreateAccountHandler)

	app.Post("/accounts/:id/deposit/:amount", DepositHandler)

	app.Post("/accounts/:id/withdraw/:amount", WithdrawHandler)

	app.Get("/accounts/:id/balance", GetBalanceHandler)

	app.Listen(":3000")
}

func findAccountByID(id string) (acc.BankAccount, error) {
	errChannel := make(chan error)
	resChannel := make(chan interface{})
	accountStorageOperations <- accs.Operation{Action: "get", Id: id, Result: resChannel, Err: errChannel}

	res := <-resChannel
	if res == nil {
		return nil, <-errChannel
	}

	return res.(acc.BankAccount), <-errChannel
}
