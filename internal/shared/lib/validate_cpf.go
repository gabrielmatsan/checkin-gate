package lib

import (
	"regexp"
	"strconv"
	"strings"
)

type CPFValidationResult struct {
	Valid      bool
	CleanedCPF *string
}

// Pega o cpf, retorna ele limpo e sem pontuação se for valido, se não for valido, retorna
func ValidateCPF(cpf string) CPFValidationResult {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	cpf = strings.ReplaceAll(cpf, " ", "")

	if len(cpf) != 11 {
		return CPFValidationResult{
			Valid:      false,
			CleanedCPF: nil,
		}
	}

	// Verifica se o CPF é composto por números repetidos
	var regexCPF = regexp.MustCompile(`^(\d)\1{10}$`)
	if regexCPF.MatchString(cpf) {
		return CPFValidationResult{
			Valid:      false,
			CleanedCPF: &cpf,
		}
	}

	// Converte cada caractere em dígito
	digits := make([]int, 11)
	for i, ch := range cpf {
		d, err := strconv.Atoi(string(ch))
		if err != nil {
			return CPFValidationResult{Valid: false, CleanedCPF: nil}
		}
		digits[i] = d
	}

	// Primeiro dígito verificador: pesos 10..2 sobre os 9 primeiros dígitos
	sum := 0
	for i := range 9 {
		sum += digits[i] * (10 - i)
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	// Segundo dígito verificador: pesos 11..2 sobre os 10 primeiros dígitos
	sum = 0
	for i := range 10 {
		sum += digits[i] * (11 - i)
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	if digits[9] != firstDigit || digits[10] != secondDigit {
		return CPFValidationResult{Valid: false, CleanedCPF: &cpf}
	}

	return CPFValidationResult{Valid: true, CleanedCPF: &cpf}
}
