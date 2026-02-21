package lib

import (
	"regexp"
	"strings"
)

type BrazilPostalCodeValidationResult struct {
	Valid             bool
	CleanedPostalCode *string
}

var onlyNumbers = regexp.MustCompile(`^[0-9]+$`)

func ValidateBrazilPostalCode(postalCode string) BrazilPostalCodeValidationResult {
	postalCode = strings.TrimSpace(postalCode)
	postalCode = strings.ReplaceAll(postalCode, "-", "")

	// Verifica se o CEP tem 8 dígitos
	if len(postalCode) != 8 {
		return BrazilPostalCodeValidationResult{
			Valid:             false,
			CleanedPostalCode: nil,
		}
	}

	// Verifica se o CEP é composto apenas por números
	if !onlyNumbers.MatchString(postalCode) {
		return BrazilPostalCodeValidationResult{
			Valid:             false,
			CleanedPostalCode: nil,
		}
	}

	// Verifica se o CEP é composto por números repetidos
	if isAllSameDigit(postalCode) {
		return BrazilPostalCodeValidationResult{
			Valid:             false,
			CleanedPostalCode: nil,
		}
	}

	return BrazilPostalCodeValidationResult{
		Valid:             true,
		CleanedPostalCode: &postalCode,
	}
}
