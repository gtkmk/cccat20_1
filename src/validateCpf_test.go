package src

import "testing"

func TestValidateCpf_Valid(t *testing.T) {
	tests := []string{
		"97456321558",
		"71428793860",
		"974.563.215-58",
		"714.287.938-60",
	}

	for _, cpf := range tests {
		if !ValidateCpf(cpf) {
			t.Errorf("ValidateCpf(%s) = false; want true", cpf)
		}
	}
}

func TestValidateCpf_Invalid(t *testing.T) {
	tests := []string{
		"",
		"11111111111",
	}

	for _, cpf := range tests {
		if ValidateCpf(cpf) {
			t.Errorf("ValidateCpf(%s) = true; want false", cpf)
		}
	}
}
