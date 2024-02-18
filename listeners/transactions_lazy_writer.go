package listeners

import (
	"context"
	"fmt"
	"main/utils"
	"strings"
	"time"
)

func TransactionsLazyWritting(ctx context.Context) {

	consumerName := "TransactionsLazyWritting"
	fmt.Printf("[%s] Iniciando Consumer de Lazy Writting Transactions\n", consumerName)

	rdb := utils.GetRedisClient()
	db := utils.GetDB()

	for {
		result, err := rdb.BRPop(ctx, 0, "transactions").Result()
		if err != nil {
			fmt.Printf("[%s] Erro ao recuperar o registro da queue: %w\n", consumerName, err)
			continue
		}
		fmt.Printf("[%s] Mensagem recebida na queue %v: %v\n", consumerName, result[0], result[1])

		messageSplit := strings.Split(result[1], ":")

		idCliente := messageSplit[0]
		descricao := messageSplit[1]
		tipo := messageSplit[2]
		valor := messageSplit[3]

		realizadaEm := time.Now().UTC().Format(time.RFC3339Nano)

		query := `INSERT INTO transacoes (id_cliente, valor, descricao, tipo, realizada_em) VALUES ($1, $2, $3, $4, $5)`
		_, err = db.Exec(query, idCliente, valor, descricao, tipo, realizadaEm)
		if err != nil {
			fmt.Printf("[%s] Erro ao inserir a transação no database principal: %v\n", consumerName, err)
			fmt.Printf("[%s] Adicionando a mensagem novamente a queue - retry: %v\n", consumerName, result[1])

			// Adicionando mensagem novamente a queue de Lazy Writting
			_, err = rdb.LPush(ctx, "transactions", result[1]).Result()
			if err != nil {
				fmt.Printf("[%s]Erro ao enviar a mensagem %v para a queue: %v\n", consumerName, result[1], err)
			}
		} else {
			fmt.Printf("[%s] Transação inserida no banco de dados principal com sucesso: %v\n", consumerName, result[1])
		}

		_, err = rdb.LPush(ctx, "extrato-cache-command", idCliente).Result()
		_, err = rdb.LPush(ctx, "saldo-cache-command", idCliente).Result()
		_, err = rdb.LPush(ctx, "consistencia-cache-command", idCliente).Result()

	}
}
