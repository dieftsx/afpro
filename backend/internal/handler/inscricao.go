// Package handler — handlers de inscrição e contrato.
package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/afpro/backend/internal/domain"
	"github.com/afpro/backend/internal/repository"
	"github.com/afpro/backend/internal/service"
)

// InscricaoHandler agrupa os handlers relacionados ao fluxo de inscrição.
type InscricaoHandler struct {
	inscricaoSvc  *service.InscricaoService
	temporadaRepo *repository.TemporadaRepository
	participRepo  *repository.ParticipanteRepository
	transacaoRepo *repository.TransacaoRepository
}

// NewInscricaoHandler cria o handler com suas dependências.
func NewInscricaoHandler(
	inscricaoSvc *service.InscricaoService,
	temporadaRepo *repository.TemporadaRepository,
	participRepo *repository.ParticipanteRepository,
	transacaoRepo *repository.TransacaoRepository,
) *InscricaoHandler {
	return &InscricaoHandler{
		inscricaoSvc:  inscricaoSvc,
		temporadaRepo: temporadaRepo,
		participRepo:  participRepo,
		transacaoRepo: transacaoRepo,
	}
}

// ─── GET /api/v1/temporada/ativa ──────────────────────────────────────────────

// GetTemporadaAtiva retorna os dados públicos da temporada ativa.
// Consumido pelo Next.js para injetar o tema visual em runtime.
func (h *InscricaoHandler) GetTemporadaAtiva(w http.ResponseWriter, r *http.Request) {
	temporada, err := h.temporadaRepo.BuscarAtiva(r.Context())
	if err != nil {
		respondError(w, http.StatusNotFound, "nenhuma temporada ativa encontrada")
		return
	}

	respondJSON(w, http.StatusOK, temporada)
}

// ─── POST /api/v1/inscricoes ─────────────────────────────────────────────────

// criarInscricaoRequest é o body esperado no Step 1.
type criarInscricaoRequest struct {
	TemporadaID string `json:"temporada_id"`

	// Dados pessoais
	NomeCompleto   string  `json:"nome_completo"`
	Email          string  `json:"email"`
	WhatsApp       string  `json:"whatsapp"`
	CPF            string  `json:"cpf"` // Enviado em texto, hasheado imediatamente no service
	DataNascimento *string `json:"data_nascimento,omitempty"` // "YYYY-MM-DD"

	// Contato de emergência
	EmergenciaNome       string  `json:"emergencia_nome"`
	EmergenciaTelefone   string  `json:"emergencia_telefone"`
	EmergenciaParentesco *string `json:"emergencia_parentesco,omitempty"`

	// Logística
	TamanhoCamisa    string `json:"tamanho_camisa"`
	TamanhoBalaclava string `json:"tamanho_balaclava"`
	TamanhoBone      string `json:"tamanho_bone"`

	// LGPD
	ConsentimentoLGPD bool `json:"consentimento_lgpd"`
}

// CriarInscricao processa o Step 1: valida, salva e retorna o contrato HTML.
func (h *InscricaoHandler) CriarInscricao(w http.ResponseWriter, r *http.Request) {
	var req criarInscricaoRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	// Validações básicas
	if req.NomeCompleto == "" || req.Email == "" || req.WhatsApp == "" || req.CPF == "" {
		respondError(w, http.StatusBadRequest, "campos obrigatórios ausentes")
		return
	}

	if !req.ConsentimentoLGPD {
		respondError(w, http.StatusBadRequest, "consentimento LGPD é obrigatório")
		return
	}

	input := service.CriarInscricaoInput{
		TemporadaID:          req.TemporadaID,
		NomeCompleto:         req.NomeCompleto,
		Email:                strings.ToLower(strings.TrimSpace(req.Email)),
		WhatsApp:             req.WhatsApp,
		CPF:                  req.CPF,
		EmergenciaNome:       req.EmergenciaNome,
		EmergenciaTelefone:   req.EmergenciaTelefone,
		EmergenciaParentesco: req.EmergenciaParentesco,
		TamanhoCamisa:        domain.TamanhoVestuario(req.TamanhoCamisa),
		TamanhoBalaclava:     domain.TamanhoVestuario(req.TamanhoBalaclava),
		TamanhoBone:          domain.TamanhoVestuario(req.TamanhoBone),
	}

	// Extrair IP e User-Agent para auditoria
	ip := r.RemoteAddr
	ua := r.UserAgent()
	input.IPOrigem = &ip
	input.UserAgent = &ua

	// Processar data de nascimento (opcional)
	if req.DataNascimento != nil && *req.DataNascimento != "" {
		t, err := time.Parse("2006-01-02", *req.DataNascimento)
		if err == nil {
			input.DataNascimento = &t
		}
	}

	result, err := h.inscricaoSvc.CriarInscricao(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, result)
}

// ─── GET /api/v1/inscricoes/:id ───────────────────────────────────────────────

// GetInscricao retorna o status atual de uma inscrição.
// Usado pelo frontend para polling (ex: aguardando confirmação do PIX).
func (h *InscricaoHandler) GetInscricao(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id obrigatório")
		return
	}

	participante, err := h.participRepo.BuscarPorID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "inscrição não encontrada")
		return
	}

	// Resposta pública — não expõe CPF hash nem dados sensíveis
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":                participante.ID,
		"nome_completo":     participante.NomeCompleto,
		"email":             participante.Email,
		"cpf_last4":         participante.CPFLast4,
		"status":            participante.Status,
		"temporada_id":      participante.TemporadaID,
		"created_at":        participante.CreatedAt,
	})
}

// ─── POST /api/v1/inscricoes/:id/contrato/aceitar ─────────────────────────────

// aceitarContratoRequest é o body esperado no Step 2.
type aceitarContratoRequest struct {
	// O CPF completo é necessário apenas aqui para criar o cliente no Asaas.
	// Não é armazenado; só trafega em HTTPS.
	CPF string `json:"cpf"`
}

// AceitarContrato processa o Step 2: aceite do contrato → gera PIX de entrada.
func (h *InscricaoHandler) AceitarContrato(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id obrigatório")
		return
	}

	var req aceitarContratoRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "payload inválido")
		return
	}

	result, err := h.inscricaoSvc.AceitarContrato(r.Context(), service.AceitarContratoInput{
		ParticipanteID: id,
		CPFCompleto:    req.CPF,
	})
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// ─── GET /api/v1/inscricoes/:id/pix ──────────────────────────────────────────

// GetPIX retorna os dados do PIX de entrada de uma inscrição.
// Usado pelo Step 3 para exibir QR Code e timer.
func (h *InscricaoHandler) GetPIX(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id obrigatório")
		return
	}

	tx, err := h.transacaoRepo.BuscarEntradaPIXPorParticipante(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "cobrança PIX não encontrada para esta inscrição")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"transacao_id":       tx.ID,
		"status":             tx.Status,
		"valor":              tx.Valor,
		"pix_qr_code_base64": tx.PIXQRCodeBase64,
		"pix_copia_cola":     tx.PIXCopiaCola,
		"pix_expiracao":      tx.PIXExpiracao,
	})
}
