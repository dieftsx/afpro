'use client';
// Step 2 — Exibição do termo de adesão e aceite digital.
import { useState, useRef } from 'react';
import { FileText, Check, ArrowRight, AlertCircle, ExternalLink } from 'lucide-react';
import { aceitarContrato } from '@/lib/api';
import type { AceitarContratoResult } from '@/types';

interface Step2Props {
  participanteId: string;
  contratoHTML: string;
  cpf: string;
  onComplete: (result: AceitarContratoResult) => void;
}

export default function Step2Contrato({
  participanteId,
  contratoHTML,
  cpf,
  onComplete,
}: Step2Props) {
  const [scrolledToBottom, setScrolledToBottom] = useState(false);
  const [accepted, setAccepted] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const contractRef = useRef<HTMLDivElement>(null);

  function handleScroll(e: React.UIEvent<HTMLDivElement>) {
    const el = e.currentTarget;
    const isAtBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 20;
    if (isAtBottom && !scrolledToBottom) {
      setScrolledToBottom(true);
    }
  }

  async function handleAceitar() {
    if (!accepted) {
      setError('Você precisa marcar "Li e Aceito" para continuar.');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await aceitarContrato(participanteId, cpf);
      onComplete(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao processar aceite do contrato.');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      {/* Header do step */}
      <div className="glass p-6 md:p-8">
        <div className="flex items-center gap-3 mb-4">
          <div
            className="w-10 h-10 rounded-xl flex items-center justify-center"
            style={{ background: 'var(--color-primary-dim)', color: 'var(--color-primary)' }}
          >
            <FileText size={20} />
          </div>
          <div>
            <h2 className="font-black text-xl">Termo de Adesão</h2>
            <p className="text-muted text-sm">Leia atentamente antes de assinar</p>
          </div>
        </div>

        {/* Instrução de scroll */}
        {!scrolledToBottom && (
          <div
            className="glass-sm p-3 flex items-center gap-2 text-sm mb-4"
            style={{ borderColor: 'var(--color-secondary)', color: 'var(--color-secondary)' }}
          >
            <ArrowRight size={14} className="rotate-90 shrink-0" />
            Role até o final do contrato para habilitar o botão de aceite.
          </div>
        )}

        {/* Contrato renderizado */}
        <div
          ref={contractRef}
          onScroll={handleScroll}
          className="glass-sm overflow-y-auto"
          style={{
            maxHeight: '420px',
            padding: '1.5rem',
            lineHeight: '1.75',
            fontSize: '0.875rem',
          }}
        >
          {/* Renderiza o HTML do contrato vindo do backend */}
          <ContractContent html={contratoHTML} />
        </div>

        {/* Indicador de progresso de leitura */}
        {scrolledToBottom && (
          <div
            className="flex items-center gap-2 mt-3 text-sm animate-fade-in"
            style={{ color: 'var(--color-primary)' }}
          >
            <Check size={16} />
            <span>Contrato lido na íntegra</span>
          </div>
        )}
      </div>

      {/* Aceite */}
      <div className="glass p-6 space-y-5">
        {/* Checkbox de aceite — só ativo após rolar até o fim */}
        <label
          className={`flex items-start gap-3 cursor-pointer group ${
            !scrolledToBottom ? 'opacity-50 cursor-not-allowed' : ''
          }`}
        >
          <div className="relative mt-0.5 shrink-0">
            <input
              type="checkbox"
              className="sr-only"
              checked={accepted}
              onChange={(e) => {
                if (!scrolledToBottom) return;
                setAccepted(e.target.checked);
                setError(null);
              }}
              disabled={!scrolledToBottom}
            />
            <div
              className="w-6 h-6 rounded-lg border-2 flex items-center justify-center transition-all duration-200"
              style={{
                borderColor: accepted ? 'var(--color-primary)' : 'var(--color-border)',
                background: accepted ? 'var(--color-primary)' : 'transparent',
              }}
            >
              {accepted && (
                <Check size={14} style={{ color: 'var(--color-bg)', strokeWidth: 3 }} />
              )}
            </div>
          </div>
          <span className="text-sm leading-relaxed">
            <strong>Li e Aceito</strong> todas as cláusulas do Termo de Adesão acima.
            Estou ciente das condições de pagamento, cancelamento e uso de imagem,
            e concordo integralmente com os termos apresentados.
          </span>
        </label>

        {/* Info sobre o próximo passo */}
        <div className="glass-sm p-4 text-sm text-muted">
          <p className="flex items-start gap-2">
            <ExternalLink
              size={14}
              className="shrink-0 mt-0.5"
              style={{ color: 'var(--color-secondary)' }}
            />
            Ao clicar em <strong>&quot;Assinar e Gerar PIX&quot;</strong>, o contrato será
            enviado para assinatura eletrônica via Clicksign e o QR Code PIX de entrada
            será gerado automaticamente.
          </p>
        </div>

        {/* Erro */}
        {error && (
          <div
            className="glass-sm p-4 flex items-center gap-3 text-sm"
            style={{ borderColor: 'hsl(0 72% 55%)', border: '1px solid' }}
          >
            <AlertCircle size={16} style={{ color: 'hsl(0 72% 55%)', flexShrink: 0 }} />
            <span style={{ color: 'hsl(0 72% 55%)' }}>{error}</span>
          </div>
        )}

        {/* Botão */}
        <button
          type="button"
          onClick={handleAceitar}
          disabled={!accepted || loading}
          className="btn-primary w-full text-lg py-4"
        >
          {loading ? (
            <span className="flex items-center gap-2">
              <LoadingSpinner />
              Gerando contrato e PIX...
            </span>
          ) : (
            <span className="flex items-center gap-2">
              <Check size={20} />
              Assinar e Gerar PIX
            </span>
          )}
        </button>
      </div>
    </div>
  );
}

// ─── Renderizador seguro do HTML do contrato ───────────────────────────────────
// O HTML vem do backend e contém o conteúdo em markdown convertido.
// Em produção, considere usar DOMPurify para sanitizar o HTML.
function ContractContent({ html }: { html: string }) {
  if (!html) {
    return (
      <p className="text-muted text-center py-8">
        Carregando contrato...
      </p>
    );
  }

  // Converte markdown simples em HTML (o backend envia markdown)
  const rendered = html
    .split('\n')
    .map((line) => {
      if (line.startsWith('# ')) return `<h1 class="text-xl font-black mb-4 mt-6" style="color:var(--color-primary)">${line.slice(2)}</h1>`;
      if (line.startsWith('## ')) return `<h2 class="text-base font-bold mb-3 mt-5" style="color:var(--color-secondary)">${line.slice(3)}</h2>`;
      if (line.startsWith('**') && line.endsWith('**')) return `<p class="font-bold mb-2">${line.slice(2, -2)}</p>`;
      if (line.startsWith('- ')) return `<li class="ml-4 mb-1 list-disc" style="color:var(--color-text)">${line.slice(2)}</li>`;
      if (line.trim() === '') return '<br />';
      return `<p class="mb-2" style="color:var(--color-text-muted)">${line}</p>`;
    })
    .join('');

  return (
    <div
      dangerouslySetInnerHTML={{ __html: rendered }}
      style={{ color: 'var(--color-text)' }}
    />
  );
}

function LoadingSpinner() {
  return (
    <svg className="animate-spin" width="20" height="20" viewBox="0 0 24 24" fill="none">
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" opacity="0.25" />
      <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
    </svg>
  );
}
