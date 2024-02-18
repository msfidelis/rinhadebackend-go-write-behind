package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"main/entities"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func AtualizarSaldo(ctx context.Context, rdb *redis.Client, clienteID string, valorTransacao float64, tipo string) (float64, float64, bool, error) {

	actionName := "AtualizarSaldo"
	// Busca o limite do cliente no Redis

	fmt.Printf("[%s] Recuperando Limite\n", actionName)
	limiteStr, err := rdb.Get(ctx, "limite:"+clienteID).Result()
	if err != nil {
		fmt.Printf("[%s] Erro ao recuperar o limite do cliente %s: %s\n", actionName, clienteID, err)
		return 0, 0, false, err
	}
	limite, _ := strconv.ParseFloat(limiteStr, 64)

	// Busca o saldo atual do cliente no Redis
	saldoStr, err := rdb.Get(ctx, "saldo:"+clienteID).Result()
	if err != nil {
		saldoStr = "0" // Assume saldo 0 se não encontrado
	}
	saldo, _ := strconv.ParseFloat(saldoStr, 64)

	// Verifica o tipo de transação e ajusta o saldo conforme necessário
	if tipo == "c" {
		saldo += valorTransacao // Incrementa o saldo para créditos
	} else if tipo == "d" {
		saldo -= valorTransacao // Decrementa o saldo para débitos
	} else {
		return 0, 0, false, fmt.Errorf("tipo de transação inválido: %s", tipo)
	}

	if saldo < (limite * -1) {
		return 0, 0, true, fmt.Errorf("Transação excede o limite do cliente")
	}

	// Atualiza o saldo no Redis
	err = rdb.Set(ctx, "saldo:"+clienteID, fmt.Sprintf("%f", saldo), 0).Err()
	if err != nil {
		return limite, saldo, false, err
	}

	return limite, saldo, false, nil
}

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
