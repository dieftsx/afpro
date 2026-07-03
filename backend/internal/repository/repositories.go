// Package repository fornece acesso ao Supabase via client REST.
// Todas as operações usam a service_role key para contornar as RLS policies.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	supa "github.com/supabase-community/supabase-go"

	"github.com/afpro/backend/internal/domain"
)

// ─── Repositório de Temporadas ────────────────────────────────────────────────

// TemporadaRepository acessa a tabela temporadas no Supabase.
type TemporadaRepository struct {
	client *supa.Client
}

// NewTemporadaRepository cria uma nova instância.
func NewTemporadaRepository(client *supa.Client) *TemporadaRepository {
	return &TemporadaRepository{client: client}
}

// BuscarAtiva retorna a temporada com is_ativa = TRUE.
// Usa a view vw_temporada_ativa para expor apenas campos públicos.
func (r *TemporadaRepository) BuscarAtiva(ctx context.Context) (*domain.Temporada, error) {
	var results []domain.Temporada

	_, err := r.client.
		From("vw_temporada_ativa").
		Select("*", "exact", false).
		Single().
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("temporada BuscarAtiva: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("nenhuma temporada ativa encontrada")
	}

	return &results[0], nil
}

// BuscarPorID busca uma temporada pelo ID.
func (r *TemporadaRepository) BuscarPorID(ctx context.Context, id string) (*domain.Temporada, error) {
	var results []domain.Temporada

	_, err := r.client.
		From("temporadas").
		Select("*", "exact", false).
		Eq("id", id).
		Single().
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("temporada BuscarPorID: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("temporada não encontrada: %s", id)
	}

	return &results[0], nil
}

// DecrementarVaga decrementa vagas_disponiveis em 1 de forma atômica.
// Usa a RPC do Supabase para garantir atomicidade.
func (r *TemporadaRepository) DecrementarVaga(ctx context.Context, temporadaID string) error {
	r.client.Rpc("decrementar_vaga_temporada", "exact", map[string]interface{}{
		"p_temporada_id": temporadaID,
	})
	return nil
}

// ─── Repositório de Participantes ─────────────────────────────────────────────

// ParticipanteRepository acessa a tabela participantes no Supabase.
type ParticipanteRepository struct {
	client *supa.Client
}

// NewParticipanteRepository cria uma nova instância.
func NewParticipanteRepository(client *supa.Client) *ParticipanteRepository {
	return &ParticipanteRepository{client: client}
}

// Criar insere uma nova participante e retorna a entidade com o ID gerado.
func (r *ParticipanteRepository) Criar(ctx context.Context, p domain.Participante) (*domain.Participante, error) {
	// Mapeamento explícito para evitar enviar campos calculados/sensíveis
	payload := map[string]interface{}{
		"temporada_id":           p.TemporadaID,
		"nome_completo":          p.NomeCompleto,
		"email":                  p.Email,
		"whatsapp":               p.WhatsApp,
		"cpf_hash":               p.CPFHash,
		"cpf_last4":              p.CPFLast4,
		"emergencia_nome":        p.EmergenciaNome,
		"emergencia_telefone":    p.EmergenciaTelefone,
		"tamanho_camisa":         string(p.TamanhoCamisa),
		"tamanho_balaclava":      string(p.TamanhoBalaclava),
		"tamanho_bone":           string(p.TamanhoBone),
		"status":                 string(p.Status),
		"consentimento_lgpd":     p.ConsentimentoLGPD,
		"consentimento_at":       p.ConsentimentoAt,
	}

	if p.DataNascimento != nil {
		payload["data_nascimento"] = p.DataNascimento.Format("2006-01-02")
	}
	if p.EmergenciaParentesco != nil {
		payload["emergencia_parentesco"] = *p.EmergenciaParentesco
	}
	if p.IPOrigem != nil {
		payload["ip_origem"] = *p.IPOrigem
	}
	if p.UserAgent != nil {
		payload["user_agent"] = *p.UserAgent
	}

	var results []domain.Participante
	_, err := r.client.
		From("participantes").
		Insert(payload, false, "", "representation", "exact").
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("participante Criar: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("participante Criar: nenhum resultado retornado")
	}

	return &results[0], nil
}

// BuscarPorID busca uma participante pelo ID.
func (r *ParticipanteRepository) BuscarPorID(ctx context.Context, id string) (*domain.Participante, error) {
	var results []domain.Participante

	_, err := r.client.
		From("participantes").
		Select("*", "exact", false).
		Eq("id", id).
		Single().
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("participante BuscarPorID: %w", err)
	}

	if len(results) == 0 {
		return nil, ErrNaoEncontrado
	}

	return &results[0], nil
}

// ExistePorCPFHashETemporada verifica se já existe inscrição com este CPF na temporada.
func (r *ParticipanteRepository) ExistePorCPFHashETemporada(ctx context.Context, cpfHash, temporadaID string) (bool, error) {
	var results []struct{ ID string `json:"id"` }

	_, err := r.client.
		From("participantes").
		Select("id", "exact", false).
		Eq("cpf_hash", cpfHash).
		Eq("temporada_id", temporadaID).
		Limit(1, "").
		ExecuteTo(&results)

	if err != nil {
		return false, fmt.Errorf("participante ExistePorCPFHash: %w", err)
	}

	return len(results) > 0, nil
}

// AtualizarStatus atualiza o status de inscrição de uma participante.
func (r *ParticipanteRepository) AtualizarStatus(ctx context.Context, id string, status domain.StatusInscricao) error {
	_, _, err := r.client.
		From("participantes").
		Update(map[string]interface{}{"status": string(status)}, "representation", "exact").
		Eq("id", id).
		Execute()

	if err != nil {
		return fmt.Errorf("participante AtualizarStatus: %w", err)
	}

	return nil
}

// ─── Repositório de Documentos Assinados ──────────────────────────────────────

// DocumentoRepository acessa a tabela documentos_assinados no Supabase.
type DocumentoRepository struct {
	client *supa.Client
}

// NewDocumentoRepository cria uma nova instância.
func NewDocumentoRepository(client *supa.Client) *DocumentoRepository {
	return &DocumentoRepository{client: client}
}

// Criar insere um novo documento.
func (r *DocumentoRepository) Criar(ctx context.Context, doc domain.DocumentoAssinado) (*domain.DocumentoAssinado, error) {
	payload := map[string]interface{}{
		"participante_id": doc.ParticipanteID,
		"temporada_id":    doc.TemporadaID,
		"provider":        string(doc.Provider),
		"status":          string(doc.Status),
	}

	if doc.ConteudoHTML != nil {
		payload["conteudo_html"] = *doc.ConteudoHTML
	}

	var results []domain.DocumentoAssinado
	_, err := r.client.
		From("documentos_assinados").
		Insert(payload, false, "", "representation", "exact").
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("documento Criar: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("documento Criar: nenhum resultado retornado")
	}

	return &results[0], nil
}

// AtualizarAposEnvio atualiza o documento com os dados retornados pelo provider.
func (r *DocumentoRepository) AtualizarAposEnvio(ctx context.Context, id, providerDocID, providerSignerID, linkAssinatura string) error {
	_, _, err := r.client.
		From("documentos_assinados").
		Update(map[string]interface{}{
			"provider_doc_id":    providerDocID,
			"provider_signer_id": providerSignerID,
			"link_assinatura":    linkAssinatura,
			"status":             string(domain.DocEnviado),
		}, "representation", "exact").
		Eq("id", id).
		Execute()

	return err
}

// AtualizarAposAssinatura marca o documento como assinado via webhook.
func (r *DocumentoRepository) AtualizarAposAssinatura(ctx context.Context, providerDocID string, payload interface{}) error {
	now := "now()"
	_, _, err := r.client.
		From("documentos_assinados").
		Update(map[string]interface{}{
			"status":              string(domain.DocAssinado),
			"assinado_em":         now,
			"webhook_payload":     payload,
			"webhook_recebido_em": now,
		}, "representation", "exact").
		Eq("provider_doc_id", providerDocID).
		Execute()

	return err
}

// BuscarPorProviderDocID busca um documento pelo ID externo do provider.
func (r *DocumentoRepository) BuscarPorProviderDocID(ctx context.Context, providerDocID string) (*domain.DocumentoAssinado, error) {
	var results []domain.DocumentoAssinado

	_, err := r.client.
		From("documentos_assinados").
		Select("*", "exact", false).
		Eq("provider_doc_id", providerDocID).
		Single().
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("documento BuscarPorProviderDocID: %w", err)
	}

	if len(results) == 0 {
		return nil, ErrNaoEncontrado
	}

	return &results[0], nil
}

// ─── Repositório de Transações ────────────────────────────────────────────────

// TransacaoRepository acessa a tabela transacoes no Supabase.
type TransacaoRepository struct {
	client *supa.Client
}

// NewTransacaoRepository cria uma nova instância.
func NewTransacaoRepository(client *supa.Client) *TransacaoRepository {
	return &TransacaoRepository{client: client}
}

// Criar insere uma nova transação.
func (r *TransacaoRepository) Criar(ctx context.Context, t domain.Transacao) (*domain.Transacao, error) {
	payload := map[string]interface{}{
		"participante_id":  t.ParticipanteID,
		"temporada_id":     t.TemporadaID,
		"tipo":             string(t.Tipo),
		"descricao":        t.Descricao,
		"valor":            t.Valor,
		"data_vencimento":  t.DataVencimento.Format("2006-01-02"),
		"provider":         string(t.Provider),
		"status":           string(t.Status),
	}

	if t.NumeroParcela != nil {
		payload["numero_parcela"] = *t.NumeroParcela
	}
	if t.ProviderCobrancaID != nil {
		payload["provider_cobranca_id"] = *t.ProviderCobrancaID
	}
	if t.ProviderCustomerID != nil {
		payload["provider_customer_id"] = *t.ProviderCustomerID
	}
	if t.PIXQRCodeBase64 != nil {
		payload["pix_qr_code_base64"] = *t.PIXQRCodeBase64
	}
	if t.PIXCopiaCola != nil {
		payload["pix_copia_cola"] = *t.PIXCopiaCola
	}
	if t.PIXExpiracao != nil {
		payload["pix_expiracao"] = t.PIXExpiracao.Format(time.RFC3339)
	}
	if t.BoletoURL != nil {
		payload["boleto_url"] = *t.BoletoURL
	}
	if t.BoletoCodigoBarras != nil {
		payload["boleto_codigo_barras"] = *t.BoletoCodigoBarras
	}

	var results []domain.Transacao
	_, err := r.client.
		From("transacoes").
		Insert(payload, false, "", "representation", "exact").
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("transacao Criar: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("transacao Criar: nenhum resultado retornado")
	}

	return &results[0], nil
}

// BuscarEntradaPIXPorParticipante retorna a transação de entrada PIX de uma participante.
func (r *TransacaoRepository) BuscarEntradaPIXPorParticipante(ctx context.Context, participanteID string) (*domain.Transacao, error) {
	var results []domain.Transacao

	_, err := r.client.
		From("transacoes").
		Select("*", "exact", false).
		Eq("participante_id", participanteID).
		Eq("tipo", string(domain.TipoEntradaPIX)).
		Neq("status", string(domain.TransacaoCancelada)).
		Single().
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("transacao BuscarEntradaPIX: %w", err)
	}

	if len(results) == 0 {
		return nil, ErrNaoEncontrado
	}

	return &results[0], nil
}

// BuscarPorProviderCobrancaID busca uma transação pelo ID externo do provider.
func (r *TransacaoRepository) BuscarPorProviderCobrancaID(ctx context.Context, cobrancaID string) (*domain.Transacao, error) {
	var results []domain.Transacao

	_, err := r.client.
		From("transacoes").
		Select("*", "exact", false).
		Eq("provider_cobranca_id", cobrancaID).
		Single().
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("transacao BuscarPorProviderID: %w", err)
	}

	if len(results) == 0 {
		return nil, ErrNaoEncontrado
	}

	return &results[0], nil
}

// MarcarComoPago atualiza a transação para status pago.
func (r *TransacaoRepository) MarcarComoPago(ctx context.Context, id string, payload interface{}) error {
	_, _, err := r.client.
		From("transacoes").
		Update(map[string]interface{}{
			"status":              string(domain.TransacaoPaga),
			"pago_em":             "now()",
			"webhook_payload":     payload,
			"webhook_recebido_em": "now()",
		}, "representation", "exact").
		Eq("id", id).
		Execute()

	return err
}

// ─── Erro Padrão ──────────────────────────────────────────────────────────────

// ErrNaoEncontrado indica que o recurso não foi encontrado no banco de dados.
var ErrNaoEncontrado = errors.New("recurso não encontrado")
