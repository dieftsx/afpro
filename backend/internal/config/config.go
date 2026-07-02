// Package config carrega e valida as variáveis de ambiente da aplicação.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config agrupa todas as configurações da aplicação.
type Config struct {
	AppEnv      string
	ServerPort  string
	FrontendURL string

	// Supabase
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string

	// Asaas (gateway de pagamento — adapter real)
	AsaasAPIKey       string
	AsaasBaseURL      string
	AsaasWebhookToken string

	// Clicksign (assinatura digital — adapter real)
	ClicksignAccessToken  string
	ClicksignBaseURL      string
	ClicksignWebhookHMAC  string
}

// Load lê o arquivo .env (se existir) e popula a struct Config.
// Retorna erro se alguma variável obrigatória estiver ausente.
func Load() (*Config, error) {
	// Ignora erro quando .env não existe (ex.: ambiente de produção com vars já injetadas)
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:        os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),

		AsaasAPIKey:       os.Getenv("ASAAS_API_KEY"),
		AsaasBaseURL:      getEnv("ASAAS_BASE_URL", "https://sandbox.asaas.com/api/v3"),
		AsaasWebhookToken: os.Getenv("ASAAS_WEBHOOK_TOKEN"),

		ClicksignAccessToken: os.Getenv("CLICKSIGN_ACCESS_TOKEN"),
		ClicksignBaseURL:     getEnv("CLICKSIGN_BASE_URL", "https://sandbox.clicksign.com"),
		ClicksignWebhookHMAC: os.Getenv("CLICKSIGN_WEBHOOK_HMAC_KEY"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate verifica se todas as variáveis obrigatórias estão presentes.
func (c *Config) validate() error {
	required := map[string]string{
		"SUPABASE_URL":              c.SupabaseURL,
		"SUPABASE_ANON_KEY":         c.SupabaseAnonKey,
		"SUPABASE_SERVICE_ROLE_KEY": c.SupabaseServiceRoleKey,
		"ASAAS_API_KEY":             c.AsaasAPIKey,
		"ASAAS_WEBHOOK_TOKEN":       c.AsaasWebhookToken,
		"CLICKSIGN_ACCESS_TOKEN":    c.ClicksignAccessToken,
		"CLICKSIGN_WEBHOOK_HMAC_KEY": c.ClicksignWebhookHMAC,
	}

	for key, val := range required {
		if val == "" {
			return fmt.Errorf("variável de ambiente obrigatória ausente: %s", key)
		}
	}

	return nil
}

// IsProduction retorna true quando APP_ENV=production.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// getEnv retorna o valor da env ou o fallback quando a variável está vazia.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
