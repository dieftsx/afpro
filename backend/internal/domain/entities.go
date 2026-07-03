// Package domain define as entidades puras do domínio afpro.
// Nenhum arquivo neste pacote deve importar dependências de infraestrutura.
package domain

import (
	"time"
)

// ─── Enums ────────────────────────────────────────────────────────────────────

type StatusInscricao string

const (
	StatusPendente          StatusInscricao = "pendente"
	StatusContratoEnviado   StatusInscricao = "contrato_enviado"
	StatusContratoAssinado  StatusInscricao = "contrato_assinado"
	StatusPago              StatusInscricao = "pago"
	StatusCancelado         StatusInscricao = "cancelado"
)

type TamanhoVestuario string

const (
	TamanhoPP  TamanhoVestuario = "PP"
	TamanhoP   TamanhoVestuario = "P"
	TamanhoM   TamanhoVestuario = "M"
	TamanhoG   TamanhoVestuario = "G"
	TamanhoGG  TamanhoVestuario = "GG"
	TamanhoXGG TamanhoVestuario = "XGG"
)

// ─── Entidade: Temporada ──────────────────────────────────────────────────────

// Temporada representa uma edição do evento "Meninas na Pesca".
// Contém as configurações visuais (tema) e financeiras da edição.
type Temporada struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
	Edicao int  `json:"edicao"`
	Slug string `json:"slug"`

	// Datas
	DataInicioEvento         time.Time `json:"data_inicio_evento"`
	DataFimEvento            time.Time `json:"data_fim_evento"`
	DataAberturaInscricoes   time.Time `json:"data_abertura_inscricoes"`
	DataFechamentoInscricoes time.Time `json:"data_fechamento_inscricoes"`

	// Vagas
	VagasTotal       int `json:"vagas_total"`
	VagasDisponiveis int `json:"vagas_disponiveis"`

	// Financeiro
	ValorTotal       float64 `json:"valor_total"`
	ValorEntradaPIX  float64 `json:"valor_entrada_pix"`
	ValorParcela     float64 `json:"valor_parcela"`
	QtdParcelas      int     `json:"qtd_parcelas"`

	// Tema visual (injetado no frontend via CSS variables)
	CorPrimaria    string  `json:"cor_primaria"`
	CorSecundaria  string  `json:"cor_secundaria"`
	CorAcento      string  `json:"cor_acento"`
	CorFundo       string  `json:"cor_fundo"`
	CorTexto       string  `json:"cor_texto"`
	LogoURL        *string `json:"logo_url"`
	BannerURL      *string `json:"banner_url"`

	// Template do contrato (com placeholders {{PARTICIPANTE_NOME}}, etc.)
	ContratoTemplate *string `json:"contrato_template,omitempty"`

	IsAtiva   bool      `json:"is_ativa"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ─── Entidade: Participante ────────────────────────────────────────────────────

// Participante representa uma inscrita no evento.
type Participante struct {
	ID          string `json:"id"`
	TemporadaID string `json:"temporada_id"`

	// Dados pessoais
	NomeCompleto   string     `json:"nome_completo"`
	Email          string     `json:"email"`
	WhatsApp       string     `json:"whatsapp"`
	CPFHash        string     `json:"-"`            // SHA-256, nunca serializado para o cliente
	CPFLast4       string     `json:"cpf_last4"`
	DataNascimento *time.Time `json:"data_nascimento,omitempty"`

	// Contato de emergência
	EmergenciaNome        string  `json:"emergencia_nome"`
	EmergenciaTelefone    string  `json:"emergencia_telefone"`
	EmergenciaParentesco  *string `json:"emergencia_parentesco,omitempty"`

	// Logística
	TamanhoCamisa    TamanhoVestuario `json:"tamanho_camisa"`
	TamanhoBalaclava TamanhoVestuario `json:"tamanho_balaclava"`
	TamanhoBone      TamanhoVestuario `json:"tamanho_bone"`

	// Status
	Status StatusInscricao `json:"status"`

	// Metadados LGPD
	IPOrigem          *string    `json:"-"`
	UserAgent         *string    `json:"-"`
	ConsentimentoLGPD bool       `json:"consentimento_lgpd"`
	ConsentimentoAt   *time.Time `json:"consentimento_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ─── Entidade: DocumentoAssinado ─────────────────────────────────────────────

type StatusDocumento string

const (
	DocGerado      StatusDocumento = "gerado"
	DocEnviado     StatusDocumento = "enviado"
	DocVisualizado StatusDocumento = "visualizado"
	DocAssinado    StatusDocumento = "assinado"
	DocCancelado   StatusDocumento = "cancelado"
)

type ProviderAssinatura string

const (
	ProviderClicksign ProviderAssinatura = "clicksign"
	ProviderZapSign   ProviderAssinatura = "zapsign"
)

// DocumentoAssinado representa o ciclo de vida de um contrato/termo de adesão.
type DocumentoAssinado struct {
	ID             string `json:"id"`
	ParticipanteID string `json:"participante_id"`
	TemporadaID    string `json:"temporada_id"`

	Provider          ProviderAssinatura `json:"provider"`
	ProviderDocID     *string            `json:"provider_doc_id,omitempty"`
	ProviderSignerID  *string            `json:"provider_signer_id,omitempty"`

	DocumentoURL         *string `json:"documento_url,omitempty"`
	LinkAssinatura       *string `json:"link_assinatura,omitempty"`
	DocumentoAssinadoURL *string `json:"documento_assinado_url,omitempty"`

	Status       StatusDocumento `json:"status"`
	ConteudoHTML *string         `json:"-"` // Não exposto via API pública

	EnviadoEm     *time.Time `json:"enviado_em,omitempty"`
	VisualizadoEm *time.Time `json:"visualizado_em,omitempty"`
	AssinadoEm    *time.Time `json:"assinado_em,omitempty"`

	WebhookPayload    interface{} `json:"-"` // JSONB — apenas para auditoria
	WebhookRecebidoEm *time.Time  `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ─── Entidade: Transacao ──────────────────────────────────────────────────────

type StatusTransacao string

const (
	TransacaoPendente    StatusTransacao = "pendente"
	TransacaoPaga        StatusTransacao = "pago"
	TransacaoExpirada    StatusTransacao = "expirado"
	TransacaoCancelada   StatusTransacao = "cancelado"
	TransacaoReembolsada StatusTransacao = "reembolsado"
)

type TipoTransacao string

const (
	TipoEntradaPIX    TipoTransacao = "entrada_pix"
	TipoParcelaBoleto TipoTransacao = "parcela_boleto"
)

type ProviderPagamento string

const (
	ProviderAsaas       ProviderPagamento = "asaas"
	ProviderMercadoPago ProviderPagamento = "mercado_pago"
)

// Transacao representa uma cobrança financeira individual.
// Cada inscrita gera 1 PIX de entrada (R$800) + 5 boletos (R$540 cada).
type Transacao struct {
	ID             string `json:"id"`
	ParticipanteID string `json:"participante_id"`
	TemporadaID    string `json:"temporada_id"`

	Tipo          TipoTransacao  `json:"tipo"`
	NumeroParcela *int           `json:"numero_parcela,omitempty"`
	Descricao     string         `json:"descricao"`
	Valor         float64        `json:"valor"`
	DataVencimento time.Time     `json:"data_vencimento"`

	Provider           ProviderPagamento `json:"provider"`
	ProviderCobrancaID *string           `json:"provider_cobranca_id,omitempty"`
	ProviderCustomerID *string           `json:"provider_customer_id,omitempty"`

	// Campos PIX (apenas para TipoEntradaPIX)
	PIXQRCodeBase64 *string    `json:"pix_qr_code_base64,omitempty"`
	PIXCopiaCola    *string    `json:"pix_copia_cola,omitempty"`
	PIXExpiracao    *time.Time `json:"pix_expiracao,omitempty"`

	// Campos Boleto (apenas para TipoParcelaBoleto)
	BoletoURL         *string `json:"boleto_url,omitempty"`
	BoletoCodigoBarras *string `json:"boleto_codigo_barras,omitempty"`

	Status StatusTransacao `json:"status"`
	PagoEm *time.Time      `json:"pago_em,omitempty"`

	WebhookPayload    interface{} `json:"-"`
	WebhookRecebidoEm *time.Time  `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
