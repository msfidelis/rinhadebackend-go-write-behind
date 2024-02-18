package routes

import (
	"encoding/json"
	"fmt"
	"main/entities"
	"main/utils"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ExtratoResponse struct {
	Saldo             Saldo                `json:"saldo"`
	UltimasTransacoes []entities.Transacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       float64 `json:"total"`
	DataExtrato string  `json:"data_extrato"`
	Limite      float64 `json:"limite"`
}

func ExtratoHandler(w http.ResponseWriter, r *http.Request) {

	featureName := "Extrato"

	var wg sync.WaitGroup
	var response ExtratoResponse
	var limite, saldo float64
	var transacoes []entities.Transacao
	var errSaldo, errTransacoes error

	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) != 4 || pathSegments[3] != "extrato" {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	clienteID := pathSegments[2]

	rdb := utils.GetRedisClient()

	// Checa se o cliente existe
	cliente, err := rdb.Get(ctx, "cliente:"+clienteID).Result()
	if err != nil || cliente == "" {
		fmt.Println(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	fmt.Printf("[%s] cliente encontrado: %v\n", featureName, clienteID)

	wg.Add(2)

	fmt.Printf("[%s] Recuperando saldo da camada de cache: %v\n", featureName, clienteID)
	// Recupera o saldo e o limite na camada de cache
	go func() {
		defer wg.Done()
		limite, saldo, errSaldo = utils.RecuperarSaldoELimite(ctx, rdb, clienteID)
	}()

	fmt.Printf("[%s] Recuperando transações da camada de cache: %v\n", featureName, clienteID)
	// Recupera as ultimas transações do client
	go func() {
		defer wg.Done()
		transacoes, errTransacoes = utils.RecuperarTransacoes(ctx, rdb, clienteID)
	}()
	fmt.Printf("[%s] Transações recuperadas: %v\n", featureName, clienteID)

	wg.Wait()

	if errSaldo != nil {
		fmt.Printf("[%s] Erro ao recuperar saldo e limite da camada de cache: %v\n", featureName, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if errTransacoes != nil {
		fmt.Printf("[%s] Erro ao recuperar transações da camada de cache: %v\n", featureName, errTransacoes)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	response.Saldo.Total = saldo
	response.Saldo.Limite = limite
	response.Saldo.DataExtrato = time.Now().UTC().Format(time.RFC3339Nano)
	response.UltimasTransacoes = transacoes

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
