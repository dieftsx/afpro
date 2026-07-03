// Package main é o entrypoint da API REST da afpro.
// Responsável por: carregar configurações, instanciar dependências (DI manual),
// montar o router e iniciar o servidor HTTP.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	supa "github.com/supabase-community/supabase-go"

	"github.com/afpro/backend/internal/adapters/assinatura"
	"github.com/afpro/backend/internal/adapters/pagamento"
	"github.com/afpro/backend/internal/config"
	"github.com/afpro/backend/internal/handler"
	"github.com/afpro/backend/internal/repository"
	"github.com/afpro/backend/internal/service"
)

func main() {
	// ── 1. Carregar configurações ──────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] Erro ao carregar configurações: %v", err)
	}

	log.Printf("[INFO] Iniciando afpro API — ambiente: %s", cfg.AppEnv)

	// ── 2. Instanciar cliente Supabase (service_role key) ─────────────────────
	supaClient, err := supa.NewClient(cfg.SupabaseURL, cfg.SupabaseServiceRoleKey, &supa.ClientOptions{})
	if err != nil {
		log.Fatalf("[FATAL] Erro ao criar cliente Supabase: %v", err)
	}

	// ── 3. Repositórios ───────────────────────────────────────────────────────
	temporadaRepo := repository.NewTemporadaRepository(supaClient)
	participRepo  := repository.NewParticipanteRepository(supaClient)
	documentoRepo := repository.NewDocumentoRepository(supaClient)
	transacaoRepo := repository.NewTransacaoRepository(supaClient)

	// ── 4. Adapters (Providers Externos) ──────────────────────────────────────
	// Pagamento: Asaas (real) — trocar por MercadoPagoAdapter para usar o outro provider
	pagamentoGW := pagamento.NewAsaasAdapter(cfg.AsaasAPIKey, cfg.AsaasBaseURL)

	// Assinatura: Clicksign (real) — trocar por ZapSignAdapter para usar o outro provider
	assinaturaGW := assinatura.NewClicksignAdapter(cfg.ClicksignAccessToken, cfg.ClicksignBaseURL)

	// ── 5. Services (Lógica de Negócio) ───────────────────────────────────────
	inscricaoSvc := service.NewInscricaoService(
		temporadaRepo,
		participRepo,
		documentoRepo,
		transacaoRepo,
		pagamentoGW,
		assinaturaGW,
	)

	financeiroSvc := service.NewFinanceiroService(
		participRepo,
		temporadaRepo,
		transacaoRepo,
		pagamentoGW,
	)

	// ── 6. Handlers HTTP ──────────────────────────────────────────────────────
	inscricaoHandler := handler.NewInscricaoHandler(
		inscricaoSvc,
		temporadaRepo,
		participRepo,
		transacaoRepo,
	)

	pagamentoWebhookHandler := handler.NewPagamentoWebhookHandler(
		cfg.AsaasWebhookToken,
		transacaoRepo,
		participRepo,
		financeiroSvc,
	)

	assinaturaWebhookHandler := handler.NewAssinaturaWebhookHandler(
		cfg.ClicksignWebhookHMAC,
		documentoRepo,
		participRepo,
	)

	// ── 7. Router ─────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// CORS — permite apenas o frontend configurado
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ── Rotas da API ──
	r.Route("/api/v1", func(r chi.Router) {
		// Endpoint público de tema (sem rate limit — alto volume esperado)
		r.Get("/temporada/ativa", inscricaoHandler.GetTemporadaAtiva)

		// Inscrições — com rate limiting por IP (anti-spam)
		r.Route("/inscricoes", func(r chi.Router) {
			r.With(httprate.LimitByIP(5, 1)).Post("/", inscricaoHandler.CriarInscricao)
			r.Get("/{id}", inscricaoHandler.GetInscricao)
			r.Post("/{id}/contrato/aceitar", inscricaoHandler.AceitarContrato)
			r.Get("/{id}/pix", inscricaoHandler.GetPIX)
		})
	})

	// ── Webhooks (recebidos dos providers externos — sem CORS) ──
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/pagamento", pagamentoWebhookHandler.HandlePagamento)
		r.Post("/assinatura", assinaturaWebhookHandler.HandleAssinatura)
	})

	// ── Health check ──
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// ── 8. Iniciar servidor ───────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("[INFO] Servidor HTTP iniciado em http://localhost%s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("[FATAL] Erro ao iniciar servidor: %v", err)
	}
}
