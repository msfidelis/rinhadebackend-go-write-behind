package main

import (
	"context"
	"fmt"
	"log"
	"main/listeners"
	"main/routes"
	"main/routines"
	"net/http"
	"time"
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
	fmt.Println("Warm Databases")
	time.Sleep(10 * time.Second)

	// Migrations
	routines.DatabaseMigration()
	routines.RedisMigration()

	//Aplicando as estratégias de Write Behind
	// através de listeners de eventos em background

	// Listeners de Lazy Writting do Cache para o Databas
	go listeners.TransactionsLazyWritting(ctx)
	go listeners.TransactionsLazyWritting(ctx)

	// Listeners de Atualização do database para o Cache
	go listeners.TransactionsExtratoCache(ctx)
	go listeners.TransactionsExtratoCache(ctx)

	// Listeners de Lazy Writter para atualização do saldo
	go listeners.TransactionsSaldoLazyWriter(ctx)

	// Utilizando o HandleFunc do Go para otimização de performance
	http.HandleFunc("/clientes/", routes.ClientesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
