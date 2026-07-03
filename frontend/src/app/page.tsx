// Landing Page — Server Component
// Busca os dados da temporada ativa e renderiza a página de captura.
import Link from 'next/link';
import { getTemporadaAtiva } from '@/lib/api';
import HeroSection from '@/components/landing/HeroSection';
import InfoSection from '@/components/landing/InfoSection';
import CTASection from '@/components/landing/CTASection';

export const revalidate = 300; // ISR: revalida a cada 5 minutos

export default async function HomePage() {
  let temporada = null;

  try {
    temporada = await getTemporadaAtiva();
  } catch {
    // Em caso de erro na API, renderiza com dados de fallback
    temporada = {
      id: '5a-temporada-2025',
      nome: '5ª Temporada Meninas na Pesca',
      edicao: 5,
      slug: '5a-temporada-2025',
      data_inicio_evento: '2025-09-01T00:00:00Z',
      data_fim_evento: '2025-09-07T00:00:00Z',
      data_abertura_inscricoes: '2025-07-01T00:00:00Z',
      data_fechamento_inscricoes: '2025-08-20T00:00:00Z',
      vagas_total: 30,
      vagas_disponiveis: 30,
      valor_total: 3500.00,
      valor_entrada_pix: 800.00,
      valor_parcela: 540.00,
      qtd_parcelas: 5,
      cor_primaria: 'hsl(168 72% 40%)',
      cor_secundaria: 'hsl(45 95% 55%)',
      cor_acento: 'hsl(168 80% 70%)',
      cor_fundo: 'hsl(220 20% 10%)',
      cor_texto: 'hsl(0 0% 97%)',
      logo_url: null,
      banner_url: null,
    };
  }

  const vagasDisponiveis = temporada?.vagas_disponiveis ?? 0;
  const inscricoesAbertas = vagasDisponiveis > 0;

  return (
    <main className="relative min-h-screen overflow-x-hidden">
      {/* Background: partículas e gradiente */}
      <BackgroundFX />

      {/* Navegação */}
      <nav className="fixed top-0 left-0 right-0 z-50 px-6 py-4">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div
              className="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold"
              style={{ background: 'var(--color-primary)', color: 'var(--color-bg)' }}
            >
              AF
            </div>
            <span className="font-bold text-sm tracking-widest uppercase opacity-80">
              AFPERO
            </span>
          </div>

          {inscricoesAbertas && (
            <Link
              href="/inscricao"
              className="btn-primary text-sm px-5 py-2.5"
            >
              Inscreva-se
            </Link>
          )}
        </div>
      </nav>

      {/* Hero */}
      <HeroSection temporada={temporada} inscricoesAbertas={inscricoesAbertas} />

      {/* Informações */}
      <InfoSection temporada={temporada} />

      {/* CTA */}
      <CTASection temporada={temporada} inscricoesAbertas={inscricoesAbertas} />

      {/* Footer */}
      <footer className="py-8 text-center">
        <div className="divider-gradient mb-8" />
        <p className="text-muted text-sm">
          © {new Date().getFullYear()} AFPERO — Associação Feminina de Pesca Esportiva de Rondônia
        </p>
        <p className="text-muted text-xs mt-1 opacity-60">
          Todos os direitos reservados.
        </p>
      </footer>
    </main>
  );
}

// ─── Background animado ────────────────────────────────────────────────────────
function BackgroundFX() {
  return (
    <div className="fixed inset-0 pointer-events-none z-0 overflow-hidden">
      {/* Gradiente radial principal */}
      <div
        className="absolute inset-0"
        style={{
          background: `
            radial-gradient(ellipse 80% 60% at 20% 20%, var(--color-primary-dim) 0%, transparent 60%),
            radial-gradient(ellipse 60% 50% at 80% 80%, var(--color-secondary-dim) 0%, transparent 60%),
            radial-gradient(ellipse 40% 40% at 50% 50%, var(--color-bg) 0%, transparent 100%)
          `,
        }}
      />

      {/* Orbes flutuantes */}
      <div
        className="absolute w-96 h-96 rounded-full animate-float opacity-20 blur-3xl"
        style={{
          top: '10%',
          left: '-5%',
          background: 'var(--color-primary)',
          animationDuration: '7s',
        }}
      />
      <div
        className="absolute w-64 h-64 rounded-full animate-float opacity-15 blur-3xl"
        style={{
          bottom: '15%',
          right: '-5%',
          background: 'var(--color-secondary)',
          animationDuration: '9s',
          animationDelay: '-3s',
        }}
      />

      {/* Grid sutil */}
      <div
        className="absolute inset-0 opacity-[0.03]"
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
