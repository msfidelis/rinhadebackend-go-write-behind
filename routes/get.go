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
	cache := utils.GetCacheInstance()

	// Checa no cache em memória da aplicação se o cliente existe
	_, found := cache.Get("cliente:" + clienteID)
	if found == false {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	wg.Add(2)

	// Recupera o saldo e o limite na camada de cache
	go func() {
		defer wg.Done()
		response.Saldo.Limite, response.Saldo.Total, errSaldo = utils.RecuperarSaldoELimite(ctx, rdb, clienteID)
	}()

	// Recupera as ultimas transações do client
	go func() {
		defer wg.Done()
		response.UltimasTransacoes, errTransacoes = utils.RecuperarTransacoes(ctx, rdb, clienteID)
	}()

	wg.Wait()

	if errSaldo != nil {
		fmt.Printf("[%s] Erro ao recuperar saldo e limite da camada de cache: %v\n", featureName, errSaldo)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if errTransacoes != nil {
		fmt.Printf("[%s] Erro ao recuperar transações da camada de cache: %v\n", featureName, errTransacoes)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	response.Saldo.DataExtrato = time.Now().UTC().Format(time.RFC3339Nano)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
