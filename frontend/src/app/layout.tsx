// Root Layout — Server Component
// Injeta as CSS variables do tema diretamente no <html> em tempo de renderização.
// Isso garante que o tema seja aplicado sem nenhum flash de conteúdo sem estilo (FOUC).
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import { getTemporadaAtiva } from '@/lib/api';
import { defaultTheme } from '@/lib/theme';

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
});

export const metadata: Metadata = {
  title: 'afpro — Associação Feminina de Pesca Esportiva de Rondônia',
  description:
    'Inscreva-se na 5ª Temporada Meninas na Pesca. A maior experiência de pesca esportiva feminina de Rondônia.',
  keywords: ['pesca esportiva', 'mulheres', 'rondônia', 'afpro', 'temporada'],
  openGraph: {
    title: 'afpro — 5ª Temporada Meninas na Pesca',
    description: 'Garanta sua vaga na maior experiência de pesca esportiva feminina de Rondônia.',
    type: 'website',
  },
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // Busca o tema da temporada ativa do backend Go.
  // Em caso de falha (API offline, etc.), usa o tema padrão da 5ª Temporada.
  let themeVars = defaultTheme;

  try {
    const temporada = await getTemporadaAtiva();
    themeVars = {
      '--color-primary': temporada.cor_primaria,
      '--color-secondary': temporada.cor_secundaria,
      '--color-accent': temporada.cor_acento,
      '--color-bg': temporada.cor_fundo,
      '--color-text': temporada.cor_texto,
    };
  } catch {
    // Silencioso — fallback para defaultTheme já configurado acima
  }

  const htmlStyle = {
    fontFamily: 'var(--font-inter), Inter, system-ui, sans-serif',
    ...themeVars,
  } as React.CSSProperties;

  return (
    <html
      lang="pt-BR"
      className={inter.variable}
      style={htmlStyle}
    >
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="" />
      </head>
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
