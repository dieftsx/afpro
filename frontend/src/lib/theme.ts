// Engine de Temas Dinâmicos — converte cores da temporada em CSS variables.
// As variáveis são injetadas no <html> via root layout em tempo de execução.
import type { Temporada } from '@/types';

export interface ThemeVars {
  '--color-primary': string;
  '--color-secondary': string;
  '--color-accent': string;
  '--color-bg': string;
  '--color-text': string;
}

/**
 * Converte os campos de cor de uma Temporada para um objeto de CSS variables.
 * Esses valores são aplicados como style no elemento <html> pelo root layout.
 */
export function buildThemeVars(temporada: Temporada): ThemeVars {
  return {
    '--color-primary': temporada.cor_primaria,
    '--color-secondary': temporada.cor_secundaria,
    '--color-accent': temporada.cor_acento,
    '--color-bg': temporada.cor_fundo,
    '--color-text': temporada.cor_texto,
  };
}

/**
 * Converte as CSS variables para uma string inline (style attribute).
 * Usado no Server Component do layout para aplicar o tema sem JavaScript no cliente.
 */
export function buildThemeStyle(vars: ThemeVars): string {
  return Object.entries(vars)
    .map(([key, value]) => `${key}: ${value}`)
    .join('; ');
}

/**
 * Fallback de tema para quando a API não está disponível.
 * Usa a paleta da 5ª Temporada (verde-água / dourado).
 */
export const defaultTheme: ThemeVars = {
  '--color-primary': 'hsl(168 72% 40%)',
  '--color-secondary': 'hsl(45 95% 55%)',
  '--color-accent': 'hsl(168 80% 70%)',
  '--color-bg': 'hsl(220 20% 10%)',
  '--color-text': 'hsl(0 0% 97%)',
};
