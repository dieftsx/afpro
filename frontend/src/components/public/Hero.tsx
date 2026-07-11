'use client';
// Hero — Seção de alto impacto visual com CTAs institucionais.
import Link from 'next/link';
import Image from 'next/image';
import { Fish, Leaf, Compass, ChevronDown } from 'lucide-react';

const PILLARS = [
  { icon: Fish, label: 'Pesca Esportiva' },
  { icon: Leaf, label: 'Conservação Ambiental' },
  { icon: Compass, label: 'Turismo Sustentável' },
];

export default function Hero() {
  return (
    <section
      id="inicio"
      className="relative min-h-screen flex flex-col items-center justify-center overflow-hidden"
    >
      {/* Background: imagem + overlay gradiente */}
      <div className="absolute inset-0 z-0">
        <Image
          src="/hero-fishing.png"
          alt="Mulheres pescando no Rio Guaporé ao entardecer em Rondônia"
          fill
          priority
          className="object-cover object-center"
          sizes="100vw"
        />
        {/* Overlay escuro degradado */}
        <div
          className="absolute inset-0"
          style={{
            background: `
              linear-gradient(180deg,
                color-mix(in srgb, var(--color-bg) 70%, transparent) 0%,
                color-mix(in srgb, var(--color-bg) 40%, transparent) 40%,
                color-mix(in srgb, var(--color-bg) 85%, transparent) 100%
              )
            `,
          }}
        />
        {/* Glow primário no canto superior esquerdo */}
        <div
          className="absolute top-0 left-0 w-[600px] h-[600px] rounded-full blur-3xl opacity-25 pointer-events-none"
          style={{ background: 'var(--color-primary)' }}
        />
        {/* Glow dourado no canto inferior direito */}
        <div
          className="absolute bottom-0 right-0 w-[400px] h-[400px] rounded-full blur-3xl opacity-20 pointer-events-none"
          style={{ background: 'var(--color-secondary)' }}
        />
      </div>

      {/* Conteúdo principal */}
      <div className="relative z-10 max-w-5xl mx-auto px-6 pt-28 pb-24 text-center flex flex-col items-center">

        {/* Badge institucional */}
        <div className="flex justify-center mb-8 animate-fade-in-up">
          <span className="badge badge-primary text-xs px-4 py-1.5 tracking-widest">
            <Fish size={11} />
            AFPERO · Fundada em 2024 · Rondônia
          </span>
        </div>

        {/* Título principal */}
        <h1 className="text-5xl sm:text-6xl md:text-8xl font-black leading-none tracking-tight mb-6 animate-fade-in-up delay-100">
          <span className="gradient-text glow-text-primary">Pesca Esportiva</span>
          <br />
          <span style={{ color: 'var(--color-text)' }}>Feminina</span>
          <br />
          <span
            className="text-3xl sm:text-4xl md:text-5xl font-extrabold"
            style={{ color: 'var(--color-text-muted)' }}
          >
            com Propósito
          </span>
        </h1>

        {/* Subtítulo */}
        <p
          className="text-base sm:text-lg md:text-xl max-w-2xl mx-auto mb-10 leading-relaxed animate-fade-in-up delay-200"
          style={{ color: 'var(--color-text-muted)' }}
        >
          A AFPERO conecta mulheres à natureza do Rio Guaporé em temporadas de pesca
          esportiva que transformam vidas, fomentam o turismo local e defendem o meio ambiente.
        </p>

        {/* Pilares */}
        <div className="flex flex-wrap justify-center gap-3 mb-12 animate-fade-in-up delay-300">
          {PILLARS.map(({ icon: Icon, label }) => (
            <div key={label} className="glass-sm flex items-center gap-2 px-4 py-2 text-sm font-semibold">
              <Icon size={15} style={{ color: 'var(--color-primary)' }} />
              <span>{label}</span>
            </div>
          ))}
        </div>

        {/* CTAs */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center animate-fade-in-up delay-400">
          <a
            href="#temporadas"
            id="hero-cta-temporadas"
            className="btn-primary text-base px-8 py-4 animate-pulse-glow"
          >
            <Fish size={18} />
            Nossas Temporadas
          </a>
          <a
            href="#patrocinadores"
            id="hero-cta-apoiadores"
            className="btn-outline text-base px-8 py-4"
          >
            Seja um Apoiador
          </a>
        </div>
      </div>

      {/* Scroll indicator */}
      <div className="absolute bottom-10 left-1/2 -translate-x-1/2 z-10 animate-float opacity-50">
        <div
          className="flex flex-col items-center gap-1"
          style={{ color: 'var(--color-primary)' }}
        >
          <span className="text-xs font-semibold tracking-widest uppercase">Explore</span>
          <ChevronDown size={20} />
        </div>
      </div>
    </section>
  );
}
