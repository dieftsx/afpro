// Package cpf fornece utilitários para validação e hashing de CPF.
// O CPF nunca é armazenado em texto claro — apenas o hash SHA-256.
package cpf

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var digitsOnly = regexp.MustCompile(`\D`)

// Sanitize remove máscara do CPF, retornando apenas os 11 dígitos.
// Ex: "123.456.789-09" → "12345678909"
func Sanitize(cpf string) string {
	return digitsOnly.ReplaceAllString(cpf, "")
}

// IsValid valida um CPF pelo algoritmo dos dígitos verificadores.
// Aceita CPFs com ou sem máscara.
func IsValid(raw string) bool {
	c := Sanitize(raw)

	if len(c) != 11 {
		return false
	}

	// Rejeita sequências inválidas conhecidas (ex: 111.111.111-11)
	if strings.Count(c, string(c[0])) == 11 {
		return false
	}

	// Valida 1º dígito verificador
	if !validaDigito(c, 9) {
		return false
	}

	// Valida 2º dígito verificador
	return validaDigito(c, 10)
}

// validaDigito calcula e confere o dígito verificador na posição `pos`.
func validaDigito(cpf string, pos int) bool {
	sum := 0
	for i := 0; i < pos; i++ {
		d, _ := strconv.Atoi(string(cpf[i]))
		sum += d * (pos + 1 - i)
	}

	remainder := (sum * 10) % 11
	if remainder == 10 || remainder == 11 {
		remainder = 0
	}

	digit, _ := strconv.Atoi(string(cpf[pos]))
	return remainder == digit
}

// Hash retorna o SHA-256 hexadecimal do CPF sanitizado.
// É o valor armazenado na coluna cpf_hash da tabela participantes.
func Hash(raw string) string {
	clean := Sanitize(raw)
	h := sha256.Sum256([]byte(clean))
	return fmt.Sprintf("%x", h)
}

// Last4 retorna os últimos 4 dígitos do CPF.
// Usado para exibição sem expor o CPF completo.
func Last4(raw string) string {
	clean := Sanitize(raw)
	if len(clean) < 4 {
		return clean
	}
	return clean[len(clean)-4:]
}
