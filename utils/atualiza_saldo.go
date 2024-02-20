package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/entities"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func AtualizarSaldo(ctx context.Context, db *sql.DB, clienteID string, valorTransacao float64, tipo string) (limite float64, saldoAtualizado float64, limiteExcedido bool, err error) {
	tx, err := db.BeginTx(ctx, nil)
	rdb := GetRedisClient()

	if err != nil {
		return 0, 0, false, fmt.Errorf("erro ao iniciar transação: %v", err)
	}

	var saldo float64
	err = tx.QueryRowContext(ctx, "SELECT limite, saldo FROM clientes WHERE id_cliente = $1", clienteID).Scan(&limite, &saldo)
	if err != nil {
		tx.Rollback()
		return 0, 0, false, fmt.Errorf("erro ao recuperar dados do cliente %s: %v", clienteID, err)
	}

	switch tipo {
	case "c":
		saldo += valorTransacao
	case "d":
		saldo -= valorTransacao
	default:
		tx.Rollback()
		return 0, 0, false, fmt.Errorf("tipo de transação inválido: %s", tipo)
	}

	if saldo < -limite {
		tx.Rollback()
		return limite, saldo, true, fmt.Errorf("transação excede o limite do cliente")
	}

	_, err = tx.ExecContext(ctx, "UPDATE clientes SET saldo = $1 WHERE id_cliente = $2", saldo, clienteID)
	if err != nil {
		tx.Rollback()
		return 0, 0, false, fmt.Errorf("erro ao atualizar saldo no banco de dados: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, false, fmt.Errorf("erro ao finalizar transação: %v", err)
	}
	if err = rdb.Set(ctx, "saldo:"+clienteID, fmt.Sprintf("%f", saldo), 0).Err(); err != nil {
		return limite, saldo, false, fmt.Errorf("erro ao atualizar saldo no Redis: %v", err)
	}
	if err = rdb.Set(ctx, "limite:"+clienteID, fmt.Sprintf("%f", limite), 0).Err(); err != nil {
		return limite, saldo, false, fmt.Errorf("erro ao atualizar saldo no Redis: %v", err)
	}
	return limite, saldo, false, nil
}

// Recupera o saldo e o limite do cliente informado da camada de cache
func RecuperarSaldoELimite(ctx context.Context, rdb *redis.Client, clienteID string) (float64, float64, error) {

	valores, err := rdb.MGet(ctx, "limite:"+clienteID, "saldo:"+clienteID).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao recuperar dados do cliente %s: %v", clienteID, err)
	}

	limite, _ := strconv.ParseFloat(valores[0].(string), 64)
	saldo, _ := strconv.ParseFloat(valores[1].(string), 64)

	return limite, saldo, nil
}

// RecuperarTransacoes recupera as últimas 10 transações do cliente da camada de cache
func RecuperarTransacoes(ctx context.Context, rdb *redis.Client, clienteID string) ([]entities.Transacao, error) {
	var transacoes []entities.Transacao

	transacoesCache, err := rdb.Get(ctx, "extrato:"+clienteID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return transacoes, nil
		}
		log.Printf("Erro ao recuperar transações do cache para cliente %s: %v\n", clienteID, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(transacoesCache), &transacoes); err != nil {
		log.Printf("Erro ao deserializar transações para cliente %s: %v\n", clienteID, err)
		return nil, err
	}

	return transacoes, nil
}
