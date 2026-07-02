// Package service implementa a lógica de negócio da inscrição.
// Orquestra os repositórios e os adapters externos seguindo as regras de negócio da AFPERO.
package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/afpero/backend/internal/adapters/assinatura"
	"github.com/afpero/backend/internal/domain"
	"github.com/afpero/backend/internal/ports"
	"github.com/afpero/backend/internal/repository"
	"github.com/afpero/backend/pkg/cpf"
)

// ─── DTOs de entrada ──────────────────────────────────────────────────────────

// CriarInscricaoInput representa os dados do Step 1 do formulário.
type CriarInscricaoInput struct {
	TemporadaID string

	// Dados pessoais
	NomeCompleto   string
	Email          string
	WhatsApp       string
	CPF            string // CPF completo (apenas neste ponto — hasheado imediatamente)
	DataNascimento *time.Time

	// Contato de emergência
	EmergenciaNome       string
	EmergenciaTelefone   string
	EmergenciaParentesco *string

	// Logística
	TamanhoCamisa    domain.TamanhoVestuario
	TamanhoBalaclava domain.TamanhoVestuario
	TamanhoBone      domain.TamanhoVestuario

	// Metadados LGPD
	IPOrigem  *string
	UserAgent *string
}

// AceitarContratoInput representa o aceite do Step 2.
type AceitarContratoInput struct {
	ParticipanteID string
	CPFCompleto    string // Necessário apenas para criar o cliente no Asaas
}

// ─── Service de Inscrição ─────────────────────────────────────────────────────

// InscricaoService orquestra o fluxo de inscrição.
type InscricaoService struct {
	temporadaRepo  *repository.TemporadaRepository
	participRepo   *repository.ParticipanteRepository
	documentoRepo  *repository.DocumentoRepository
	transacaoRepo  *repository.TransacaoRepository
	pagamento      ports.PagamentoGateway
	assinaturaGW   ports.AssinaturaGateway
}

// NewInscricaoService cria o service com suas dependências injetadas.
func NewInscricaoService(
	temporadaRepo *repository.TemporadaRepository,
	participRepo *repository.ParticipanteRepository,
	documentoRepo *repository.DocumentoRepository,
	transacaoRepo *repository.TransacaoRepository,
	pagamento ports.PagamentoGateway,
	assinaturaGW ports.AssinaturaGateway,
) *InscricaoService {
	return &InscricaoService{
		temporadaRepo: temporadaRepo,
		participRepo:  participRepo,
		documentoRepo: documentoRepo,
		transacaoRepo: transacaoRepo,
		pagamento:     pagamento,
		assinaturaGW:  assinaturaGW,
	}
}

// ─── Step 1: Criar Inscrição ──────────────────────────────────────────────────

// CriarInscricaoResult é o retorno do Step 1.
type CriarInscricaoResult struct {
	ParticipanteID string `json:"participante_id"`
	ContratoHTML   string `json:"contrato_html"` // Contrato preenchido para exibição no Step 2
}

// CriarInscricao executa o Step 1 do fluxo:
// valida o CPF, verifica duplicata, cria a participante e renderiza o contrato.
func (s *InscricaoService) CriarInscricao(ctx context.Context, input CriarInscricaoInput) (*CriarInscricaoResult, error) {
	// 1. Validar CPF
	if !cpf.IsValid(input.CPF) {
		return nil, fmt.Errorf("CPF inválido")
	}

	cpfHash := cpf.Hash(input.CPF)
	cpfLast4 := cpf.Last4(input.CPF)

	// 2. Verificar se inscrições estão abertas e há vagas
	temporada, err := s.temporadaRepo.BuscarPorID(ctx, input.TemporadaID)
	if err != nil {
		return nil, fmt.Errorf("temporada não encontrada: %w", err)
	}

	now := time.Now()
	if now.Before(temporada.DataAberturaInscricoes) {
		return nil, fmt.Errorf("inscrições ainda não estão abertas")
	}
	if now.After(temporada.DataFechamentoInscricoes) {
		return nil, fmt.Errorf("prazo de inscrições encerrado")
	}
	if temporada.VagasDisponiveis <= 0 {
		return nil, fmt.Errorf("não há vagas disponíveis nesta temporada")
	}

	// 3. Verificar duplicata de CPF na mesma temporada
	existe, err := s.participRepo.ExistePorCPFHashETemporada(ctx, cpfHash, input.TemporadaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if existe {
		return nil, fmt.Errorf("já existe uma inscrição com este CPF nesta temporada")
	}

	// 4. Criar participante (status=pendente)
	agora := time.Now()
	p := domain.Participante{
		TemporadaID:          input.TemporadaID,
		NomeCompleto:         input.NomeCompleto,
		Email:                input.Email,
		WhatsApp:             input.WhatsApp,
		CPFHash:              cpfHash,
		CPFLast4:             cpfLast4,
		DataNascimento:       input.DataNascimento,
		EmergenciaNome:       input.EmergenciaNome,
		EmergenciaTelefone:   input.EmergenciaTelefone,
		EmergenciaParentesco: input.EmergenciaParentesco,
		TamanhoCamisa:        input.TamanhoCamisa,
		TamanhoBalaclava:     input.TamanhoBalaclava,
		TamanhoBone:          input.TamanhoBone,
		Status:               domain.StatusPendente,
		IPOrigem:             input.IPOrigem,
		UserAgent:            input.UserAgent,
		ConsentimentoLGPD:    true,
		ConsentimentoAt:      &agora,
	}

	participante, err := s.participRepo.Criar(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar participante: %w", err)
	}

	// 5. Renderizar o contrato com os dados reais
	contratoHTML, err := s.renderizarContrato(ctx, participante, temporada, input.CPF)
	if err != nil {
		return nil, fmt.Errorf("erro ao renderizar contrato: %w", err)
	}

	return &CriarInscricaoResult{
		ParticipanteID: participante.ID,
		ContratoHTML:   contratoHTML,
	}, nil
}

// ─── Step 2: Aceitar Contrato ─────────────────────────────────────────────────

// AceitarContratoResult é o retorno do Step 2 (dados para o checkout PIX).
type AceitarContratoResult struct {
	TransacaoID   string    `json:"transacao_id"`
	PIXQRCode     string    `json:"pix_qr_code_base64"`
	PIXCopiaCola  string    `json:"pix_copia_cola"`
	PIXExpiracao  time.Time `json:"pix_expiracao"`
	ValorEntrada  float64   `json:"valor_entrada"`
}

// AceitarContrato executa o Step 2:
// 1. Cria o documento no Clicksign
// 2. Atualiza status da inscrição para contrato_enviado
// 3. Gera o PIX de entrada no Asaas
func (s *InscricaoService) AceitarContrato(ctx context.Context, input AceitarContratoInput) (*AceitarContratoResult, error) {
	// 1. Buscar participante
	participante, err := s.participRepo.BuscarPorID(ctx, input.ParticipanteID)
	if err != nil {
		return nil, fmt.Errorf("participante não encontrado: %w", err)
	}

	if participante.Status != domain.StatusPendente {
		return nil, fmt.Errorf("inscrição já processada (status: %s)", participante.Status)
	}

	// 2. Buscar temporada para obter o template e os valores
	temporada, err := s.temporadaRepo.BuscarPorID(ctx, participante.TemporadaID)
	if err != nil {
		return nil, fmt.Errorf("temporada não encontrada: %w", err)
	}

	// 3. Renderizar e enviar contrato ao Clicksign
	contratoHTML, err := s.renderizarContrato(ctx, participante, temporada, input.CPFCompleto)
	if err != nil {
		return nil, fmt.Errorf("erro ao renderizar contrato: %w", err)
	}

	// 4. Salvar documento localmente (status=gerado)
	doc, err := s.documentoRepo.Criar(ctx, domain.DocumentoAssinado{
		ParticipanteID: participante.ID,
		TemporadaID:    participante.TemporadaID,
		Provider:       domain.ProviderClicksign,
		Status:         domain.DocGerado,
		ConteudoHTML:   &contratoHTML,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar documento: %w", err)
	}

	// 5. Enviar ao Clicksign
	docBase64 := assinatura.HTMLParaBase64(contratoHTML)
	nomeArquivo := fmt.Sprintf("contrato_%s_%s.html",
		strings.ToLower(strings.ReplaceAll(participante.NomeCompleto, " ", "_")),
		participante.ID[:8],
	)

	docResp, err := s.assinaturaGW.EnviarParaAssinatura(ctx, ports.DocumentoRequest{
		ParticipanteID:     participante.ID,
		TemporadaID:        participante.TemporadaID,
		NomeSignatario:     participante.NomeCompleto,
		EmailSignatario:    participante.Email,
		WhatsAppSignatario: participante.WhatsApp,
		ConteudoBase64:     docBase64,
		NomeArquivo:        nomeArquivo,
		MensagemEmail: fmt.Sprintf(
			"Olá %s! Acesse o link abaixo para assinar seu contrato de inscrição na %s.",
			participante.NomeCompleto, temporada.Nome,
		),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar ao Clicksign: %w", err)
	}

	// 6. Atualizar documento com dados do Clicksign
	if err := s.documentoRepo.AtualizarAposEnvio(ctx,
		doc.ID,
		docResp.DocID,
		docResp.SignerID,
		docResp.LinkAssinatura,
	); err != nil {
		return nil, fmt.Errorf("erro ao atualizar documento: %w", err)
	}

	// 7. Atualizar status da inscrição
	if err := s.participRepo.AtualizarStatus(ctx, participante.ID, domain.StatusContratoEnviado); err != nil {
		return nil, fmt.Errorf("erro ao atualizar status da inscrição: %w", err)
	}

	// 8. Criar cliente no Asaas (necessita do CPF completo)
	customerID, err := s.pagamento.CriarCliente(ctx, *participante)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente no Asaas: %w", err)
	}

	// 9. Gerar PIX de entrada (R$800, expira em 30 minutos)
	pixResp, err := s.pagamento.GerarPIX(ctx, ports.PIXRequest{
		ParticipanteID: participante.ID,
		CustomerID:     customerID,
		Valor:          temporada.ValorEntradaPIX,
		Descricao:      fmt.Sprintf("Entrada — %s", temporada.Nome),
		ExpiracaoMin:   30,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar PIX: %w", err)
	}

	// 10. Salvar transação de entrada PIX
	tx, err := s.transacaoRepo.Criar(ctx, domain.Transacao{
		ParticipanteID:     participante.ID,
		TemporadaID:        participante.TemporadaID,
		Tipo:               domain.TipoEntradaPIX,
		Descricao:          fmt.Sprintf("Entrada PIX — %s", temporada.Nome),
		Valor:              temporada.ValorEntradaPIX,
		DataVencimento:     time.Now().AddDate(0, 0, 1),
		Provider:           domain.ProviderAsaas,
		ProviderCobrancaID: &pixResp.CobrancaID,
		ProviderCustomerID: &customerID,
		PIXQRCodeBase64:    &pixResp.QRCodeBase64,
		PIXCopiaCola:       &pixResp.CopiaCola,
		PIXExpiracao:       &pixResp.Expiracao,
		Status:             domain.TransacaoPendente,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar transação PIX: %w", err)
	}

	return &AceitarContratoResult{
		TransacaoID:  tx.ID,
		PIXQRCode:    pixResp.QRCodeBase64,
		PIXCopiaCola: pixResp.CopiaCola,
		PIXExpiracao: pixResp.Expiracao,
		ValorEntrada: temporada.ValorEntradaPIX,
	}, nil
}

// ─── Renderização do Contrato ─────────────────────────────────────────────────

// placeholders mapeia os tokens do template às informações reais.
type contratoData struct {
	TemporadaNome        string
	ParticipanteNome     string
	CPFLast4             string
	ParticipanteEmail    string
	ParticipanteWhatsApp string
	DataInicio           string
	DataFim              string
	ValorTotal           string
	ValorEntrada         string
	DataVencimentoEntrada string
	QtdParcelas          string
	ValorParcela         string
	DataAssinatura       string
}

// renderizarContrato substitui os placeholders do template e retorna o HTML final.
func (s *InscricaoService) renderizarContrato(
	ctx context.Context,
	p *domain.Participante,
	t *domain.Temporada,
	_ string, // CPF completo — reservado para uso futuro (ex: parceiro de nota fiscal)
) (string, error) {
	if t.ContratoTemplate == nil || *t.ContratoTemplate == "" {
		return "", fmt.Errorf("template de contrato não configurado para esta temporada")
	}

	// Converte placeholders {{...}} para sintaxe Go template {{.FieldName}}
	tmplText := *t.ContratoTemplate
	replacer := strings.NewReplacer(
		"{{TEMPORADA_NOME}}",          "{{.TemporadaNome}}",
		"{{PARTICIPANTE_NOME}}",       "{{.ParticipanteNome}}",
		"{{CPF_LAST4}}",               "{{.CPFLast4}}",
		"{{PARTICIPANTE_EMAIL}}",      "{{.ParticipanteEmail}}",
		"{{PARTICIPANTE_WHATSAPP}}",   "{{.ParticipanteWhatsApp}}",
		"{{DATA_INICIO}}",             "{{.DataInicio}}",
		"{{DATA_FIM}}",                "{{.DataFim}}",
		"{{VALOR_TOTAL}}",             "{{.ValorTotal}}",
		"{{VALOR_ENTRADA}}",           "{{.ValorEntrada}}",
		"{{DATA_VENCIMENTO_ENTRADA}}", "{{.DataVencimentoEntrada}}",
		"{{QTD_PARCELAS}}",            "{{.QtdParcelas}}",
		"{{VALOR_PARCELA}}",           "{{.ValorParcela}}",
		"{{DATA_ASSINATURA}}",         "{{.DataAssinatura}}",
	)
	tmplText = replacer.Replace(tmplText)

	tmpl, err := template.New("contrato").Parse(tmplText)
	if err != nil {
		return "", fmt.Errorf("erro ao parsear template: %w", err)
	}

	data := contratoData{
		TemporadaNome:         t.Nome,
		ParticipanteNome:      p.NomeCompleto,
		CPFLast4:              p.CPFLast4,
		ParticipanteEmail:     p.Email,
		ParticipanteWhatsApp:  p.WhatsApp,
		DataInicio:            t.DataInicioEvento.Format("02/01/2006"),
		DataFim:               t.DataFimEvento.Format("02/01/2006"),
		ValorTotal:            fmt.Sprintf("%.2f", t.ValorTotal),
		ValorEntrada:          fmt.Sprintf("%.2f", t.ValorEntradaPIX),
		DataVencimentoEntrada: time.Now().AddDate(0, 0, 1).Format("02/01/2006"),
		QtdParcelas:           fmt.Sprintf("%d", t.QtdParcelas),
		ValorParcela:          fmt.Sprintf("%.2f", t.ValorParcela),
		DataAssinatura:        time.Now().Format("02/01/2006"),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("erro ao executar template: %w", err)
	}

	return buf.String(), nil
}
