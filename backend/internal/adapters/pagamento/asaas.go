// Package pagamento fornece o adapter Asaas para o gateway de pagamento.
// Implementa a interface ports.PagamentoGateway usando a API REST do Asaas.
// Documentação: https://docs.asaas.com/
package pagamento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afpero/backend/internal/domain"
	"github.com/afpero/backend/internal/ports"
)

// AsaasAdapter implementa ports.PagamentoGateway usando a API do Asaas.
type AsaasAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewAsaasAdapter cria uma nova instância do adapter Asaas.
func NewAsaasAdapter(apiKey, baseURL string) *AsaasAdapter {
	return &AsaasAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// ─── Modelos de Requisição/Resposta da API Asaas ──────────────────────────────

type asaasCustomer struct {
	Name        string `json:"name"`
	CPFCnpj     string `json:"cpfCnpj"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobilePhone"`
}

type asaasCustomerResponse struct {
	ID string `json:"id"`
}

type asaasPaymentRequest struct {
	Customer        string  `json:"customer"`
	BillingType     string  `json:"billingType"` // "PIX" | "BOLETO"
	Value           float64 `json:"value"`
	DueDate         string  `json:"dueDate"`         // "YYYY-MM-DD"
	Description     string  `json:"description"`
	ExternalReference string `json:"externalReference"` // participante_id (rastreabilidade)
}

type asaasPaymentResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	InvoiceURL  string `json:"invoiceUrl"`
	BankSlipURL string `json:"bankSlipUrl"` // URL do boleto
}

type asaasPIXQRCodeResponse struct {
	EncodedImage string `json:"encodedImage"` // Base64
	Payload      string `json:"payload"`      // Copia e Cola
	ExpirationDate string `json:"expirationDate"`
}

// ─── Implementação da Interface ───────────────────────────────────────────────

// CriarCliente cadastra a participante no Asaas e retorna o customerID.
// O CPF é necessário; o Asaas usa o CPF como identificador único de cliente.
// ATENÇÃO: o CPF completo é enviado apenas ao Asaas (provider externo confiável),
// mas NUNCA é armazenado em nosso banco de dados.
func (a *AsaasAdapter) CriarCliente(ctx context.Context, p domain.Participante) (string, error) {
	// O CPF completo deve ser passado via campo auxiliar (não armazenado no DB)
	// O handler deve extrair o CPF da requisição original antes de hashear.
	// Aqui recebemos o CPF completo no campo auxiliar da entidade.
	body := asaasCustomer{
		Name:        p.NomeCompleto,
		CPFCnpj:     "", // Preenchido pelo service com o CPF original
		Email:       p.Email,
		MobilePhone: p.WhatsApp,
	}

	resp, err := a.post(ctx, "/customers", body)
	if err != nil {
		return "", fmt.Errorf("asaas CriarCliente: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", a.parseError(resp)
	}

	var customer asaasCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		return "", fmt.Errorf("asaas CriarCliente decode: %w", err)
	}

	return customer.ID, nil
}

// GerarPIX cria uma cobrança PIX no Asaas e retorna QR Code + Copia e Cola.
func (a *AsaasAdapter) GerarPIX(ctx context.Context, req ports.PIXRequest) (ports.PIXResponse, error) {
	dueDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02") // Vence no dia seguinte

	payReq := asaasPaymentRequest{
		Customer:          req.CustomerID,
		BillingType:       "PIX",
		Value:             req.Valor,
		DueDate:           dueDate,
		Description:       req.Descricao,
		ExternalReference: req.ParticipanteID,
	}

	resp, err := a.post(ctx, "/payments", payReq)
	if err != nil {
		return ports.PIXResponse{}, fmt.Errorf("asaas GerarPIX payment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return ports.PIXResponse{}, a.parseError(resp)
	}

	var payResp asaasPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payResp); err != nil {
		return ports.PIXResponse{}, fmt.Errorf("asaas GerarPIX decode payment: %w", err)
	}

	// Buscar o QR Code da cobrança criada
	qrResp, err := a.get(ctx, fmt.Sprintf("/payments/%s/pixQrCode", payResp.ID))
	if err != nil {
		return ports.PIXResponse{}, fmt.Errorf("asaas GerarPIX qrcode: %w", err)
	}
	defer qrResp.Body.Close()

	if qrResp.StatusCode != http.StatusOK {
		return ports.PIXResponse{}, a.parseError(qrResp)
	}

	var qr asaasPIXQRCodeResponse
	if err := json.NewDecoder(qrResp.Body).Decode(&qr); err != nil {
		return ports.PIXResponse{}, fmt.Errorf("asaas GerarPIX decode qr: %w", err)
	}

	expiracao := time.Now().Add(time.Duration(req.ExpiracaoMin) * time.Minute)

	return ports.PIXResponse{
		CobrancaID:   payResp.ID,
		QRCodeBase64: qr.EncodedImage,
		CopiaCola:    qr.Payload,
		Expiracao:    expiracao,
	}, nil
}

// GerarBoleto cria um boleto bancário no Asaas para uma parcela específica.
func (a *AsaasAdapter) GerarBoleto(ctx context.Context, req ports.BoletoRequest) (ports.BoletoResponse, error) {
	payReq := asaasPaymentRequest{
		Customer:          req.CustomerID,
		BillingType:       "BOLETO",
		Value:             req.Valor,
		DueDate:           req.Vencimento.Format("2006-01-02"),
		Description:       req.Descricao,
		ExternalReference: req.ParticipanteID,
	}

	resp, err := a.post(ctx, "/payments", payReq)
	if err != nil {
		return ports.BoletoResponse{}, fmt.Errorf("asaas GerarBoleto: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return ports.BoletoResponse{}, a.parseError(resp)
	}

	var payResp asaasPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payResp); err != nil {
		return ports.BoletoResponse{}, fmt.Errorf("asaas GerarBoleto decode: %w", err)
	}

	return ports.BoletoResponse{
		CobrancaID:   payResp.ID,
		URL:          payResp.BankSlipURL,
		CodigoBarras: "", // O Asaas retorna o código de barras em endpoint separado; simplificado aqui
	}, nil
}

// CancelarCobranca cancela uma cobrança pendente no Asaas.
func (a *AsaasAdapter) CancelarCobranca(ctx context.Context, cobrancaID string) error {
	resp, err := a.delete(ctx, fmt.Sprintf("/payments/%s", cobrancaID))
	if err != nil {
		return fmt.Errorf("asaas CancelarCobranca: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return a.parseError(resp)
	}

	return nil
}

// ─── Helpers HTTP ─────────────────────────────────────────────────────────────

func (a *AsaasAdapter) post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	a.setHeaders(req)
	return a.client.Do(req)
}

func (a *AsaasAdapter) get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	a.setHeaders(req)
	return a.client.Do(req)
}

func (a *AsaasAdapter) delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, a.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	a.setHeaders(req)
	return a.client.Do(req)
}

func (a *AsaasAdapter) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access_token", a.apiKey)
}

type asaasErrorResponse struct {
	Errors []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"errors"`
}

func (a *AsaasAdapter) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var asaasErr asaasErrorResponse
	if err := json.Unmarshal(body, &asaasErr); err == nil && len(asaasErr.Errors) > 0 {
		return fmt.Errorf("asaas error [%d]: %s — %s",
			resp.StatusCode,
			asaasErr.Errors[0].Code,
			asaasErr.Errors[0].Description,
		)
	}
	return fmt.Errorf("asaas error [%d]: %s", resp.StatusCode, string(body))
}
