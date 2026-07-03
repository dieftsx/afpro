'use client';
// Step 3 — Checkout PIX com QR Code, copia-e-cola e countdown timer.
import { useState, useEffect, useCallback } from 'react';
import Image from 'next/image';
import {
  CheckCircle2,
  Copy,
  Check,
  Clock,
  AlertTriangle,
  RefreshCw,
  QrCode,
} from 'lucide-react';
import type { AceitarContratoResult } from '@/types';

interface Step3Props {
  pixData: AceitarContratoResult;
}

export default function Step3Checkout({ pixData }: Step3Props) {
  const [copied, setCopied] = useState(false);
  const [timeLeft, setTimeLeft] = useState<number>(0);
  const [expired, setExpired] = useState(false);

  // Calcular o tempo restante em segundos
  const calcTimeLeft = useCallback(() => {
    const expiresAt = new Date(pixData.pix_expiracao).getTime();
    const now = Date.now();
    const diff = Math.floor((expiresAt - now) / 1000);
    return Math.max(0, diff);
  }, [pixData.pix_expiracao]);

  // Countdown
  useEffect(() => {
    const initial = calcTimeLeft();
    setTimeLeft(initial);
    if (initial <= 0) {
      setExpired(true);
      return;
    }

    const interval = setInterval(() => {
      const remaining = calcTimeLeft();
      setTimeLeft(remaining);
      if (remaining <= 0) {
        setExpired(true);
        clearInterval(interval);
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [calcTimeLeft]);

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(pixData.pix_copia_cola);
      setCopied(true);
      setTimeout(() => setCopied(false), 3000);
    } catch {
      // Fallback para browsers sem API de clipboard
      const el = document.createElement('textarea');
      el.value = pixData.pix_copia_cola;
      document.body.appendChild(el);
      el.select();
      document.execCommand('copy');
      document.body.removeChild(el);
      setCopied(true);
      setTimeout(() => setCopied(false), 3000);
    }
  }

  const minutes = Math.floor(timeLeft / 60);
  const seconds = timeLeft % 60;
  const timerPercentage = (timeLeft / (30 * 60)) * 100;
  const isUrgent = timeLeft < 5 * 60; // < 5 minutos = alerta vermelho

  if (expired) {
    return <ExpiredState />;
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header de sucesso */}
      <div className="glass p-6 md:p-8 text-center">
        <div className="flex justify-center mb-4">
          <div
            className="w-16 h-16 rounded-full flex items-center justify-center animate-pulse-glow"
            style={{ background: 'var(--color-primary-dim)' }}
          >
            <CheckCircle2 size={32} style={{ color: 'var(--color-primary)' }} />
          </div>
        </div>

        <h2 className="text-2xl font-black mb-2">
          Contrato assinado! 🎉
        </h2>
        <p className="text-muted">
          Agora realize o pagamento da entrada via PIX para confirmar sua vaga.
        </p>

        {/* Valor */}
        <div className="mt-6 glass-sm inline-flex px-6 py-3 gap-2 items-center">
          <span className="text-muted text-sm">Valor da entrada:</span>
          <span className="gradient-text font-black text-2xl">
            {formatBRL(pixData.valor_entrada)}
          </span>
        </div>
      </div>

      {/* Timer + QR Code */}
      <div className="glass p-6 md:p-8">
        {/* Countdown */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2 text-sm font-semibold">
              <Clock
                size={16}
                style={{ color: isUrgent ? 'hsl(0 72% 55%)' : 'var(--color-secondary)' }}
              />
              <span style={{ color: isUrgent ? 'hsl(0 72% 55%)' : 'var(--color-text)' }}>
                {isUrgent ? '⚠️ Expira em breve!' : 'Este PIX expira em:'}
              </span>
            </div>
            <div
              className="font-black text-2xl tabular-nums"
              style={{
                color: isUrgent ? 'hsl(0 72% 55%)' : 'var(--color-primary)',
              }}
            >
              {String(minutes).padStart(2, '0')}:{String(seconds).padStart(2, '0')}
            </div>
          </div>

          {/* Barra de progresso do timer */}
          <div className="progress-track">
            <div
              className="progress-fill"
              style={{
                width: `${timerPercentage}%`,
                background: isUrgent
                  ? 'linear-gradient(90deg, hsl(0 72% 55%), hsl(0 80% 70%))'
                  : undefined,
                transition: 'width 1s linear',
              }}
            />
          </div>
        </div>

        {/* QR Code */}
        <div className="flex flex-col items-center gap-6">
          <div className="qr-container">
            {pixData.pix_qr_code_base64 ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={`data:image/png;base64,${pixData.pix_qr_code_base64}`}
                alt="QR Code PIX"
                width={200}
                height={200}
                style={{ imageRendering: 'pixelated' }}
              />
            ) : (
              <div className="w-[200px] h-[200px] flex items-center justify-center">
                <QrCode size={48} style={{ color: 'hsl(0 0% 20%)' }} />
              </div>
            )}
          </div>

          <p className="text-muted text-sm text-center">
            Escaneie o QR Code com o aplicativo do seu banco ou use o código abaixo.
          </p>
        </div>
      </div>

      {/* Copia e cola */}
      <div className="glass p-6">
        <p className="text-sm font-semibold mb-3 text-muted uppercase tracking-wider">
          Código PIX — Copia e Cola
        </p>

        <div className="glass-sm p-4 flex items-center gap-3">
          <code
            className="flex-1 text-xs break-all leading-relaxed font-mono"
            style={{ color: 'var(--color-accent)' }}
          >
            {pixData.pix_copia_cola}
          </code>

          <button
            type="button"
            onClick={handleCopy}
            className="shrink-0 btn-outline text-sm px-4 py-2.5"
            style={{
              borderColor: copied ? 'var(--color-primary)' : undefined,
              color: copied ? 'var(--color-primary)' : undefined,
            }}
          >
            {copied ? (
              <span className="flex items-center gap-1.5">
                <Check size={14} />
                Copiado!
              </span>
            ) : (
              <span className="flex items-center gap-1.5">
                <Copy size={14} />
                Copiar
              </span>
            )}
          </button>
        </div>
      </div>

      {/* Instruções */}
      <div className="glass-sm p-5 space-y-2">
        <p className="text-sm font-bold mb-3">Como pagar:</p>
        {[
          'Abra o aplicativo do seu banco.',
          'Acesse a opção de pagamento via PIX.',
          'Escaneie o QR Code ou cole o código acima.',
          `Confirme o pagamento de ${formatBRL(pixData.valor_entrada)}.`,
          'Pronto! Você receberá a confirmação por e-mail em até 5 minutos.',
        ].map((step, i) => (
          <div key={i} className="flex items-start gap-3 text-sm">
            <span
              className="w-5 h-5 rounded-full flex items-center justify-center text-xs font-bold shrink-0 mt-0.5"
              style={{
                background: 'var(--color-primary-dim)',
                color: 'var(--color-primary)',
              }}
            >
              {i + 1}
            </span>
            <span className="text-muted">{step}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

// ─── Estado de expiração ───────────────────────────────────────────────────────
function ExpiredState() {
  return (
    <div className="glass p-8 md:p-12 text-center space-y-6 animate-fade-in">
      <div className="flex justify-center">
        <div
          className="w-20 h-20 rounded-full flex items-center justify-center"
          style={{ background: 'color-mix(in srgb, hsl(0 72% 55%) 20%, transparent)' }}
        >
          <AlertTriangle size={40} style={{ color: 'hsl(0 72% 55%)' }} />
        </div>
      </div>

      <div>
        <h2 className="text-2xl font-black mb-2" style={{ color: 'hsl(0 72% 55%)' }}>
          PIX expirado
        </h2>
        <p className="text-muted">
          O QR Code expirou após 30 minutos. Não se preocupe — sua inscrição foi salva!
          Entre em contato com a AFPERO para reagendar o pagamento.
        </p>
      </div>

      <a
        href="https://wa.me/5569000000000"
        target="_blank"
        rel="noopener noreferrer"
        className="btn-primary inline-flex gap-2"
      >
        <RefreshCw size={18} />
        Falar com a AFPERO
      </a>
    </div>
  );
}

function formatBRL(value: number) {
  return value.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
}
