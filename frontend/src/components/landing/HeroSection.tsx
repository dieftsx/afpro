'use client';
// HeroSection — seção principal da landing page com contador de vagas e CTA.
import Link from 'next/link';
import { Fish, Calendar, MapPin, Users } from 'lucide-react';
import type { Temporada } from '@/types';

interface HeroSectionProps {
  temporada: Temporada | null;
  inscricoesAbertas: boolean;
}

export default function HeroSection({ temporada, inscricoesAbertas }: HeroSectionProps) {
  const vagasPorcentagem = temporada
    ? ((temporada.vagas_total - temporada.vagas_disponiveis) / temporada.vagas_total) * 100
    : 0;

  return (
    <section className="relative z-10 min-h-screen flex flex-col items-center justify-center px-6 pt-24 pb-16">
      <div className="max-w-4xl mx-auto text-center">

        {/* Badge de edição */}
        <div className="flex justify-center mb-8 animate-fade-in-up">
          <span className="badge badge-primary">
            <Fish size={12} />
            {temporada?.nome ?? '5ª Temporada Meninas na Pesca'}
          </span>
        </div>

        {/* Título principal */}
        <h1 className="text-5xl md:text-7xl font-black tracking-tight mb-6 animate-fade-in-up delay-100 glow-text-primary">
          <span className="gradient-text">Meninas</span>
          <br />
          <span style={{ color: 'var(--color-text)' }}>na Pesca</span>
        </h1>

        {/* Subtítulo */}
        <p className="text-lg md:text-xl text-muted max-w-2xl mx-auto mb-10 animate-fade-in-up delay-200 leading-relaxed">
          A experiência de pesca esportiva mais especial de Rondônia.
          Uma temporada criada por mulheres, para mulheres — na natureza, com propósito.
        </p>

        {/* Stats rápidos */}
        {temporada && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10 animate-fade-in-up delay-300">
            {[
              {
                icon: <Calendar size={18} />,
                label: 'Evento',
                value: formatDate(temporada.data_inicio_evento),
              },
              {
                icon: <MapPin size={18} />,
                label: 'Local',
                value: 'Rondônia',
              },
              {
                icon: <Users size={18} />,
                label: 'Vagas',
                value: `${temporada.vagas_disponiveis} restantes`,
              },
              {
                icon: <Fish size={18} />,
                label: 'Investimento',
                value: formatBRL(temporada.valor_total),
              },
            ].map((stat, i) => (
              <div key={i} className="glass-sm p-4 text-center">
                <div className="text-primary flex justify-center mb-1">{stat.icon}</div>
                <div className="text-xs text-muted mb-0.5 uppercase tracking-wider">{stat.label}</div>
                <div className="font-bold text-sm">{stat.value}</div>
              </div>
            ))}
          </div>
        )}

        {/* Barra de ocupação de vagas */}
        {temporada && (
          <div className="mb-8 animate-fade-in-up delay-300 max-w-sm mx-auto">
            <div className="flex justify-between text-xs text-muted mb-2">
              <span>{temporada.vagas_total - temporada.vagas_disponiveis} inscritas</span>
              <span>{temporada.vagas_disponiveis} vagas restantes</span>
            </div>
            <div className="progress-track">
              <div
                className="progress-fill"
                style={{ width: `${vagasPorcentagem}%` }}
              />
            </div>
          </div>
        )}

        {/* CTAs */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center animate-fade-in-up delay-400">
          {inscricoesAbertas ? (
            <Link href="/inscricao" className="btn-primary text-lg px-8 py-4 animate-pulse-glow">
              <Fish size={20} />
              Garantir minha vaga
            </Link>
          ) : (
            <div className="glass-sm px-8 py-4 text-center">
              <span className="text-muted font-semibold">
                {temporada ? 'Inscrições encerradas' : 'Carregando...'}
              </span>
            </div>
          )}

          <a
            href="#informacoes"
            className="btn-outline text-lg px-8 py-4"
          >
            Saiba mais
          </a>
        </div>
      </div>

      {/* Scroll indicator */}
      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 animate-float opacity-40">
        <div
          className="w-6 h-10 rounded-full border-2 flex items-start justify-center pt-2"
          style={{ borderColor: 'var(--color-primary)' }}
        >
          <div
            className="w-1.5 h-2.5 rounded-full animate-float"
            style={{ background: 'var(--color-primary)', animationDuration: '1.5s' }}
          />
        </div>
      </div>
    </section>
  );
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr);
  return d.toLocaleDateString('pt-BR', { day: '2-digit', month: 'long' });
}

function formatBRL(value: number) {
  return value.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
}
