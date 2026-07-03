// Package assinatura fornece o adapter Clicksign para assinatura digital.
// Implementa ports.AssinaturaGateway usando a API REST do Clicksign.
// Documentação: https://developers.clicksign.com/
package assinatura

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afpro/backend/internal/ports"
)

// ClicksignAdapter implementa ports.AssinaturaGateway usando a API do Clicksign.
type ClicksignAdapter struct {
	accessToken string
	baseURL     string
	client      *http.Client
}

// NewClicksignAdapter cria uma nova instância do adapter Clicksign.
func NewClicksignAdapter(accessToken, baseURL string) *ClicksignAdapter {
	return &ClicksignAdapter{
		accessToken: accessToken,
		baseURL:     baseURL,
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

// ─── Modelos de Requisição/Resposta da API Clicksign ──────────────────────────

// Documento
type clicksignDocumentRequest struct {
	Document clicksignDocumentData `json:"document"`
}

type clicksignDocumentData struct {
	Path        string `json:"path"`        // Ex: "/contratos/afpro/contrato_joana.pdf"
	ContentBase64 string `json:"content_base64"`
	DeadlineAt  string `json:"deadline_at"` // ISO 8601
	AutoClose   bool   `json:"auto_close"`
	Locale      string `json:"locale"`    // "pt-BR"
	Sequence    bool   `json:"sequence"`  // false = assina em qualquer ordem
}

type clicksignDocumentResponse struct {
	Document struct {
		Key    string `json:"key"`
		Status string `json:"status"`
	} `json:"document"`
}

// Signatário
type clicksignSignerRequest struct {
	Signer clicksignSignerData `json:"signer"`
}

type clicksignSignerData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"` // +5569912345678
	AuthAs      string `json:"auth_as"`      // "whatsapp" | "email" | "api"
	HasDocumentation bool `json:"has_documentation"`
}

type clicksignSignerResponse struct {
	Signer struct {
		Key string `json:"key"`
	} `json:"signer"`
}

// Associação documento ↔ signatário
type clicksignSignatureRequest struct {
	List struct {
		DocumentKey string `json:"document_key"`
		SignerKey   string `json:"signer_key"`
		SignAs      string `json:"sign_as"` // "sign" = assinar
		Message     string `json:"message"`
	} `json:"list"`
}

type clicksignSignatureResponse struct {
	List struct {
		RequestSignatureKey string `json:"request_signature_key"`
		DocumentURL         string `json:"document_url"`
		SignatureURL        string `json:"signature_url"`
	} `json:"list"`
}

// Status do documento
type clicksignDocStatusResponse struct {
	Document struct {
		Key    string `json:"key"`
		Status string `json:"status"` // "running" | "closed" | "canceled"
	} `json:"document"`
}

// ─── Implementação da Interface ───────────────────────────────────────────────

// EnviarParaAssinatura realiza o fluxo completo do Clicksign:
// 1. Upload do documento (PDF em Base64)
// 2. Cadastro do signatário
// 3. Associação documento ↔ signatário
// 4. Retorna o link de assinatura
func (c *ClicksignAdapter) EnviarParaAssinatura(ctx context.Context, req ports.DocumentoRequest) (ports.DocumentoResponse, error) {
	// ── Passo 1: Upload do documento ──
	deadline := time.Now().AddDate(0, 1, 0).Format(time.RFC3339) // 30 dias para assinar
	docPath := fmt.Sprintf("/afpro/contratos/%s/%s", req.TemporadaID, req.NomeArquivo)

	docReq := clicksignDocumentRequest{
		Document: clicksignDocumentData{
			Path:          docPath,
			ContentBase64: req.ConteudoBase64,
			DeadlineAt:    deadline,
			AutoClose:     true,
			Locale:        "pt-BR",
			Sequence:      false,
		},
	}

	docResp, err := c.post(ctx, "/api/v1/documents", docReq)
	if err != nil {
		return ports.DocumentoResponse{}, fmt.Errorf("clicksign upload documento: %w", err)
	}
	defer docResp.Body.Close()

	if docResp.StatusCode != http.StatusCreated && docResp.StatusCode != http.StatusOK {
		return ports.DocumentoResponse{}, c.parseError(docResp)
	}

	var doc clicksignDocumentResponse
	if err := json.NewDecoder(docResp.Body).Decode(&doc); err != nil {
		return ports.DocumentoResponse{}, fmt.Errorf("clicksign decode documento: %w", err)
	}

	// ── Passo 2: Cadastrar signatário ──
	whatsapp := fmt.Sprintf("+55%s", req.WhatsAppSignatario) // Clicksign espera +55...

	signerReq := clicksignSignerRequest{
		Signer: clicksignSignerData{
			Name:             req.NomeSignatario,
			Email:            req.EmailSignatario,
			PhoneNumber:      whatsapp,
			AuthAs:           "email", // Autenticação via link por e-mail
			HasDocumentation: true,
		},
	}

	signerResp, err := c.post(ctx, "/api/v1/signers", signerReq)
	if err != nil {
		return ports.DocumentoResponse{}, fmt.Errorf("clicksign criar signatário: %w", err)
	}
	defer signerResp.Body.Close()

	if signerResp.StatusCode != http.StatusCreated && signerResp.StatusCode != http.StatusOK {
		return ports.DocumentoResponse{}, c.parseError(signerResp)
	}

	var signer clicksignSignerResponse
	if err := json.NewDecoder(signerResp.Body).Decode(&signer); err != nil {
		return ports.DocumentoResponse{}, fmt.Errorf("clicksign decode signatário: %w", err)
	}

	// ── Passo 3: Associar documento ao signatário ──
	var assocReq clicksignSignatureRequest
	assocReq.List.DocumentKey = doc.Document.Key
	assocReq.List.SignerKey = signer.Signer.Key
	assocReq.List.SignAs = "sign"
	assocReq.List.Message = req.MensagemEmail

	assocResp, err := c.post(ctx, "/api/v1/lists", assocReq)
	if err != nil {
		return ports.DocumentoResponse{}, fmt.Errorf("clicksign associar signatário: %w", err)
	}
	defer assocResp.Body.Close()

	if assocResp.StatusCode != http.StatusCreated && assocResp.StatusCode != http.StatusOK {
		return ports.DocumentoResponse{}, c.parseError(assocResp)
	}

	var assoc clicksignSignatureResponse
	if err := json.NewDecoder(assocResp.Body).Decode(&assoc); err != nil {
		return ports.DocumentoResponse{}, fmt.Errorf("clicksign decode associação: %w", err)
	}

	return ports.DocumentoResponse{
		DocID:          doc.Document.Key,
		SignerID:       signer.Signer.Key,
		LinkAssinatura: assoc.List.SignatureURL,
	}, nil
}

// BuscarStatus consulta o status atual de um documento no Clicksign.
func (c *ClicksignAdapter) BuscarStatus(ctx context.Context, docID string) (string, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/api/v1/documents/%s", docID))
	if err != nil {
		return "", fmt.Errorf("clicksign BuscarStatus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.parseError(resp)
	}

	var status clicksignDocStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", fmt.Errorf("clicksign BuscarStatus decode: %w", err)
	}

	return status.Document.Status, nil
}

// CancelarDocumento cancela o processo de assinatura de um documento.
func (c *ClicksignAdapter) CancelarDocumento(ctx context.Context, docID string) error {
	resp, err := c.delete(ctx, fmt.Sprintf("/api/v1/documents/%s/cancel", docID))
	if err != nil {
		return fmt.Errorf("clicksign CancelarDocumento: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return c.parseError(resp)
	}

	return nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// HTMLParaBase64 converte um HTML de contrato para Base64.
// O Clicksign aceita PDFs; numa implementação de produção, use chromedp ou
// uma biblioteca de conversão HTML→PDF antes de chamar esta função.
// Para o MVP, enviamos o HTML diretamente encodado em Base64.
func HTMLParaBase64(html string) string {
	return base64.StdEncoding.EncodeToString([]byte(html))
}

func (c *ClicksignAdapter) post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s?access_token=%s", c.baseURL, path, c.accessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

func (c *ClicksignAdapter) get(ctx context.Context, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s?access_token=%s", c.baseURL, path, c.accessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *ClicksignAdapter) delete(ctx context.Context, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s?access_token=%s", c.baseURL, path, c.accessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *ClicksignAdapter) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("clicksign error [%d]: %s", resp.StatusCode, string(body))
}
