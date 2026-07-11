'use client';
// Seasons — Temporadas de pesca: timeline anual + destaque da edição atual.
import Link from 'next/link';
import { Calendar, MapPin, Fish, Star, ChevronRight } from 'lucide-react';

const TEMPORADAS_HISTORICO = [
  { num: '1ª', periodo: 'Maio 2024', vagas: 15, status: 'Encerrada' },
  { num: '2ª', periodo: 'Novembro 2024', vagas: 20, status: 'Encerrada' },
  { num: '3ª', periodo: 'Maio 2025', vagas: 25, status: 'Encerrada' },
  { num: '4ª', periodo: 'Novembro 2025', vagas: 30, status: 'Encerrada' },
];

export default function Seasons() {
  return (
    <section id="temporadas" className="relative z-10 py-24 px-6">
      {/* Background decorativo */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: `radial-gradient(ellipse 70% 50% at 0% 50%, color-mix(in srgb, var(--color-primary) 6%, transparent), transparent)`,
        }}
      />

      <div className="max-w-7xl mx-auto relative">
        {/* Header da seção */}
        <div className="text-center mb-16">
          <div className="divider-gradient mb-8" />
          <span className="badge badge-primary mb-4 mx-auto inline-flex">
            <Calendar size={11} />
            Calendário de Eventos
          </span>
          <h2 className="text-3xl md:text-5xl font-black mb-4 mt-2">
            <span>Nossas </span>
            <span className="gradient-text">Temporadas</span>
          </h2>
          <p className="text-muted max-w-xl mx-auto leading-relaxed">
            Duas temporadas por ano — sempre em Maio e Novembro — no Distrito de Porto
            Rolim do Guaporé, às margens do Rio Guaporé, em Rondônia.
          </p>
        </div>

        {/* Local fixo */}
        <div className="flex justify-center mb-16">
          <div className="glass px-8 py-4 flex items-center gap-4 max-w-xl w-full">
            <div
              className="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
              style={{ background: 'var(--color-primary-dim)', color: 'var(--color-primary)' }}
            >
              <MapPin size={18} />
            </div>
            <div>
              <p className="font-black text-sm">Local Fixo de Todas as Temporadas</p>
              <p className="text-muted text-xs mt-0.5">
                Distrito de Porto Rolim do Guaporé · Rio Guaporé · Rondônia (RO)
              </p>
            </div>
          </div>
        </div>

        <div className="grid lg:grid-cols-5 gap-8 items-start">
          {/* Timeline histórica — coluna esquerda (2 cols) */}
          <div className="lg:col-span-2">
            <h3 className="font-black text-lg mb-6 text-muted uppercase tracking-widest text-sm">
              Histórico
            </h3>
            <div className="relative">
              {/* Linha vertical */}
              <div
                className="absolute left-[19px] top-0 bottom-0 w-px"
                style={{ background: 'var(--color-border)' }}
              />
              <div className="space-y-4">
                {TEMPORADAS_HISTORICO.map((t, i) => (
                  <div key={i} className="flex items-center gap-4 relative pl-2">
                    {/* Dot */}
                    <div
                      className="w-9 h-9 rounded-full flex items-center justify-center text-xs font-black shrink-0 z-10"
                      style={{
                        background: 'var(--color-surface-high)',
                        border: '2px solid var(--color-border)',
                        color: 'var(--color-text-muted)',
                      }}
                    >
                      {i + 1}
                    </div>
                    <div className="glass-sm flex-1 px-4 py-3 flex items-center justify-between">
                      <div>
                        <p className="font-bold text-sm">{t.num} Temporada</p>
                        <p className="text-xs text-muted">{t.periodo} · {t.vagas} vagas</p>
                      </div>
                      <span
                        className="text-xs px-2 py-0.5 rounded-full font-semibold"
                        style={{
                          background: 'var(--color-surface-high)',
                          color: 'var(--color-text-muted)',
                        }}
                      >
                        {t.status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Próxima janela */}
            <div className="mt-6 glass-sm px-4 py-3 flex items-center gap-3">
              <Calendar size={14} style={{ color: 'var(--color-secondary)' }} />
              <p className="text-xs text-muted">
                <span className="font-bold" style={{ color: 'var(--color-secondary)' }}>
                  Próxima janela:
                </span>{' '}
                Maio de 2027
              </p>
            </div>
          </div>

          {/* Card destaque da edição atual — coluna direita (3 cols) */}
          <div className="lg:col-span-3">
            <h3 className="font-black text-lg mb-6 uppercase tracking-widest text-sm gradient-text">
              Edição Atual
            </h3>

            <div
              className="glass relative overflow-hidden rounded-2xl"
              style={{
                border: '2px solid var(--color-primary)',
                boxShadow: '0 0 60px var(--color-primary-dim)',
              }}
            >
              {/* Glow interior */}
              <div
                className="absolute inset-0 opacity-10 pointer-events-none"
                style={{
                  background:
                    'radial-gradient(ellipse 80% 60% at 50% 0%, var(--color-primary), transparent)',
                }}
              />

              {/* Conteúdo */}
              <div className="relative p-8 md:p-10">
                {/* Badge destaque */}
                <div className="flex items-center gap-3 mb-6">
                  <span
                    className="badge text-xs px-4 py-2 animate-pulse-glow"
                    style={{
                      background: 'var(--color-primary)',
                      color: 'var(--color-bg)',
                      border: 'none',
                    }}
                  >
                    <Star size={11} />
                    INSCRIÇÕES ABERTAS
                  </span>
                </div>

                <div className="mb-8">
                  <h3 className="text-4xl md:text-5xl font-black gradient-text mb-2">
                    5ª Temporada
                  </h3>
                  <p className="text-xl font-extrabold" style={{ color: 'var(--color-text)' }}>
                    Meninas na Pesca
                  </p>
                  <p className="text-muted mt-1">Novembro de 2026</p>
                </div>

                {/* Detalhes */}
                <div className="grid sm:grid-cols-3 gap-4 mb-8">
                  {[
                    { icon: Calendar, label: 'Período', value: 'Novembro / 2026' },
                    { icon: MapPin, label: 'Local', value: 'Porto Rolim do Guaporé' },
                    { icon: Fish, label: 'Vagas', value: '30 participantes' },
                  ].map(({ icon: Icon, label, value }) => (
                    <div
                      key={label}
                      className="glass-sm p-4 text-center"
                    >
                      <Icon size={16} className="mx-auto mb-2" style={{ color: 'var(--color-primary)' }} />
                      <p className="text-xs text-muted uppercase tracking-wider mb-1">{label}</p>
                      <p className="font-bold text-sm">{value}</p>
                    </div>
                  ))}
                </div>

                {/* Incluso */}
                <div className="space-y-2 mb-8">
                  {[
                    '🏕️ Hospedagem e estrutura completa',
                    '🚤 Embarcações equipadas',
                    '👕 Kit personalizado (camisa, boné, balaclava)',
                    '🎓 Monitores especializados',
                    '🐟 Pesca esportiva (pesque e solte)',
                  ].map((item) => (
                    <div key={item} className="flex items-center gap-2 text-sm text-muted">
                      <ChevronRight size={13} style={{ color: 'var(--color-primary)' }} />
                      {item}
                    </div>
                  ))}
                </div>

                {/* CTA */}
                <Link
                  href="/inscricao"
                  id="seasons-cta-garantir-vaga"
                  className="btn-primary w-full justify-center text-lg py-4 animate-pulse-glow"
                >
                  <Fish size={20} />
                  Garantir Minha Vaga
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
