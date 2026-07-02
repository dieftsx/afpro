// Package pagamento — Stub do Mercado Pago.
// Este arquivo é um stub documentado. Para ativar, implemente os métodos
// usando a API do Mercado Pago (https://www.mercadopago.com.br/developers)
// e troque o adapter no main.go.
package pagamento

import (
	"context"
	"errors"

	"github.com/afpero/backend/internal/domain"
	"github.com/afpero/backend/internal/ports"
)

// ErrMercadoPagoNotImplemented é retornado por todos os métodos do stub.
var ErrMercadoPagoNotImplemented = errors.New(
	"MercadoPagoAdapter: não implementado. " +
		"Para usar o Mercado Pago, implemente os métodos deste adapter " +
		"e troque o provider em cmd/api/main.go",
)

// MercadoPagoAdapter é o stub do gateway Mercado Pago.
// Implementa ports.PagamentoGateway para satisfazer a interface,
// mas todos os métodos retornam ErrMercadoPagoNotImplemented.
type MercadoPagoAdapter struct{}

// NewMercadoPagoAdapter cria o stub.
func NewMercadoPagoAdapter() *MercadoPagoAdapter {
	return &MercadoPagoAdapter{}
}

// CriarCliente — STUB. Não implementado.
func (m *MercadoPagoAdapter) CriarCliente(_ context.Context, _ domain.Participante) (string, error) {
	return "", ErrMercadoPagoNotImplemented
}

// GerarPIX — STUB. Não implementado.
func (m *MercadoPagoAdapter) GerarPIX(_ context.Context, _ ports.PIXRequest) (ports.PIXResponse, error) {
	return ports.PIXResponse{}, ErrMercadoPagoNotImplemented
}

// GerarBoleto — STUB. Não implementado.
func (m *MercadoPagoAdapter) GerarBoleto(_ context.Context, _ ports.BoletoRequest) (ports.BoletoResponse, error) {
	return ports.BoletoResponse{}, ErrMercadoPagoNotImplemented
}

// CancelarCobranca — STUB. Não implementado.
func (m *MercadoPagoAdapter) CancelarCobranca(_ context.Context, _ string) error {
	return ErrMercadoPagoNotImplemented
}
