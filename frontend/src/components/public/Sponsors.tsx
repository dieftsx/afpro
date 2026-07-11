'use client';
// Sponsors — Portal de Patrocinadores: captação, transparência financeira e contrapartidas (ROI).
import { Mail, TrendingUp, DollarSign, Shirt, Ship, Image as ImageIcon, Share2, Video, BarChart3 } from 'lucide-react';

const MODALIDADES = [
  {
    nivel: 'Patrocínio Master',
    cor: 'var(--color-secondary)',
    valor: 'R$ 50.000+',
    badge: '⭐ Master',
    beneficios: [
      'Logo na frente da camisa (destaque exclusivo)',
      'Logo na proa das embarcações',
      'Banner principal no evento',
      'Participação em vídeo institucional',
      'Menção em todas as publicações',
      'Co-branding em todos os materiais',
    ],
  },
  {
    nivel: 'Patrocínio Gold',
    cor: 'var(--color-primary)',
    valor: 'R$ 20.000+',
    badge: '🥇 Gold',
    beneficios: [
      'Logo na manga da camisa',
      'Logo nas embarcações (lateral)',
      'Banner lateral no evento',
      'Participação em Reels e stories',
      'Menção nos posts da temporada',
    ],
  },
  {
    nivel: 'Patrocínio Silver',
    cor: 'var(--color-accent)',
    valor: 'R$ 8.000+',
    badge: '🥈 Silver',
    beneficios: [
      'Logo no dorso da camisa',
      'Banner interno no evento',
      'Menção nas redes sociais',
      'Publicação de agradecimento',
    ],
  },
];

const CONTRAPARTIDAS = [
  {
    icon: Shirt,
    titulo: 'Uniformes',
    desc: 'Logo bordado nas camisetas, bonés e balaclavas de todas as participantes — visibilidade de marca durante e após o evento.',
  },
  {
    icon: Ship,
    titulo: 'Embarcações',
    desc: 'Adesivação nas lanchas e barcos utilizados nas temporadas, alcançando pescadoras, guias e fotógrafos.',
  },
  {
    icon: ImageIcon,
    titulo: 'Banners & Estrutura',
    desc: 'Banners e backdrops no ponto de pesca e área de convivência, capturados em fotos e vídeos do evento.',
  },
  {
    icon: Share2,
    titulo: 'Redes Sociais',
    desc: 'Menções nas publicações oficiais durante toda a temporada — stories, reels e posts no feed da AFPERO.',
  },
  {
    icon: Video,
    titulo: 'Vídeo Institucional',
    desc: 'Patrocinadores master e gold têm participação no vídeo de memória oficial da temporada, distribuído em multiplataformas.',
  },
  {
    icon: BarChart3,
    titulo: 'Relatório de Impacto',
    desc: 'Relatório personalizado ao final da temporada com alcance, engajamento e retorno de mídia para cada patrocinador.',
  },
];

export default function Sponsors() {
  return (
    <section id="patrocinadores" className="relative z-10 py-24 px-6">
      {/* Background decorativo */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: `radial-gradient(ellipse 60% 60% at 50% 100%, color-mix(in srgb, var(--color-secondary) 7%, transparent), transparent)`,
        }}
      />

      <div className="max-w-7xl mx-auto relative">

        {/* Header da seção */}
        <div className="text-center mb-16">
          <div className="divider-gradient mb-8" />
          <span className="badge badge-secondary mb-4 mx-auto inline-flex">
            <TrendingUp size={11} />
            Patrocínios & Parcerias
          </span>
          <h2 className="text-3xl md:text-5xl font-black mb-4 mt-2">
            <span>Portal do </span>
            <span className="gradient-text">Apoiador</span>
          </h2>
          <p className="text-muted max-w-2xl mx-auto leading-relaxed">
            Invista em um movimento que une mulheres, natureza e Rondônia. Sua marca ao lado de uma
            causa que transforma vidas e gera impacto real.
          </p>
        </div>

        {/* Bloco de Transparência Financeira */}
        <div className="glass p-8 md:p-10 mb-16 relative overflow-hidden">
          <div
            className="absolute inset-0 opacity-5 pointer-events-none"
            style={{ background: 'radial-gradient(circle at 80% 20%, var(--color-secondary), transparent 70%)' }}
          />
          <div className="relative grid md:grid-cols-2 gap-10 items-center">
            {/* Texto */}
            <div>
              <div className="flex items-center gap-3 mb-4">
                <div
                  className="w-10 h-10 rounded-xl flex items-center justify-center"
                  style={{ background: 'var(--color-secondary-dim)', color: 'var(--color-secondary)' }}
                >
                  <DollarSign size={18} />
                </div>
                <h3 className="font-black text-xl">Transparência Financeira</h3>
              </div>
              <p className="text-muted text-sm leading-relaxed mb-4">
                A AFPERO opera com total transparência. O custo global anual estimado para a
                realização das duas temporadas é divulgado publicamente para orientar patrocinadores,
                apoiadores e solicitações de emendas parlamentares.
              </p>
              <p className="text-muted text-sm leading-relaxed">
                Os recursos captados cobrem logística de campo, hospedagem das participantes,
                transporte fluvial e terrestre, infraestrutura de evento e a produção do
                kit personalizado.
              </p>
            </div>

            {/* Destaque do valor */}
            <div className="text-center">
              <p className="text-sm font-semibold text-muted uppercase tracking-widest mb-2">
                Orçamento Global Anual Estimado
              </p>
              <div
                className="text-5xl md:text-6xl font-black gradient-text mb-2"
                style={{ letterSpacing: '-0.03em' }}
              >
                R$ 280.000
              </div>
              <p className="text-xs text-muted mb-6">
                Duas temporadas · Logística · Hospedagem · Transporte · Estrutura
              </p>

              {/* Divisão por categoria */}
              <div className="space-y-2">
                {[
                  { label: 'Logística e Campo', pct: 35 },
                  { label: 'Hospedagem', pct: 25 },
                  { label: 'Transporte', pct: 20 },
                  { label: 'Estrutura e Kit', pct: 20 },
                ].map(({ label, pct }) => (
                  <div key={label}>
                    <div className="flex justify-between text-xs mb-1">
                      <span className="text-muted font-semibold">{label}</span>
                      <span style={{ color: 'var(--color-secondary)' }} className="font-bold">{pct}%</span>
                    </div>
                    <div className="progress-track">
                      <div className="progress-fill" style={{ width: `${pct}%` }} />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Grid de Contrapartidas (ROI para marcas) */}
        <div className="mb-16">
          <h3 className="text-2xl font-black text-center mb-10">
            O que sua <span className="gradient-text">marca recebe</span>
          </h3>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {CONTRAPARTIDAS.map(({ icon: Icon, titulo, desc }) => (
              <div
                key={titulo}
                className="glass p-6 group hover:scale-[1.02] transition-transform duration-300"
              >
                <div
                  className="w-11 h-11 rounded-xl flex items-center justify-center mb-4"
                  style={{ background: 'var(--color-primary-dim)', color: 'var(--color-primary)' }}
                >
                  <Icon size={20} />
                </div>
                <h4 className="font-bold text-base mb-2">{titulo}</h4>
                <p className="text-muted text-sm leading-relaxed">{desc}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Cards de Modalidades */}
        <div className="mb-16">
          <h3 className="text-2xl font-black text-center mb-10">
            Escolha a <span className="gradient-text">modalidade</span> ideal
          </h3>
          <div className="grid md:grid-cols-3 gap-6">
            {MODALIDADES.map((mod) => (
              <div
                key={mod.nivel}
                className="glass p-6 flex flex-col relative overflow-hidden"
                style={{ border: `1.5px solid color-mix(in srgb, ${mod.cor} 30%, transparent)` }}
              >
                <div
                  className="absolute top-0 left-0 right-0 h-1"
                  style={{ background: mod.cor }}
                />
                <p className="text-sm font-black mb-1">{mod.badge}</p>
                <h4 className="font-black text-lg mb-1" style={{ color: mod.cor }}>
                  {mod.nivel}
                </h4>
                <p className="text-2xl font-black mb-4">{mod.valor}</p>
                <ul className="space-y-2 flex-1">
                  {mod.beneficios.map((b) => (
                    <li key={b} className="flex items-start gap-2 text-sm text-muted">
                      <span style={{ color: mod.cor }} className="mt-0.5 shrink-0">✓</span>
                      {b}
                    </li>
                  ))}
                </ul>
                <a
                  href="mailto:contato@afpero.org.br?subject=Patrocínio AFPERO"
                  className="btn-outline mt-6 text-sm justify-center"
                  style={{ color: mod.cor, borderColor: mod.cor }}
                >
                  Quero Patrocinar
                </a>
              </div>
            ))}
          </div>
        </div>

        {/* Outras formas de apoio */}
        <div className="glass p-8 text-center relative overflow-hidden">
          <div
            className="absolute inset-0 opacity-5 pointer-events-none"
            style={{ background: 'radial-gradient(circle, var(--color-primary), transparent 70%)' }}
          />
          <div className="relative">
            <h3 className="text-2xl font-black mb-3">
              Outras formas de <span className="gradient-text">apoio</span>
            </h3>
            <p className="text-muted text-sm max-w-xl mx-auto mb-6 leading-relaxed">
              Além dos patrocínios diretos, a AFPERO aceita apoio via{' '}
              <strong className="font-semibold" style={{ color: 'var(--color-primary)' }}>emendas parlamentares</strong>,{' '}
              <strong className="font-semibold" style={{ color: 'var(--color-secondary)' }}>doações</strong> e{' '}
              <strong className="font-semibold" style={{ color: 'var(--color-accent)' }}>permutas de serviços</strong>.
              Entre em contato para conversarmos sobre o melhor formato para a sua contribuição.
            </p>
            <a
              href="mailto:contato@afpero.org.br?subject=Parceria AFPERO"
              id="sponsors-cta-contato"
              className="btn-primary inline-flex"
            >
              <Mail size={18} />
              Entrar em Contato
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
