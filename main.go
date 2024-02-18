package main

import (
	"context"
	"log"
	"main/listeners"
	"main/routes"
	"main/routines"
	"net/http"
)

type Transacao struct {
	Valor     float64 `json:"valor"`
	Tipo      string  `json:"tipo"`
	Descricao string  `json:"descricao"`
}

type Resposta struct {
	Limite float64 `json:"limite"`
	Saldo  float64 `json:"saldo"`
}

var ctx = context.Background()

func main() {
	// Migrations
	routines.DatabaseMigration()
	routines.RedisMigration()

	// Listeners
	go listeners.TransactionsLazyWritting(ctx)
	go listeners.TransactionsLazyWritting(ctx)

	go listeners.TransactionsExtratoCache(ctx)
	go listeners.TransactionsExtratoCache(ctx)

	go listeners.TransactionsSaldoLazyWriter(ctx)

	http.HandleFunc("/clientes/", routes.ClientesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
