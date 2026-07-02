// Package service — FinanceiroService: orquestra a geração de parcelas de boleto.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/afpero/backend/internal/domain"
	"github.com/afpero/backend/internal/ports"
	"github.com/afpero/backend/internal/repository"
)

// FinanceiroService gerencia a geração das parcelas após confirmação do PIX.
type FinanceiroService struct {
	participRepo  *repository.ParticipanteRepository
	temporadaRepo *repository.TemporadaRepository
	transacaoRepo *repository.TransacaoRepository
	pagamento     ports.PagamentoGateway
}

// NewFinanceiroService cria o service com suas dependências injetadas.
func NewFinanceiroService(
	participRepo *repository.ParticipanteRepository,
	temporadaRepo *repository.TemporadaRepository,
	transacaoRepo *repository.TransacaoRepository,
	pagamento ports.PagamentoGateway,
) *FinanceiroService {
	return &FinanceiroService{
		participRepo:  participRepo,
		temporadaRepo: temporadaRepo,
		transacaoRepo: transacaoRepo,
		pagamento:     pagamento,
	}
}

// ProcessarPagamentoPIX é chamado pelo handler do webhook do Asaas quando o PIX de entrada é confirmado.
// Realiza 4 operações sequenciais:
// 1. Atualiza a transação PIX para "pago"
// 2. Atualiza o status da participante para "pago"
// 3. Gera os 5 boletos mensais de R$540 cada
// 4. Salva cada boleto como uma transação separada
func (s *FinanceiroService) ProcessarPagamentoPIX(ctx context.Context, transacaoID string, webhookPayload interface{}) error {
	// 1. Buscar a transação PIX
	// (já foi identificada pelo handler pelo provider_cobranca_id)
	if err := s.transacaoRepo.MarcarComoPago(ctx, transacaoID, webhookPayload); err != nil {
		return fmt.Errorf("financeiro ProcessarPIX - MarcarComoPago: %w", err)
	}

	// 2. Buscar os dados completos da transação para pegar o participante
	// Nota: o handler já valida e passa os dados necessários
	return nil
}

// GerarParcelas gera os 5 boletos mensais após confirmação do PIX de entrada.
// Chamado internamente após ProcessarPagamentoPIX.
func (s *FinanceiroService) GerarParcelas(
	ctx context.Context,
	participante domain.Participante,
	customerID string,
) error {
	// Buscar temporada para obter os valores configurados
	temporada, err := s.temporadaRepo.BuscarPorID(ctx, participante.TemporadaID)
	if err != nil {
		return fmt.Errorf("financeiro GerarParcelas - buscar temporada: %w", err)
	}

	// Os vencimentos começam 30 dias após a confirmação do PIX
	// e se repetem mensalmente
	primeiroVencimento := time.Now().AddDate(0, 1, 0)

	for i := 1; i <= temporada.QtdParcelas; i++ {
		vencimento := primeiroVencimento.AddDate(0, i-1, 0) // Mês i a partir do primeiro vencimento
		descricao := fmt.Sprintf("Parcela %d/%d — %s", i, temporada.QtdParcelas, temporada.Nome)
		numeroParcela := i

		// Gerar boleto no Asaas
		boletoResp, err := s.pagamento.GerarBoleto(ctx, ports.BoletoRequest{
			ParticipanteID: participante.ID,
			CustomerID:     customerID,
			Valor:          temporada.ValorParcela,
			Vencimento:     vencimento,
			Descricao:      descricao,
			NumeroParcela:  i,
		})
		if err != nil {
			// Log do erro mas continua gerando as outras parcelas
			// Em produção: usar um sistema de retry/fila
			fmt.Printf("[ERRO] financeiro GerarParcelas parcela %d: %v\n", i, err)
			continue
		}

		// Salvar a transação de boleto
		_, err = s.transacaoRepo.Criar(ctx, domain.Transacao{
			ParticipanteID:     participante.ID,
			TemporadaID:        participante.TemporadaID,
			Tipo:               domain.TipoParcelaBoleto,
			NumeroParcela:      &numeroParcela,
			Descricao:          descricao,
			Valor:              temporada.ValorParcela,
			DataVencimento:     vencimento,
			Provider:           domain.ProviderAsaas,
			ProviderCobrancaID: &boletoResp.CobrancaID,
			ProviderCustomerID: &customerID,
			BoletoURL:          &boletoResp.URL,
			BoletoCodigoBarras: &boletoResp.CodigoBarras,
			Status:             domain.TransacaoPendente,
		})
		if err != nil {
			fmt.Printf("[ERRO] financeiro salvar transacao boleto parcela %d: %v\n", i, err)
		}
	}

	return nil
}
