// Package ports define as interfaces (contratos) que os adapters externos devem implementar.
// Seguindo o padrão de Ports & Adapters (Arquitetura Hexagonal).
package ports

import (
	"context"
	"time"

	"github.com/afpero/backend/internal/domain"
)

// ─── Port: Gateway de Pagamento ───────────────────────────────────────────────

// PIXRequest encapsula os parâmetros para gerar uma cobrança PIX.
type PIXRequest struct {
	ParticipanteID string
	CustomerID     string // ID do cliente já criado no provider
	Valor          float64
	Descricao      string
	ExpiracaoMin   int // Minutos até expiração do QR Code
}

// PIXResponse contém os dados da cobrança PIX gerada.
type PIXResponse struct {
	CobrancaID   string
	QRCodeBase64 string    // Imagem do QR Code em Base64
	CopiaCola    string    // Código "Copia e Cola" para pagamento
	Expiracao    time.Time // Timestamp exato de expiração
}

// BoletoRequest encapsula os parâmetros para gerar um boleto bancário.
type BoletoRequest struct {
	ParticipanteID string
	CustomerID     string
	Valor          float64
	Vencimento     time.Time
	Descricao      string
	NumeroParcela  int // 1–5
}

// BoletoResponse contém os dados do boleto gerado.
type BoletoResponse struct {
	CobrancaID   string
	URL          string // Link para visualização/impressão do boleto
	CodigoBarras string // Linha digitável
}

// PagamentoGateway é a interface que todos os adapters de pagamento devem implementar.
// Implementações: adapters/pagamento/asaas.go (real), adapters/pagamento/mercado_pago.go (stub)
type PagamentoGateway interface {
	// CriarCliente cadastra ou recupera o cliente no provider de pagamento.
	// Retorna o customerID do provider para uso nas cobranças subsequentes.
	CriarCliente(ctx context.Context, p domain.Participante) (customerID string, err error)

	// GerarPIX cria uma cobrança PIX com QR Code e código copia-e-cola.
	// Deve ser chamado logo após o aceite do contrato.
	GerarPIX(ctx context.Context, req PIXRequest) (PIXResponse, error)

	// GerarBoleto cria um boleto bancário com data de vencimento definida.
	// Chamado 5 vezes após a confirmação do pagamento do PIX de entrada.
	GerarBoleto(ctx context.Context, req BoletoRequest) (BoletoResponse, error)

	// CancelarCobranca cancela uma cobrança pendente no provider.
	CancelarCobranca(ctx context.Context, cobrancaID string) error
}

// ─── Port: Gateway de Assinatura Digital ──────────────────────────────────────

// DocumentoRequest encapsula os dados necessários para enviar um contrato para assinatura.
type DocumentoRequest struct {
	// Identificação interna
	ParticipanteID string
	TemporadaID    string

	// Dados do signatário
	NomeSignatario     string
	EmailSignatario    string
	WhatsAppSignatario string // Formato: 5569912345678 (DDI + DDD + número)

	// Conteúdo do documento
	// O Clicksign aceita o upload do arquivo e retorna o key para associar o signatário.
	ConteudoBase64 string // PDF gerado em Base64
	NomeArquivo    string // Ex: "contrato_5a_temporada_joana_silva.pdf"

	// Metadados
	MensagemEmail string // Mensagem customizada no e-mail enviado ao signatário
}

// DocumentoResponse contém os identificadores retornados pelo provider após o envio.
type DocumentoResponse struct {
	DocID          string // Identificador do documento no provider
	SignerID       string // Identificador do signatário no provider
	LinkAssinatura string // URL que deve ser enviada/exibida à participante
}

// AssinaturaGateway é a interface que todos os adapters de assinatura devem implementar.
// Implementações: adapters/assinatura/clicksign.go (real), adapters/assinatura/zapsign.go (stub)
type AssinaturaGateway interface {
	// EnviarParaAssinatura faz upload do documento e cadastra o signatário no provider.
	// Retorna os identificadores e o link de assinatura a ser enviado à participante.
	EnviarParaAssinatura(ctx context.Context, req DocumentoRequest) (DocumentoResponse, error)

	// BuscarStatus consulta o status atual do documento no provider.
	// Útil para polling ou reconciliação manual.
	BuscarStatus(ctx context.Context, docID string) (status string, err error)

	// CancelarDocumento cancela o processo de assinatura de um documento.
	CancelarDocumento(ctx context.Context, docID string) error
}
