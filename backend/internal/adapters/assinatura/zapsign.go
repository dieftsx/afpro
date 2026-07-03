// Package assinatura — Stub do ZapSign.
// Para ativar, implemente os métodos usando https://docs.zapsign.com.br/
package assinatura

import (
	"context"
	"errors"

	"github.com/afpro/backend/internal/ports"
)

// ErrZapSignNotImplemented é retornado por todos os métodos do stub.
var ErrZapSignNotImplemented = errors.New(
	"ZapSignAdapter: não implementado. " +
		"Para usar o ZapSign, implemente os métodos deste adapter " +
		"e troque o provider em cmd/api/main.go",
)

// ZapSignAdapter é o stub do provider ZapSign.
type ZapSignAdapter struct{}

// NewZapSignAdapter cria o stub.
func NewZapSignAdapter() *ZapSignAdapter { return &ZapSignAdapter{} }

func (z *ZapSignAdapter) EnviarParaAssinatura(_ context.Context, _ ports.DocumentoRequest) (ports.DocumentoResponse, error) {
	return ports.DocumentoResponse{}, ErrZapSignNotImplemented
}

func (z *ZapSignAdapter) BuscarStatus(_ context.Context, _ string) (string, error) {
	return "", ErrZapSignNotImplemented
}

func (z *ZapSignAdapter) CancelarDocumento(_ context.Context, _ string) error {
	return ErrZapSignNotImplemented
}
