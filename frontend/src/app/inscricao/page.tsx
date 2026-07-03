'use client';
// Página de inscrição — container do Multi-step form.
// Esta é a página principal que orquestra os 3 steps.
import { useState } from 'react';
import { Fish, ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import Step1Cadastro from '@/components/inscricao/Step1Cadastro';
import Step2Contrato from '@/components/inscricao/Step2Contrato';
import Step3Checkout from '@/components/inscricao/Step3Checkout';
import type { AceitarContratoResult, CriarInscricaoResult, StepInscricao } from '@/types';

export default function InscricaoPage() {
  const [step, setStep] = useState<StepInscricao>(1);

  // Dados que fluem entre steps
  const [step1Result, setStep1Result] = useState<CriarInscricaoResult | null>(null);
  const [step2Result, setStep2Result] = useState<AceitarContratoResult | null>(null);

  // CPF mantido em memória apenas durante a sessão (nunca persistido no frontend)
  const [cpfSession, setCpfSession] = useState<string>('');

  const stepLabels = ['Cadastro', 'Contrato', 'Pagamento'];
  const progress = ((step - 1) / 2) * 100;

  function handleStep1Complete(result: CriarInscricaoResult, cpf: string) {
    setStep1Result(result);
    setCpfSession(cpf);
    setStep(2);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  function handleStep2Complete(result: AceitarContratoResult) {
    setStep2Result(result);
    setStep(3);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  return (
    <main className="relative min-h-screen flex flex-col">
      {/* Background */}
      <div className="fixed inset-0 pointer-events-none z-0">
        <div
          className="absolute inset-0 opacity-30"
          style={{
            background: 'radial-gradient(ellipse 70% 50% at 50% 0%, var(--color-primary-dim), transparent)',
          }}
        />
      </div>

      <div className="relative z-10 max-w-2xl mx-auto w-full px-4 py-8 flex-1 flex flex-col">

        {/* Header */}
        <header className="mb-8">
          <div className="flex items-center justify-between mb-6">
            <Link
              href="/"
              className="flex items-center gap-2 text-muted hover:text-accent transition-colors text-sm"
              style={{ color: 'var(--color-text-muted)' }}
            >
              <ArrowLeft size={16} />
              Voltar
            </Link>

            <div className="flex items-center gap-2">
              <Fish size={16} style={{ color: 'var(--color-primary)' }} />
              <span className="text-sm font-bold opacity-70 uppercase tracking-widest">AFPERO</span>
            </div>
          </div>

          {/* Step indicator */}
          <div className="glass-sm p-4">
            {/* Labels */}
            <div className="flex justify-between mb-3">
              {stepLabels.map((label, i) => {
                const stepNum = (i + 1) as StepInscricao;
                const isActive = stepNum === step;
                const isDone = stepNum < step;

                return (
                  <div key={i} className="flex items-center gap-2">
                    <div
                      className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold transition-all duration-300"
                      style={{
                        background: isDone || isActive
                          ? 'linear-gradient(135deg, var(--color-primary), var(--color-secondary))'
                          : 'var(--color-surface-high)',
                        color: isDone || isActive ? 'var(--color-bg)' : 'var(--color-text-muted)',
                      }}
                    >
                      {isDone ? '✓' : stepNum}
                    </div>
                    <span
                      className="text-xs font-semibold hidden sm:block"
                      style={{
                        color: isActive
                          ? 'var(--color-primary)'
                          : isDone
                          ? 'var(--color-text)'
                          : 'var(--color-text-muted)',
                      }}
                    >
                      {label}
                    </span>
                  </div>
                );
              })}
            </div>

            {/* Barra de progresso */}
            <div className="progress-track">
              <div className="progress-fill" style={{ width: `${progress}%` }} />
            </div>
          </div>
        </header>

        {/* Conteúdo do step atual */}
        <div className="flex-1 animate-fade-in">
          {step === 1 && (
            <Step1Cadastro onComplete={handleStep1Complete} />
          )}

          {step === 2 && step1Result && (
            <Step2Contrato
              participanteId={step1Result.participante_id}
              contratoHTML={step1Result.contrato_html}
              cpf={cpfSession}
              onComplete={handleStep2Complete}
            />
          )}

          {step === 3 && step2Result && (
            <Step3Checkout
              pixData={step2Result}
            />
          )}
        </div>
      </div>
    </main>
  );
}
