// Funções utilitárias para comunicação com a API Go
import type {
  Temporada,
  CriarInscricaoPayload,
  CriarInscricaoResult,
  AceitarContratoResult,
} from '@/types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Erro desconhecido' }));
    throw new Error(err.error ?? `Erro ${res.status}`);
  }

  return res.json() as Promise<T>;
}

/** Busca o tema e configurações da temporada ativa */
export async function getTemporadaAtiva(): Promise<Temporada> {
  return apiFetch<Temporada>('/api/v1/temporada/ativa');
}

/** Cria uma nova inscrição (Step 1) */
export async function criarInscricao(
  payload: CriarInscricaoPayload
): Promise<CriarInscricaoResult> {
  return apiFetch<CriarInscricaoResult>('/api/v1/inscricoes', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

/** Aceita o contrato e gera o PIX de entrada (Step 2) */
export async function aceitarContrato(
  participanteId: string,
  cpf: string
): Promise<AceitarContratoResult> {
  return apiFetch<AceitarContratoResult>(
    `/api/v1/inscricoes/${participanteId}/contrato/aceitar`,
    {
      method: 'POST',
      body: JSON.stringify({ cpf }),
    }
  );
}

/** Busca os dados do PIX de uma inscrição (para polling/atualização) */
export async function getPIX(
  participanteId: string
): Promise<AceitarContratoResult> {
  return apiFetch<AceitarContratoResult>(
    `/api/v1/inscricoes/${participanteId}/pix`
  );
}
