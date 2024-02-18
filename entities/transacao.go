package entities

type Transacao struct {
	Valor       float64 `json:"valor"`
	Tipo        string  `json:"tipo"`
	Descricao   string  `json:"descricao"`
	RealizadaEm string  `json:"realizada_em"`
}
