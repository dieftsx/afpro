'use client';
// CTASection — seção final de call-to-action com urgência e botão de inscrição.
import Link from 'next/link';
import { ArrowRight, Fish, Sparkles } from 'lucide-react';
import type { Temporada } from '@/types';

interface CTASectionProps {
  temporada: Temporada | null;
  inscricoesAbertas: boolean;
}

export default function CTASection({ temporada, inscricoesAbertas }: CTASectionProps) {
  if (!inscricoesAbertas) return null;

  return (
    <section className="relative z-10 py-24 px-6">
      <div className="max-w-4xl mx-auto text-center">
        <div className="glass p-12 md:p-16 relative overflow-hidden">
          {/* Glow de fundo */}
          <div
            className="absolute inset-0 opacity-20 blur-3xl"
            style={{
              background: 'radial-gradient(circle, var(--color-primary) 0%, var(--color-secondary) 100%)',
            }}
          />

          <div className="relative z-10">
            <div className="flex justify-center mb-6">
              <span className="badge badge-secondary text-base px-5 py-2">
                <Sparkles size={14} />
                {temporada?.vagas_disponiveis ?? 0} vagas disponíveis
              </span>
            </div>

            <h2 className="text-4xl md:text-6xl font-black mb-4">
              Sua vaga está
              <br />
              <span className="gradient-text">esperando por você</span>
            </h2>

            <p className="text-muted text-lg mb-10 max-w-lg mx-auto">
              Não deixe para depois. As vagas são limitadas e o processo é simples:
              preencha o formulário, assine o contrato e garanta a entrada.
            </p>

            {/* Passos rápidos */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-10">
              {[
                { n: '1', label: 'Dados pessoais' },
                { n: '2', label: 'Assinar contrato' },
                { n: '3', label: 'Pagar entrada PIX' },
              ].map((step, i) => (
                <div key={i} className="flex items-center gap-3 glass-sm px-4 py-3">
                  <div
                    className="w-8 h-8 rounded-full flex items-center justify-center text-sm font-black"
                    style={{
                      background: 'linear-gradient(135deg, var(--color-primary), var(--color-secondary))',
                      color: 'var(--color-bg)',
                    }}
                  >
                    {step.n}
                  </div>
                  <span className="text-sm font-semibold">{step.label}</span>
                  {i < 2 && (
                    <ArrowRight size={14} className="text-muted hidden sm:block" />
                  )}
                </div>
              ))}
            </div>

            <Link
              href="/inscricao"
              className="btn-primary text-xl px-10 py-5 animate-pulse-glow inline-flex"
            >
              <Fish size={22} />
              Inscrever-se agora
            </Link>

            <p className="text-muted text-xs mt-6">
              Processo 100% digital. Contrato assinado eletronicamente. Dados protegidos pela LGPD.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
