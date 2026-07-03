'use client';
// InfoSection — seção de informações sobre o evento, pagamento e logística.
import { CreditCard, FileText, Shirt, Shield, Star, Clock } from 'lucide-react';
import type { Temporada } from '@/types';

interface InfoSectionProps {
  temporada: Temporada | null;
}

export default function InfoSection({ temporada }: InfoSectionProps) {
  const valorEntrada = temporada?.valor_entrada_pix ?? 800;
  const valorParcela = temporada?.valor_parcela ?? 540;
  const qtdParcelas  = temporada?.qtd_parcelas ?? 5;
  const valorTotal   = temporada?.valor_total ?? 3500;

  return (
    <section id="informacoes" className="relative z-10 py-24 px-6">
      <div className="max-w-6xl mx-auto">

        {/* Título da seção */}
        <div className="text-center mb-16">
          <div className="divider-gradient mb-8" />
          <h2 className="text-3xl md:text-5xl font-black mb-4">
            <span className="gradient-text">Tudo que você</span>
            <br />
            <span>precisa saber</span>
          </h2>
          <p className="text-muted max-w-xl mx-auto">
            Uma experiência completa, organizada com cuidado e atenção para cada participante.
          </p>
        </div>

        {/* Cards de features */}
        <div className="grid md:grid-cols-3 gap-6 mb-16">
          {features.map((feat, i) => (
            <div
              key={i}
              className="glass p-6 group hover:scale-[1.02] transition-transform duration-300"
              style={{ animationDelay: `${i * 0.1}s` }}
            >
              <div
                className="w-12 h-12 rounded-xl flex items-center justify-center mb-4 transition-colors duration-300"
                style={{
                  background: 'var(--color-primary-dim)',
                  color: 'var(--color-primary)',
                }}
              >
                {feat.icon}
              </div>
              <h3 className="font-bold text-lg mb-2">{feat.title}</h3>
              <p className="text-muted text-sm leading-relaxed">{feat.desc}</p>
            </div>
          ))}
        </div>

        {/* Plano de pagamento */}
        <div className="glass p-8 md:p-12">
          <div className="flex items-center gap-3 mb-8">
            <CreditCard style={{ color: 'var(--color-secondary)' }} size={24} />
            <h2 className="text-2xl font-black">Plano de Pagamento</h2>
          </div>

          <div className="grid md:grid-cols-2 gap-8 items-center">
            {/* Visualização das parcelas */}
            <div className="space-y-3">
              {/* Entrada PIX */}
              <PaymentRow
                label="Entrada via PIX"
                value={formatBRL(valorEntrada)}
                highlight
                icon={<Clock size={14} />}
                detail="Gerada ao aceitar o contrato — expira em 30 min"
              />

              {/* Parcelas */}
              {Array.from({ length: qtdParcelas }, (_, i) => (
                <PaymentRow
                  key={i}
                  label={`Parcela ${i + 1}/${qtdParcelas}`}
                  value={formatBRL(valorParcela)}
                  detail={`Boleto mensal`}
                />
              ))}

              {/* Total */}
              <div className="divider-gradient my-4" />
              <div className="flex justify-between items-center font-black text-lg">
                <span>Total</span>
                <span className="gradient-text">{formatBRL(valorTotal)}</span>
              </div>
            </div>

            {/* Texto explicativo */}
            <div className="space-y-4">
              <div className="glass-sm p-5">
                <div className="flex gap-3 items-start">
                  <Shield className="shrink-0 mt-0.5" style={{ color: 'var(--color-primary)' }} size={18} />
                  <div>
                    <p className="font-semibold text-sm mb-1">Contrato digital seguro</p>
                    <p className="text-muted text-sm">
                      Após o preenchimento dos dados, você receberá um link para assinar o
                      termo de adesão digitalmente. O PIX de entrada só é gerado após a assinatura.
                    </p>
                  </div>
                </div>
              </div>

              <div className="glass-sm p-5">
                <div className="flex gap-3 items-start">
                  <FileText className="shrink-0 mt-0.5" style={{ color: 'var(--color-secondary)' }} size={18} />
                  <div>
                    <p className="font-semibold text-sm mb-1">Boletos automáticos</p>
                    <p className="text-muted text-sm">
                      Após a confirmação do PIX de entrada, os {qtdParcelas} boletos mensais
                      são gerados automaticamente e enviados para o seu e-mail.
                    </p>
                  </div>
                </div>
              </div>

              <div className="glass-sm p-5">
                <div className="flex gap-3 items-start">
                  <Star className="shrink-0 mt-0.5" style={{ color: 'var(--color-accent)' }} size={18} />
                  <div>
                    <p className="font-semibold text-sm mb-1">Incluso na inscrição</p>
                    <p className="text-muted text-sm">
                      Camisa, balaclava e boné personalizados nos tamanhos escolhidos.
                      Estrutura completa e monitores especializados durante o evento.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ─── Componentes internos ──────────────────────────────────────────────────────

function PaymentRow({
  label,
  value,
  detail,
  highlight,
  icon,
}: {
  label: string;
  value: string;
  detail?: string;
  highlight?: boolean;
  icon?: React.ReactNode;
}) {
  return (
    <div
      className={`flex items-center justify-between p-3.5 rounded-xl transition-all ${
        highlight
          ? 'border-2'
          : 'bg-surface'
      }`}
      style={
        highlight
          ? {
              borderColor: 'var(--color-primary)',
              background: 'var(--color-primary-dim)',
            }
          : undefined
      }
    >
      <div>
        <div className="flex items-center gap-1.5 font-semibold text-sm">
          {icon && <span style={{ color: 'var(--color-primary)' }}>{icon}</span>}
          {label}
          {highlight && (
            <span className="badge badge-primary ml-1">PIX</span>
          )}
        </div>
        {detail && <div className="text-xs text-muted mt-0.5">{detail}</div>}
      </div>
      <span
        className={`font-bold ${highlight ? 'gradient-text text-lg' : ''}`}
      >
        {value}
      </span>
    </div>
  );
}

const features = [
  {
    icon: <Shirt size={22} />,
    title: 'Kit completo',
    desc: 'Camisa, balaclava e boné com identidade visual exclusiva da temporada, nos tamanhos que você escolher.',
  },
  {
    icon: <Shield size={22} />,
    title: 'Segurança garantida',
    desc: 'Monitores especializados, equipamentos de segurança e protocolos rigorosos para toda a equipe.',
  },
  {
    icon: <Star size={22} />,
    title: 'Experiência única',
    desc: 'Imersão completa na pesca esportiva com mulheres incríveis. Natureza, amizade e conquista.',
  },
];

function formatBRL(value: number) {
  return value.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
}
