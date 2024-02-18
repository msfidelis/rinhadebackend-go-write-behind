package entities

type Cliente struct {
	ID     string `json:"id"`
	Saldo  string `json:"saldo"`
	Limite string `json:"limite"`
}
