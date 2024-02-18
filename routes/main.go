package routes

import (
	"net/http"
)

func ClientesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		TransacaoHandler(w, r)
		return
	case "GET":
		ExtratoHandler(w, r)
		return
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
}
