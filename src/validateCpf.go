package src

import (
	"regexp"
	"strconv"
)

const CPFLENGTH = 11

func ValidateCpf(cpf string) bool {
	if cpf == "" {
		return false
	}
	cpf = Clean(cpf)
	if len(cpf) != CPFLENGTH {
		return false
	}
	if AllDigitsAreTheSame(cpf) {
		return false
	}
	digit1 := CalculateDigit(cpf, 10)
	digit2 := CalculateDigit(cpf, 11)

	return ExtractDigit(cpf) == strconv.Itoa(digit1)+strconv.Itoa(digit2)
}

func Clean(cpf string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(cpf, "")
}

func AllDigitsAreTheSame(cpf string) bool {
	firstDigit := string(cpf[0])
	for _, digit := range cpf {
		if string(digit) != firstDigit {
			return false
		}
	}
	return true
}

func CalculateDigit(cpf string, factor int) int {
	total := 0
	for _, digit := range cpf {
		if factor > 1 {
			val, _ := strconv.Atoi(string(digit))
			total += val * factor
			factor--
		}
	}
	rest := total % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}

func ExtractDigit(cpf string) string {
	return cpf[len(cpf)-2:]
}
