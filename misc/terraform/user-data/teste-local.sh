#!/usr/bin/env bash

# Use este script para executar testes locais

RESULTS_WORKSPACE="/tmp/gatling/results"
GATLING_HOME="/opt/gatling/"
GATLING_BIN_DIR=/opt/gatling/bin
GATLING_WORKSPACE="/tmp/gatling"

runGatling() {
    sh $GATLING_BIN_DIR/gatling.sh -rm local -s RinhaBackendCrebitosSimulation \
        -rd "Rinha de Backend - 2024/Q1: Crébito" \
        -rf $RESULTS_WORKSPACE \
        -sf "$GATLING_WORKSPACE/simulations"
}

startTest() {
    mkdir -p $GATLING_HOME/user-files/simulations/rinhadebackend/
    cp ../../gatling/user-files/simulations/rinhadebackend/RinhaBackendCrebitosSimulation.scala $GATLING_HOME/user-files/simulations/rinhadebackend/
    for i in {1..20}; do
        # 2 requests to wake the 2 api instances up :)
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        runGatling && \
        break || sleep 2;
    done
}

startTest