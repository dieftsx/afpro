'use client';
// About — História e Propósito da AFPERO com layout split (texto + visual).
import { Heart, Users, Leaf, TrendingUp, CheckCircle } from 'lucide-react';

const PUBLICS = [
  { emoji: '🎣', label: 'Pescadoras experientes' },
  { emoji: '🌱', label: 'Iniciantes curiosas' },
  { emoji: '💼', label: 'Empreendedoras locais' },
  { emoji: '👨‍👩‍👧', label: 'Famílias' },
  { emoji: '📸', label: 'Influenciadoras digitais' },
];

const VALUES = [
  {
    icon: Heart,
    title: 'Inclusão Feminina',
    desc: 'Espaço exclusivo e acolhedor para mulheres de todos os perfis vivenciarem a pesca esportiva sem barreiras.',
    color: 'var(--color-primary)',
  },
  {
    icon: TrendingUp,
    title: 'Impacto Econômico Local',
    desc: 'Cada temporada movimenta a economia de Porto Rolim do Guaporé, gerando renda para guias, hospedagens e comércio.',
    color: 'var(--color-secondary)',
  },
  {
    icon: Leaf,
    title: 'Preservação Ambiental',
    desc: 'Praticamos o sistema pesque e solte, garantindo a sustentabilidade do ecossistema do Rio Guaporé para as próximas gerações.',
    color: 'var(--color-accent)',
  },
];

export default function About() {
  return (
    <section id="sobre" className="relative z-10 py-24 px-6">
      {/* Background decorativo */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: `radial-gradient(ellipse 60% 50% at 100% 50%, color-mix(in srgb, var(--color-secondary) 8%, transparent), transparent)`,
        }}
      />

      <div className="max-w-7xl mx-auto relative">
        {/* Header da seção */}
        <div className="text-center mb-16">
          <div className="divider-gradient mb-8" />
          <span className="badge badge-secondary mb-4 mx-auto inline-flex">
            <Users size={11} />
            Sobre a Associação
          </span>
          <h2 className="text-3xl md:text-5xl font-black mb-4 mt-2">
            <span className="gradient-text">Uma história</span>
            <br />
            <span>de propósito e crescimento</span>
          </h2>
          <p className="text-muted max-w-2xl mx-auto leading-relaxed">
            Nascida da paixão pela pesca e pela força feminina, a AFPERO foi fundada para
            formalizar e expandir um movimento que já transformava vidas desde 2024.
          </p>
        </div>

        {/* Layout split: Narrative + Timeline */}
        <div className="grid md:grid-cols-2 gap-12 items-start mb-20">
          {/* Coluna esquerda: Narrativa */}
          <div className="space-y-6">
            <div className="glass p-8">
              <h3 className="text-xl font-black mb-4 gradient-text">A origem do programa</h3>
              <p className="text-muted leading-relaxed text-sm mb-4">
                O programa <strong className="text-primary font-semibold">Meninas na Pesca</strong> surgiu
                em 2024 como uma iniciativa pioneira no estado de Rondônia — um espaço criado por
                mulheres, para mulheres, no coração do Rio Guaporé.
              </p>
              <p className="text-muted leading-relaxed text-sm mb-4">
                O crescimento expressivo das temporadas e o impacto positivo na vida das participantes
                revelaram a necessidade de uma estrutura jurídica e institucional sólida. Foi assim que
                nasceu a <strong className="text-secondary font-semibold">AFPERO</strong> — para expandir
                as ações, viabilizar recursos públicos e privados, e garantir a continuidade do programa.
              </p>
              <p className="text-muted leading-relaxed text-sm">
                Hoje, a associação é o elo entre a paixão pela pesca, o desenvolvimento feminino,
                o turismo regional e a responsabilidade ambiental.
              </p>
            </div>

            {/* Público-alvo */}
            <div className="glass p-6">
              <h3 className="font-black text-base mb-4 flex items-center gap-2">
                <Users size={16} style={{ color: 'var(--color-primary)' }} />
                Quem participa?
              </h3>
              <div className="grid grid-cols-1 gap-2">
                {PUBLICS.map(({ emoji, label }) => (
                  <div key={label} className="flex items-center gap-3 glass-sm px-4 py-2.5">
                    <span className="text-lg">{emoji}</span>
                    <span className="text-sm font-semibold">{label}</span>
                    <CheckCircle
                      size={14}
                      className="ml-auto shrink-0"
                      style={{ color: 'var(--color-primary)' }}
                    />
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Coluna direita: Valores em cards */}
          <div className="space-y-6">
            {VALUES.map(({ icon: Icon, title, desc, color }) => (
              <div
                key={title}
                className="glass p-6 group hover:scale-[1.02] transition-transform duration-300 relative overflow-hidden"
              >
                {/* Glow de fundo no hover */}
                <div
                  className="absolute -inset-1 opacity-0 group-hover:opacity-10 blur-xl transition-opacity duration-500 pointer-events-none"
                  style={{ background: color }}
                />
                <div className="relative flex gap-5 items-start">
                  <div
                    className="w-12 h-12 rounded-2xl flex items-center justify-center shrink-0 shadow-lg"
                    style={{ background: `color-mix(in srgb, ${color} 20%, transparent)`, color }}
                  >
                    <Icon size={22} />
                  </div>
                  <div>
                    <h3 className="font-black text-base mb-1.5" style={{ color }}>
                      {title}
                    </h3>
                    <p className="text-muted text-sm leading-relaxed">{desc}</p>
                  </div>
                </div>
              </div>
            ))}

            {/* Detalhe: pesque e solte */}
            <div
              className="glass-sm p-5 border-l-4 flex items-start gap-4"
              style={{ borderColor: 'var(--color-primary)' }}
            >
              <Leaf size={20} style={{ color: 'var(--color-primary)' }} className="shrink-0 mt-0.5" />
              <div>
                <p className="font-bold text-sm mb-1" style={{ color: 'var(--color-primary)' }}>
                  🐟 Política de Pesque e Solte
                </p>
                <p className="text-muted text-xs leading-relaxed">
                  Toda pesca praticada nas temporadas da AFPERO segue o protocolo de captura e
                  liberação responsável, assegurando a saúde do ecossistema do Rio Guaporé.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
