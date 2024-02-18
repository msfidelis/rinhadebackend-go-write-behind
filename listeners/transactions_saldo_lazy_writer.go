package listeners

import (
	"context"
	"fmt"
	"main/utils"
	"strconv"
)

func TransactionsSaldoLazyWriter(ctx context.Context) {
	consumerName := "TransactionsSaldoLazyWriter"
	fmt.Printf("[%s] Iniciando Consumer de Lazy Writting Transactions\n", consumerName)

	rdb := utils.GetRedisClient()
	db := utils.GetDB()

	for {
		result, err := rdb.BRPop(ctx, 0, "saldo-cache-command").Result()
		if err != nil {
			fmt.Printf("[%s] Erro ao recuperar o registro da queue: %w\n", consumerName, err)
			continue
		}

		fmt.Printf("[%s] Mensagem recebida na queue %v: %v\n", consumerName, result[0], result[1])

		// Busca o saldo atual do cliente no Redis
		saldoStr, err := rdb.Get(ctx, "saldo:"+result[1]).Result()
		if err != nil {
			fmt.Printf("[%s] Erro ao recuperar o saldoe: %w\n", err)
		}
		saldo, _ := strconv.ParseFloat(saldoStr, 64)

		// Iniciando a Transação
		tx, err := db.Begin()
		if err != nil {
			fmt.Printf("[%s] Erro ao iniciar a transação no banco de dados: %w\n", err)
		}

		query := `UPDATE clientes SET saldo = $1 WHERE id_cliente = $2`
		_, err = tx.Exec(query, saldo, result[1])
		if err != nil {
			tx.Rollback()
			fmt.Printf("[%s] Erro ao consultar transações do cliente %v: %v\n", consumerName, result[0], result[1])
			return
		}

		err = tx.Commit()
		if err != nil {
			fmt.Printf("[%s] Erro ao commitar as transações no banco de dados principal %v: %v\n", consumerName, result[0], result[1])
		}

		fmt.Printf("[%s] Saldo atualizado no banco de dados principal %v: %v\n", consumerName, result[0], result[1])
	}
}
