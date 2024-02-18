package routines

import (
	"context"
	"fmt"
	"main/entities"
	"main/utils"
)

// Função para iniciar as chaves mock no Redis
func RedisMigration() {
	consumerName := "RedisMigration"

	ctx := context.Background()
	rdb := utils.GetRedisClient()
	db := utils.GetDB()

	queryTransacoes := `select id_cliente, saldo, limite from clientes order by id_cliente`
	rows, err := db.Query(queryTransacoes)
	if err != nil {
		fmt.Printf("[%s] Erro ao recuperar os clientes do database principal %v:\n", consumerName)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t entities.Cliente
		err := rows.Scan(&t.ID, &t.Saldo, &t.Limite)
		if err != nil {
			fmt.Printf("[%s] Erro ao consultar transações do database principal: %w\n", consumerName, err)
			continue // implementar o retry
		}

		// Atualiza os clientes no redis
		err = rdb.Set(ctx, "cliente:"+t.ID, fmt.Sprintf("%f", t.ID), 0).Err()
		if err != nil {
			fmt.Printf("[%s] Erro atualizar o cliente do database principal para o cache: %v\n", consumerName, err)
		}
		fmt.Printf("[%s] Cliente atualizado do database principal para o cache: %v\n", consumerName, t.ID)

		// Atualiza o saldo no Redis
		err = rdb.Set(ctx, "saldo:"+t.ID, fmt.Sprintf("%s", t.Saldo), 0).Err()
		if err != nil {
			fmt.Printf("[%s] Erro atualizar o saldo do database principal para o cache: %v\n", consumerName, err)
		}
		fmt.Printf("[%s] Saldo atualizado do database principal para o cache: %v\n", consumerName, t.ID)

		// Atualiza o limite no Redis
		err = rdb.Set(ctx, "limite:"+t.ID, fmt.Sprintf("%s", t.Limite), 0).Err()
		if err != nil {
			fmt.Printf("[%s] Erro atualizar o limite do database principal para o cache: %v\n", consumerName, err)
		}
		fmt.Printf("[%s] Limite atualizado do database principal para o cache: %v\n", consumerName, t.ID)

		_, err = rdb.LPush(ctx, "extrato-cache-command", t.ID).Result()
		_, err = rdb.LPush(ctx, "saldo-cache-command", t.ID).Result()
		_, err = rdb.LPush(ctx, "consistencia-cache-command", t.ID).Result()
	}
}
