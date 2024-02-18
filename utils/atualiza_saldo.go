package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"main/entities"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func AtualizarSaldo(ctx context.Context, rdb *redis.Client, clienteID string, valorTransacao float64, tipo string) (limite float64, saldoAtualizado float64, limiteExcedido bool, err error) {
	var saldo float64

	// Recupera o limite e o saldo atual do cliente de forma eficiente
	valores, err := rdb.MGet(ctx, "limite:"+clienteID, "saldo:"+clienteID).Result()
	if err != nil {
		return 0, 0, false, fmt.Errorf("erro ao recuperar dados do cliente %s: %v", clienteID, err)
	}

	limite, err = strconv.ParseFloat(valores[0].(string), 64)
	if err != nil {
		return 0, 0, false, fmt.Errorf("erro ao converter limite para float: %v", err)
	}

	if valores[1] != nil { // Verifica se o saldo existe
		saldo, err = strconv.ParseFloat(valores[1].(string), 64)
		if err != nil {
			return 0, 0, false, fmt.Errorf("erro ao converter saldo para float: %v", err)
		}
	}

	switch tipo {
	case "c":
		saldo += valorTransacao
	case "d":
		saldo -= valorTransacao
	default:
		return 0, 0, false, fmt.Errorf("tipo de transação inválido: %s", tipo)
	}

	if saldo < -limite {
		return limite, saldo, true, nil // Indica que o limite seria excedido
	}

	// Atualiza o saldo no Redis
	if err = rdb.Set(ctx, "saldo:"+clienteID, fmt.Sprintf("%f", saldo), 0).Err(); err != nil {
		return limite, saldo, false, fmt.Errorf("erro ao atualizar saldo no Redis: %v", err)
	}

	return limite, saldo, false, nil
}

// Recupera o saldo e o limite do cliente informado da camada de cache
func RecuperarSaldoELimite(ctx context.Context, rdb *redis.Client, clienteID string) (float64, float64, error) {
	// Busca o limite do cliente no Redis
	limiteStr, err := rdb.Get(ctx, "limite:"+clienteID).Result()
	if err != nil {
		return 0, 0, err
	}
	limite, _ := strconv.ParseFloat(limiteStr, 64)

	// Busca o saldo atual do cliente no Redis
	saldoStr, err := rdb.Get(ctx, "saldo:"+clienteID).Result()
	if err != nil {
		saldoStr = "0" // Assume saldo 0 se não encontrado
	}
	saldo, _ := strconv.ParseFloat(saldoStr, 64)

	return limite, saldo, nil
}

// Recupera as ultimas 10 transações do cliente da camada de cache
func RecuperarTransacoes(ctx context.Context, rdb *redis.Client, clienteID string) ([]entities.Transacao, error) {
	var transacoes []entities.Transacao
	transacoesCache, err := rdb.Get(ctx, "extrato:"+clienteID).Result()
	if err != nil {
		fmt.Println(err)
		return transacoes, err
	}

	err = json.Unmarshal([]byte(transacoesCache), &transacoes)
	if err != nil {
		fmt.Println(err)
		return transacoes, err
	}

	return transacoes, nil
}
