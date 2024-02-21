package model

type Address struct {
	PostalCode string `json:"cep"`
	Street     string `json:"logradouro"`
	Complement string `json:"complemento"`
	District   string `json:"bairro"`
	City       string `json:"localidade"`
	State      string `json:"uf"`
}
