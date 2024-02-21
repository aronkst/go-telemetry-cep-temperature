package utils_test

import (
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
)

func TestCleanString(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"São Paulo", "SaoPaulo"},
		{"Espaço    Teste", "EspacoTeste"},
		{"Olá, Mundo!", "OlaMundo"},
		{"", ""},
		{" ", ""},
		{"Cão", "Cao"},
	}

	for _, c := range cases {
		got := utils.CleanString(c.in)
		if got != c.want {
			t.Errorf("cleanString(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
