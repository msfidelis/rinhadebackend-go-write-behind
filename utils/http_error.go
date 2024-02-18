package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

// Handler simples de erros http. Escrito na mão pois não foi utilizado
// um microframework.
func HttpError(w http.ResponseWriter, message string, status int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	resError := ResponseError{
		Message: "Sem limite disponível",
	}
	if err := json.NewEncoder(w).Encode(resError); err != nil {
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
	return
}
