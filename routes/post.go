package routes

import (
	"context"
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

var ctx = context.Background()

func TransacaoHandler(w http.ResponseWriter, r *http.Request) {

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

	// Checa se o cliente existe
	cliente, err := rdb.Get(ctx, "cliente:"+clienteID).Result()
	if err != nil || cliente == "" {
		fmt.Println(err)
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
	limite, saldo, semlimite, err := utils.AtualizarSaldo(ctx, rdb, clienteID, transacao.Valor, transacao.Tipo)
	if semlimite {
		utils.HttpError(w, "cliente sem limite disponível", http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		fmt.Println(err)
		utils.HttpError(w, "Erro ao processar a transação", http.StatusInternalServerError)
		return
	}

	resposta := Resposta{Limite: limite, Saldo: saldo}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resposta)

	message := fmt.Sprintf("%v:%v:%v:%v", clienteID, transacao.Descricao, transacao.Tipo, transacao.Valor)

	// Adicionando mensagem a queue de Lazy Writting
	_, err = rdb.LPush(ctx, "transactions", message).Result()
	if err != nil {
		fmt.Printf("error pushing to queue: %w", err)
	}
}
