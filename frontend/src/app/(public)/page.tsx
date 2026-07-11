// Landing Page Institucional — Server Component
// Compõe as seções modulares do site público da AFPERO.
// Compartilha o sistema de temas dinâmicos injetados pelo Root Layout (variáveis CSS globais).
import type { Metadata } from 'next';
import Hero from '@/components/public/Hero';
import About from '@/components/public/About';
import Seasons from '@/components/public/Seasons';
import Sponsors from '@/components/public/Sponsors';
import { Fish, ExternalLink, Mail } from 'lucide-react';
import Link from 'next/link';

export const metadata: Metadata = {
  title: 'AFPERO — Associação Feminina de Pesca Esportiva de Rondônia',
  description:
    'Conheça a AFPERO: pesca esportiva feminina, turismo sustentável e conservação ambiental no Rio Guaporé, Rondônia. Temporadas em Maio e Novembro.',
};

// ISR: revalida a cada 10 minutos (tema dinâmico vem do Root Layout)
export const revalidate = 600;

export default function InstitucionalPage() {
  return (
    <main className="relative min-h-screen overflow-x-hidden">
      {/* Background: orbes flutuantes globais */}
      <BackgroundFX />

      {/* Seção 1: Hero */}
      <Hero />

      {/* Seção 2: Sobre a Associação */}
      <About />

      {/* Divisor */}
      <div className="max-w-7xl mx-auto px-6">
        <div className="divider-gradient" />
      </div>

      {/* Seção 3: Temporadas de Pesca */}
      <Seasons />

      {/* Divisor */}
      <div className="max-w-7xl mx-auto px-6">
        <div className="divider-gradient" />
      </div>

      {/* Seção 4: Portal de Patrocinadores */}
      <Sponsors />

      {/* Footer institucional */}
      <footer className="relative z-10 py-12 px-6">
        <div className="max-w-7xl mx-auto">
          <div className="divider-gradient mb-10" />
          <div className="grid sm:grid-cols-3 gap-8 mb-10">
            {/* Brand */}
            <div>
              <div className="flex items-center gap-3 mb-4">
                <div
                  className="w-10 h-10 rounded-xl flex items-center justify-center"
                  style={{
                    background: 'linear-gradient(135deg, var(--color-primary), var(--color-secondary))',
                    color: 'var(--color-bg)',
                  }}
                >
                  <Fish size={18} />
                </div>
                <span className="font-black text-sm tracking-widest uppercase gradient-text">AFPERO</span>
              </div>
              <p className="text-muted text-xs leading-relaxed">
                Associação Feminina de Pesca Esportiva de Rondônia. Fundada em 2024 para
                formalizar e expandir o programa Meninas na Pesca.
              </p>
            </div>

            {/* Links rápidos */}
            <div>
              <p className="font-black text-xs uppercase tracking-widest text-muted mb-4">Navegação</p>
              <nav className="space-y-2">
                {[
                  { href: '#sobre', label: 'Sobre a Associação' },
                  { href: '#temporadas', label: 'Temporadas' },
                  { href: '#patrocinadores', label: 'Patrocínios' },
                  { href: '/inscricao', label: 'Inscrições' },
                ].map((link) => (
                  <a
                    key={link.href}
                    href={link.href}
                    className="block text-sm text-muted hover:text-primary transition-colors duration-200"
                  >
                    {link.label}
                  </a>
                ))}
              </nav>
            </div>

            {/* Contato */}
            <div>
              <p className="font-black text-xs uppercase tracking-widest text-muted mb-4">Contato</p>
              <div className="space-y-3">
                <a
                  href="mailto:contato@afpero.org.br"
                  className="flex items-center gap-2 text-sm text-muted hover:text-primary transition-colors"
                >
                  <Mail size={13} />
                  contato@afpero.org.br
                </a>
                <a
                  href="https://instagram.com/afpero"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-sm text-muted hover:text-primary transition-colors"
                >
                  <ExternalLink size={13} />
                  @afpero
                </a>
                <Link href="/inscricao" className="btn-primary text-xs px-4 py-2 mt-2 inline-flex">
                  Inscreva-se agora
                </Link>
              </div>
            </div>
          </div>

          {/* Copyright */}
          <div className="divider-gradient mb-6" />
          <div className="flex flex-col sm:flex-row items-center justify-between gap-2 text-xs text-muted opacity-70">
            <p>© {new Date().getFullYear()} AFPERO — Associação Feminina de Pesca Esportiva de Rondônia</p>
            <p>Todos os direitos reservados · CNPJ em processo de registro</p>
          </div>
        </div>
      </footer>
    </main>
  );
}

// ─── Background FX global ─────────────────────────────────────────────────────
function BackgroundFX() {
  return (
    <div className="fixed inset-0 pointer-events-none z-0 overflow-hidden">
      {/* Gradiente de fundo */}
      <div
        className="absolute inset-0"
        style={{
          background: `
            radial-gradient(ellipse 80% 60% at 20% 20%, var(--color-primary-dim) 0%, transparent 60%),
            radial-gradient(ellipse 60% 50% at 80% 80%, var(--color-secondary-dim) 0%, transparent 60%)
          `,
        }}
      />
      {/* Orbe primário flutuante */}
      <div
        className="absolute w-[500px] h-[500px] rounded-full animate-float opacity-10 blur-3xl"
        style={{
          top: '30%',
          left: '-10%',
          background: 'var(--color-primary)',
          animationDuration: '10s',
        }}
      />
      {/* Orbe secundário */}
      <div
        className="absolute w-[350px] h-[350px] rounded-full animate-float opacity-10 blur-3xl"
        style={{
          bottom: '20%',
          right: '-5%',
          background: 'var(--color-secondary)',
          animationDuration: '12s',
          animationDelay: '-5s',
        }}
      />
      {/* Grid sutil */}
      <div
        className="absolute inset-0 opacity-[0.025]"
        style={{
          backgroundImage: `
            linear-gradient(var(--color-text) 1px, transparent 1px),
            linear-gradient(90deg, var(--color-text) 1px, transparent 1px)
          `,
          backgroundSize: '60px 60px',
        }}
      />
    </div>
  );
}
