'use client';
// Step 1 — Coleta de dados pessoais e logística.
import { useState } from 'react';
import { ArrowRight, User, Phone, Shirt, AlertCircle } from 'lucide-react';
import { criarInscricao, getTemporadaAtiva } from '@/lib/api';
import type { CriarInscricaoResult, TamanhoVestuario } from '@/types';

const TAMANHOS: TamanhoVestuario[] = ['PP', 'P', 'M', 'G', 'GG', 'XGG'];
const PARENTESCOS = ['Mãe', 'Pai', 'Marido', 'Filho/a', 'Irmã/ão', 'Amigo/a', 'Outro'];

interface Step1Props {
  onComplete: (result: CriarInscricaoResult, cpf: string) => void;
}

export default function Step1Cadastro({ onComplete }: Step1Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [form, setForm] = useState({
    nome_completo: '',
    email: '',
    whatsapp: '',
    cpf: '',
    data_nascimento: '',
    emergencia_nome: '',
    emergencia_telefone: '',
    emergencia_parentesco: '',
    tamanho_camisa: '' as TamanhoVestuario | '',
    tamanho_balaclava: '' as TamanhoVestuario | '',
    tamanho_bone: '' as TamanhoVestuario | '',
    consentimento_lgpd: false,
  });

  function set(field: keyof typeof form) {
    return (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
      const value =
        e.target.type === 'checkbox'
          ? (e.target as HTMLInputElement).checked
          : e.target.value;
      setForm((prev) => ({ ...prev, [field]: value }));
      setError(null);
    };
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!form.consentimento_lgpd) {
      setError('Você precisa aceitar a política de privacidade para continuar.');
      return;
    }

    setLoading(true);
    try {
      // Buscar a temporada ativa para obter o temporada_id
      const temporada = await getTemporadaAtiva();

      const result = await criarInscricao({
        temporada_id: temporada.id,
        nome_completo: form.nome_completo.trim(),
        email: form.email.trim().toLowerCase(),
        whatsapp: form.whatsapp.replace(/\D/g, ''),
        cpf: form.cpf.replace(/\D/g, ''),
        data_nascimento: form.data_nascimento || undefined,
        emergencia_nome: form.emergencia_nome.trim(),
        emergencia_telefone: form.emergencia_telefone.replace(/\D/g, ''),
        emergencia_parentesco: form.emergencia_parentesco || undefined,
        tamanho_camisa: form.tamanho_camisa as TamanhoVestuario,
        tamanho_balaclava: form.tamanho_balaclava as TamanhoVestuario,
        tamanho_bone: form.tamanho_bone as TamanhoVestuario,
        consentimento_lgpd: true,
      });

      onComplete(result, form.cpf.replace(/\D/g, ''));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao processar inscrição.');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-8" noValidate>
      {/* ── Dados Pessoais ── */}
      <section className="glass p-6 md:p-8">
        <SectionTitle icon={<User size={18} />} title="Dados Pessoais" />

        <div className="grid gap-5">
          <Field label="Nome Completo *">
            <input
              className="input-field"
              type="text"
              placeholder="Seu nome completo"
              value={form.nome_completo}
              onChange={set('nome_completo')}
              required
              autoComplete="name"
            />
          </Field>

          <div className="grid sm:grid-cols-2 gap-5">
            <Field label="E-mail *">
              <input
                className="input-field"
                type="email"
                placeholder="seu@email.com"
                value={form.email}
                onChange={set('email')}
                required
                autoComplete="email"
              />
            </Field>

            <Field label="CPF *">
              <input
                className="input-field"
                type="text"
                placeholder="000.000.000-00"
                value={form.cpf}
                onChange={set('cpf')}
                required
                maxLength={14}
                inputMode="numeric"
              />
            </Field>
          </div>

          <div className="grid sm:grid-cols-2 gap-5">
            <Field label="WhatsApp *">
              <input
                className="input-field"
                type="tel"
                placeholder="(69) 9 0000-0000"
                value={form.whatsapp}
                onChange={set('whatsapp')}
                required
                autoComplete="tel"
              />
            </Field>

            <Field label="Data de Nascimento">
              <input
                className="input-field"
                type="date"
                value={form.data_nascimento}
                onChange={set('data_nascimento')}
              />
            </Field>
          </div>
        </div>
      </section>

      {/* ── Contato de Emergência ── */}
      <section className="glass p-6 md:p-8">
        <SectionTitle icon={<Phone size={18} />} title="Contato de Emergência" />

        <div className="grid gap-5">
          <div className="grid sm:grid-cols-2 gap-5">
            <Field label="Nome *">
              <input
                className="input-field"
                type="text"
                placeholder="Nome do contato"
                value={form.emergencia_nome}
                onChange={set('emergencia_nome')}
                required
              />
            </Field>

            <Field label="Telefone *">
              <input
                className="input-field"
                type="tel"
                placeholder="(69) 9 0000-0000"
                value={form.emergencia_telefone}
                onChange={set('emergencia_telefone')}
                required
              />
            </Field>
          </div>

          <Field label="Grau de Parentesco">
            <select
              className="input-field"
              value={form.emergencia_parentesco}
              onChange={set('emergencia_parentesco')}
            >
              <option value="">Selecionar...</option>
              {PARENTESCOS.map((p) => (
                <option key={p} value={p}>{p}</option>
              ))}
            </select>
          </Field>
        </div>
      </section>

      {/* ── Logística (Vestuário) ── */}
      <section className="glass p-6 md:p-8">
        <SectionTitle icon={<Shirt size={18} />} title="Tamanhos do Kit" />
        <p className="text-muted text-sm mb-5">
          Camisa, balaclava e boné personalizados serão entregues no evento nos tamanhos abaixo.
        </p>

        <div className="grid sm:grid-cols-3 gap-5">
          <SizeSelect
            label="Camisa *"
            value={form.tamanho_camisa}
            onChange={(v) => setForm((p) => ({ ...p, tamanho_camisa: v as TamanhoVestuario }))}
          />
          <SizeSelect
            label="Balaclava *"
            value={form.tamanho_balaclava}
            onChange={(v) => setForm((p) => ({ ...p, tamanho_balaclava: v as TamanhoVestuario }))}
          />
          <SizeSelect
            label="Boné *"
            value={form.tamanho_bone}
            onChange={(v) => setForm((p) => ({ ...p, tamanho_bone: v as TamanhoVestuario }))}
          />
        </div>
      </section>

      {/* ── LGPD ── */}
      <div className="glass-sm p-5">
        <label className="flex items-start gap-3 cursor-pointer group">
          <div className="relative mt-0.5">
            <input
              type="checkbox"
              className="sr-only"
              checked={form.consentimento_lgpd}
              onChange={set('consentimento_lgpd')}
              required
            />
            <div
              className="w-5 h-5 rounded border-2 flex items-center justify-center transition-all duration-200"
              style={{
                borderColor: form.consentimento_lgpd ? 'var(--color-primary)' : 'var(--color-border)',
                background: form.consentimento_lgpd ? 'var(--color-primary)' : 'transparent',
              }}
            >
              {form.consentimento_lgpd && (
                <svg width="10" height="8" viewBox="0 0 10 8" fill="none">
                  <path d="M1 4L3.5 6.5L9 1" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ color: 'var(--color-bg)' }} />
                </svg>
              )}
            </div>
          </div>
          <span className="text-sm text-muted leading-relaxed">
            Li e concordo com a{' '}
            <span className="text-primary font-semibold" style={{ color: 'var(--color-primary)' }}>
              Política de Privacidade
            </span>{' '}
            da afpro e autorizo o uso dos meus dados para processamento desta inscrição,
            conforme a Lei Geral de Proteção de Dados (LGPD).
          </span>
        </label>
      </div>

      {/* ── Erro ── */}
      {error && (
        <div
          className="glass-sm p-4 flex items-center gap-3 text-sm"
          style={{ borderColor: 'hsl(0 72% 55%)', border: '1px solid' }}
        >
          <AlertCircle size={18} style={{ color: 'hsl(0 72% 55%)', flexShrink: 0 }} />
          <span style={{ color: 'hsl(0 72% 55%)' }}>{error}</span>
        </div>
      )}

      {/* ── Submit ── */}
      <button
        type="submit"
        className="btn-primary w-full text-lg py-4"
        disabled={loading}
      >
        {loading ? (
          <span className="flex items-center gap-2">
            <LoadingSpinner />
            Processando...
          </span>
        ) : (
          <span className="flex items-center gap-2">
            Continuar para o Contrato
            <ArrowRight size={20} />
          </span>
        )}
      </button>
    </form>
  );
}

// ─── Sub-componentes ───────────────────────────────────────────────────────────

function SectionTitle({ icon, title }: { icon: React.ReactNode; title: string }) {
  return (
    <div className="flex items-center gap-2 mb-6 pb-4" style={{ borderBottom: '1px solid var(--color-border)' }}>
      <span style={{ color: 'var(--color-primary)' }}>{icon}</span>
      <h2 className="font-bold text-lg">{title}</h2>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="input-label">{label}</label>
      {children}
    </div>
  );
}

function SizeSelect({
  label,
  value,
  onChange,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <Field label={label}>
      <select
        className="input-field"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required
      >
        <option value="">Tamanho...</option>
        {TAMANHOS.map((t) => (
          <option key={t} value={t}>{t}</option>
        ))}
      </select>
    </Field>
  );
}

function LoadingSpinner() {
  return (
    <svg
      className="animate-spin"
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
    >
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" opacity="0.25" />
      <path
        d="M12 2a10 10 0 0 1 10 10"
        stroke="currentColor"
        strokeWidth="3"
        strokeLinecap="round"
      />
    </svg>
  );
}
