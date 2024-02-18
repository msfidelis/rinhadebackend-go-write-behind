package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/entities"
	"main/utils"
)

// Este listener captura eventos de insert de registros através de uma fila de comando
// De forma proativa, as ultimas 10 transações do cliente informado são colocadas em cache
func TransactionsExtratoCache(ctx context.Context) {

	consumerName := "TransactionsExtratoCache"
	fmt.Printf("[%s] Iniciando Consumer de Caching de Transactions\n", consumerName)

	rdb := utils.GetRedisClient()
	db := utils.GetDB()

	for {
		result, err := rdb.BRPop(ctx, 0, "extrato-cache-command").Result()
		if err != nil {
			fmt.Printf("[%s] Erro ao recuperar o registro da queue: %w\n", consumerName, err)
			continue
		}

		fmt.Printf("[%s] Mensagem recebida na queue %v: %v\n", consumerName, result[0], result[1])

		idCliente := result[1]

		queryTransacoes := `SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE id_cliente = $1 ORDER BY realizada_em DESC LIMIT 10`
		rows, err := db.Query(queryTransacoes, idCliente)
		if err != nil {
			fmt.Printf("[%s] Erro ao consultar transações do cliente %v: %v\n", consumerName, result[0], result[1])
			return
		}
		defer rows.Close()

		// Preenche a lista de transações
		var transacoes []entities.Transacao
		for rows.Next() {
			var t entities.Transacao
			err := rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm)
			if err != nil {
				fmt.Printf("[%s] Erro ao consultar transações do cliente %v: %v\n", consumerName, result[0], result[1])
				continue // implementar o retry
			}
			transacoes = append(transacoes, t)
		}

		// Convertendo a lista de transações em uma String JSON
		transacoesJSON, err := json.Marshal(transacoes)
		if err != nil {
			fmt.Printf("[%s] Erro ao converter a lista de transações para JSON %v: %v\n", consumerName, result[0], result[1])
		}

		// Salva o extrato em cache
		extratoKey := fmt.Sprintf("extrato:%s", idCliente)
		err = rdb.Set(ctx, extratoKey, transacoesJSON, 0).Err()
		if err != nil {
			log.Fatalf("Erro ao criar a chave de cliente %s: %v", idCliente, err)
		}
	}
}
