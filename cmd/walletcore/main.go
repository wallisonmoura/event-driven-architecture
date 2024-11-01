package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wallison/fc-ms-wallet/internal/database"
	"github.com/wallison/fc-ms-wallet/internal/event"
	"github.com/wallison/fc-ms-wallet/internal/usecase/create_account"
	"github.com/wallison/fc-ms-wallet/internal/usecase/create_client"
	"github.com/wallison/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/wallison/fc-ms-wallet/internal/web"
	"github.com/wallison/fc-ms-wallet/internal/web/webserver"
	"github.com/wallison/fc-ms-wallet/pkg/events"
	"github.com/wallison/fc-ms-wallet/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)
	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createaccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createaccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}
