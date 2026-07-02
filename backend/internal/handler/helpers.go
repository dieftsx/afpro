// Package handler contém os handlers HTTP da API REST da AFPERO.
// Usa o chi router. Cada handler é responsável apenas por:
// - Decodificar a requisição
// - Chamar o service
// - Retornar a resposta JSON
package handler

import (
	"encoding/json"
	"net/http"
)

// ─── Helpers HTTP ─────────────────────────────────────────────────────────────

// respondJSON serializa o payload como JSON e escreve na resposta.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// respondError escreve uma resposta de erro padronizada.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// decodeJSON decodifica o body da requisição para o destino.
func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
