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

	var wg sync.WaitGroup
	var response ExtratoResponse
	var limite, saldo float64
	var transacoes []entities.Transacao
	var err1, err2 error

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

	wg.Add(2)

	// Recupera o saldo e o limite na camada de cache
	go func() {
		defer wg.Done()
		limite, saldo, err1 = utils.RecuperarSaldoELimite(ctx, rdb, clienteID)
	}()

	// Recupera as ultimas transações do client
	go func() {
		defer wg.Done()
		transacoes, err2 = utils.RecuperarTransacoes(ctx, rdb, clienteID)
	}()

	wg.Wait()

	if err1 != nil {
		fmt.Println("Erro ao recuperar saldo e limite:", err1)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err2 != nil {
		fmt.Println("Erro ao recuperar as ultimas transações:", err1)
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
