package routes

import (
	"encoding/json"
	"fmt"
	"main/utils"
	"net/http"
	"strings"
)

type Request struct {
	Valor     float64 `json:"valor"`
	Tipo      string  `json:"tipo"`
	Descricao string  `json:"descricao"`
}

type Resposta struct {
	Limite float64 `json:"limite"`
	Saldo  float64 `json:"saldo"`
}

func TransacaoHandler(w http.ResponseWriter, r *http.Request) {

	featureName := "NewTransaction"

	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Split manual da URL, teste de performance
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) != 4 || pathSegments[3] != "transacoes" {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	clienteID := string(pathSegments[2])

	rdb := utils.GetRedisClient()
	db := utils.GetDB()
	cache := utils.GetCacheInstance()

	// Checa no cache em memória da aplicação se o cliente existe
	_, found := cache.Get("cliente:" + clienteID)
	if !found {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Decode da Transação -> Maior ofensor de performance até o momento
	var transacao Request
	if err := json.NewDecoder(r.Body).Decode(&transacao); err != nil {
		utils.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atualização de saldo
	// limite, saldo, semlimite, err := utils.AtualizarSaldo(r.Context(), rdb, clienteID, transacao.Valor, transacao.Tipo)
	limite, saldo, semlimite, err := utils.AtualizarSaldo(r.Context(), db, clienteID, transacao.Valor, transacao.Tipo)
	if semlimite {
		fmt.Printf("[%s] Cliente sem limite disponível: %w\n", featureName, clienteID)
		utils.HttpError(w, "Cliente sem limite disponível", http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		fmt.Printf("[%s] Erro ao processar a transação: %w\n", featureName, err)
		utils.HttpError(w, "Erro ao processar a transação", http.StatusInternalServerError)
		return
	}

	resposta := Resposta{Limite: limite, Saldo: saldo}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resposta); err != nil {
		utils.HttpError(w, "Erro ao enviar resposta", http.StatusInternalServerError)
		return
	}

	// Adiciona a transação à fila de processamento assíncrono - Lazy Writting
	message := fmt.Sprintf("%v:%v:%v:%v", clienteID, transacao.Descricao, transacao.Tipo, transacao.Valor)
	if err := rdb.LPush(r.Context(), "transactions", message).Err(); err != nil {
		fmt.Printf("[%s] Erro ao publicar a mensagem na fila de persistência: %v\n", featureName, err)
	}
}
