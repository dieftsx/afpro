// Package handler — webhooks de pagamento (Asaas) e assinatura (Clicksign).
package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/afpro/backend/internal/domain"
	"github.com/afpro/backend/internal/repository"
	"github.com/afpro/backend/internal/service"
)

// ─── Webhook de Pagamento (Asaas) ─────────────────────────────────────────────

// PagamentoWebhookHandler processa confirmações de pagamento enviadas pelo Asaas.
type PagamentoWebhookHandler struct {
	webhookToken  string // Valor configurado no painel Asaas e na env ASAAS_WEBHOOK_TOKEN
	transacaoRepo *repository.TransacaoRepository
	participRepo  *repository.ParticipanteRepository
	financeiroSvc *service.FinanceiroService
}

// NewPagamentoWebhookHandler cria o handler.
func NewPagamentoWebhookHandler(
	webhookToken string,
	transacaoRepo *repository.TransacaoRepository,
	participRepo *repository.ParticipanteRepository,
	financeiroSvc *service.FinanceiroService,
) *PagamentoWebhookHandler {
	return &PagamentoWebhookHandler{
		webhookToken:  webhookToken,
		transacaoRepo: transacaoRepo,
		participRepo:  participRepo,
		financeiroSvc: financeiroSvc,
	}
}

// Estrutura do evento enviado pelo Asaas
type asaasWebhookEvent struct {
	Event   string `json:"event"` // "PAYMENT_RECEIVED" | "PAYMENT_OVERDUE" | etc.
	Payment struct {
		ID          string  `json:"id"`
		Status      string  `json:"status"`      // "RECEIVED" | "CONFIRMED" | "OVERDUE"
		BillingType string  `json:"billingType"` // "PIX" | "BOLETO"
		Value       float64 `json:"value"`
		ExternalReference string `json:"externalReference"` // participante_id
	} `json:"payment"`
}

// HandlePagamento processa os webhooks do Asaas.
// POST /webhooks/pagamento
func (h *PagamentoWebhookHandler) HandlePagamento(w http.ResponseWriter, r *http.Request) {
	// 1. Validar autenticação via header (Asaas envia access_token no header)
	token := r.Header.Get("asaas-access-token")
	if token != h.webhookToken {
		http.Error(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	// 2. Ler body e preservar para auditoria
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "erro ao ler body")
		return
	}
	defer r.Body.Close()

	var event asaasWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		respondError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	// 3. Processar apenas eventos de pagamento confirmado
	if event.Event != "PAYMENT_RECEIVED" && event.Event != "PAYMENT_CONFIRMED" {
		// Outros eventos (OVERDUE, REFUNDED, etc.) — aceitar silenciosamente
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	// 4. Buscar a transação pelo ID da cobrança no provider
	transacao, err := h.transacaoRepo.BuscarPorProviderCobrancaID(ctx, event.Payment.ID)
	if err != nil {
		// Se não encontrarmos, pode ser uma cobrança de outra temporada/ambiente
		respondError(w, http.StatusNotFound, "transação não encontrada")
		return
	}

	// Já processada — responde OK para evitar reentrada
	if transacao.Status == domain.TransacaoPaga {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 5. Marcar transação como paga
	rawPayload := json.RawMessage(body)
	if err := h.transacaoRepo.MarcarComoPago(ctx, transacao.ID, rawPayload); err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao atualizar transação")
		return
	}

	// 6. Se for o PIX de entrada, desencadear a geração dos boletos
	if transacao.Tipo == domain.TipoEntradaPIX {
		participante, err := h.participRepo.BuscarPorID(ctx, transacao.ParticipanteID)
		if err != nil {
			fmt.Printf("[ERRO] webhook pagamento: participante não encontrado %s\n", transacao.ParticipanteID)
			// Não retorna erro para o Asaas — já marcamos o pagamento
			w.WriteHeader(http.StatusOK)
			return
		}

		// Atualizar status da inscrita para "pago"
		if err := h.participRepo.AtualizarStatus(ctx, participante.ID, domain.StatusPago); err != nil {
			fmt.Printf("[ERRO] webhook pagamento: atualizar status participante %v\n", err)
		}

		// Gerar os 5 boletos mensais (operação assíncrona em produção — use uma fila)
		customerID := ""
		if transacao.ProviderCustomerID != nil {
			customerID = *transacao.ProviderCustomerID
		}

		if err := h.financeiroSvc.GerarParcelas(ctx, *participante, customerID); err != nil {
			fmt.Printf("[ERRO] webhook pagamento: gerar parcelas %v\n", err)
			// Não retorna erro — o pagamento já foi registrado. Retry manual via reconciliação.
		}
	}

	w.WriteHeader(http.StatusOK)
}

// ─── Webhook de Assinatura (Clicksign) ────────────────────────────────────────

// AssinaturaWebhookHandler processa notificações de assinatura enviadas pelo Clicksign.
type AssinaturaWebhookHandler struct {
	hmacSecret    string // Segredo HMAC configurado no painel Clicksign
	documentoRepo *repository.DocumentoRepository
	participRepo  *repository.ParticipanteRepository
}

// NewAssinaturaWebhookHandler cria o handler.
func NewAssinaturaWebhookHandler(
	hmacSecret string,
	documentoRepo *repository.DocumentoRepository,
	participRepo *repository.ParticipanteRepository,
) *AssinaturaWebhookHandler {
	return &AssinaturaWebhookHandler{
		hmacSecret:    hmacSecret,
		documentoRepo: documentoRepo,
		participRepo:  participRepo,
	}
}

// Estrutura do evento enviado pelo Clicksign
type clicksignWebhookEvent struct {
	Event struct {
		Name string `json:"name"` // "sign" | "auto_close" | "cancel"
		Data struct {
			Document struct {
				Key    string `json:"key"`
				Status string `json:"status"` // "closed" = todos assinaram
			} `json:"document"`
		} `json:"data"`
	} `json:"event"`
}

// HandleAssinatura processa os webhooks do Clicksign.
// POST /webhooks/assinatura
func (h *AssinaturaWebhookHandler) HandleAssinatura(w http.ResponseWriter, r *http.Request) {
	// 1. Ler body antes de validar (precisamos do body para o HMAC)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "erro ao ler body")
		return
	}
	defer r.Body.Close()

	// 2. Validar HMAC-SHA256 (Clicksign envia no header X-Clicksign-Hmac-SHA256)
	if !h.validarHMAC(body, r.Header.Get("X-Clicksign-Hmac-SHA256")) {
		http.Error(w, "assinatura HMAC inválida", http.StatusUnauthorized)
		return
	}

	var event clicksignWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		respondError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	// 3. Processar apenas eventos de documento fechado (todos assinaram)
	if event.Event.Name != "auto_close" && event.Event.Data.Document.Status != "closed" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	docKey := event.Event.Data.Document.Key

	// 4. Buscar documento pelo ID do provider
	documento, err := h.documentoRepo.BuscarPorProviderDocID(ctx, docKey)
	if err != nil {
		respondError(w, http.StatusNotFound, "documento não encontrado")
		return
	}

	// Já processado
	if documento.Status == domain.DocAssinado {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 5. Atualizar documento para "assinado"
	rawPayload := json.RawMessage(body)
	if err := h.documentoRepo.AtualizarAposAssinatura(ctx, docKey, rawPayload); err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao atualizar documento")
		return
	}

	// 6. Atualizar status da inscrição para "contrato_assinado"
	if err := h.participRepo.AtualizarStatus(ctx, documento.ParticipanteID, domain.StatusContratoAssinado); err != nil {
		fmt.Printf("[ERRO] webhook assinatura: atualizar status %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
}

// validarHMAC verifica a autenticidade do webhook do Clicksign.
func (h *AssinaturaWebhookHandler) validarHMAC(body []byte, headerHMAC string) bool {
	if h.hmacSecret == "" || headerHMAC == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(h.hmacSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(headerHMAC))
}
