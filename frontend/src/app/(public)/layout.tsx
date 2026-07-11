// Layout do Route Group (public) — Site Institucional
// Herda as variáveis de tema do Root Layout e adiciona a navbar pública.
import type { Metadata } from 'next';
import Link from 'next/link';
import { Fish } from 'lucide-react';

export const metadata: Metadata = {
  title: {
    template: '%s | AFPERO — Associação Feminina de Pesca Esportiva de Rondônia',
    default: 'AFPERO — Associação Feminina de Pesca Esportiva de Rondônia',
  },
  description:
    'A AFPERO promove a inclusão feminina na pesca esportiva, o turismo sustentável e a conservação ambiental no Rio Guaporé, Rondônia.',
  keywords: [
    'pesca esportiva feminina',
    'afpero',
    'rondônia',
    'rio guaporé',
    'conservação ambiental',
    'turismo sustentável',
    'meninas na pesca',
  ],
  openGraph: {
    type: 'website',
    siteName: 'AFPERO',
  },
};

export default function PublicLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      {/* Navbar pública fixa */}
      <header
        className="fixed top-0 left-0 right-0 z-50"
        style={{ background: 'color-mix(in srgb, var(--color-bg) 80%, transparent)' }}
      >
        <div
          className="absolute inset-x-0 bottom-0 h-px"
          style={{ background: 'var(--color-border)' }}
        />
        <nav className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-3 group" aria-label="AFPERO - Página Inicial">
            <div
              className="w-10 h-10 rounded-xl flex items-center justify-center shadow-lg transition-transform group-hover:scale-105"
              style={{
                background: 'linear-gradient(135deg, var(--color-primary), var(--color-secondary))',
                color: 'var(--color-bg)',
              }}
            >
              <Fish size={18} />
            </div>
            <div className="flex flex-col leading-tight">
              <span className="font-black text-sm tracking-widest uppercase gradient-text">AFPERO</span>
              <span className="text-xs hidden sm:block" style={{ color: 'var(--color-text-muted)' }}>
                Pesca Esportiva Feminina
              </span>
            </div>
          </Link>

          {/* Links de navegação */}
          <div className="hidden md:flex items-center gap-8">
            {[
              { href: '#sobre', label: 'Sobre' },
              { href: '#temporadas', label: 'Temporadas' },
              { href: '#patrocinadores', label: 'Patrocinadores' },
            ].map((link) => (
              <a
                key={link.href}
                href={link.href}
                className="text-sm font-semibold transition-colors duration-200 hover:text-primary"
                style={{ color: 'var(--color-text-muted)' }}
              >
                {link.label}
              </a>
            ))}
          </div>

          {/* CTA Navbar */}
          <Link
            href="/inscricao"
            className="btn-primary text-sm px-5 py-2.5"
            id="nav-cta-inscricao"
          >
            Garantir Vaga
          </Link>
        </nav>
      </header>

      {children}
    </>
  );
}
