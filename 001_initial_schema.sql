-- =============================================================================
-- AFPERO — "5ª Temporada Meninas na Pesca"
-- Schema SQL (DDL) — Supabase / PostgreSQL
-- Versão: 1.0.0
-- =============================================================================
-- INSTRUÇÕES DE USO:
-- 1. Abrir o Supabase Dashboard > SQL Editor
-- 2. Colar e executar este script completo
-- 3. Confirmar que todas as tabelas aparecem em "Table Editor"
-- =============================================================================


-- ─────────────────────────────────────────────────────────────────────────────
-- 0. EXTENSÕES
-- ─────────────────────────────────────────────────────────────────────────────

-- UUID v4 para PKs (já habilitado por padrão no Supabase, mas garantimos aqui)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- pgcrypto para hash de CPF (segurança de dados pessoais)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";


-- ─────────────────────────────────────────────────────────────────────────────
-- 1. TIPOS ENUM
-- ─────────────────────────────────────────────────────────────────────────────

-- Status do ciclo de vida da inscrição
CREATE TYPE status_inscricao_enum AS ENUM (
    'pendente',           -- Formulário submetido, aguardando assinatura do contrato
    'contrato_enviado',   -- Contrato enviado ao provider de assinatura
    'contrato_assinado',  -- Contrato assinado; PIX de entrada gerado
    'pago',               -- PIX de entrada confirmado; parcelas de boleto geradas
    'cancelado'           -- Inscrição cancelada (manualmente ou por expiração)
);

-- Status de cada transação financeira
CREATE TYPE status_transacao_enum AS ENUM (
    'pendente',    -- Cobrança gerada, aguardando pagamento
    'pago',        -- Pagamento confirmado via webhook
    'expirado',    -- Prazo de pagamento expirado
    'cancelado',   -- Cobrança cancelada
    'reembolsado'  -- Pagamento estornado
);

-- Tipos de cobrança
CREATE TYPE tipo_transacao_enum AS ENUM (
    'entrada_pix',    -- R$ 800,00 — entrada via PIX (gerada após assinatura)
    'parcela_boleto'  -- R$ 540,00 — parcelas mensais (5x, geradas após confirmação do PIX)
);

-- Providers de gateway de pagamento (adapter pattern)
CREATE TYPE provider_pagamento_enum AS ENUM (
    'asaas',
    'mercado_pago'
);

-- Providers de assinatura digital (adapter pattern)
CREATE TYPE provider_assinatura_enum AS ENUM (
    'zapsign',
    'clicksign'
);

-- Status do documento no ciclo de assinatura
CREATE TYPE status_documento_enum AS ENUM (
    'gerado',      -- PDF gerado e salvo no Storage
    'enviado',     -- Enviado ao provider de assinatura
    'visualizado', -- Participante abriu o link de assinatura
    'assinado',    -- Assinatura concluída e confirmada via webhook
    'cancelado'    -- Processo cancelado
);

-- Tamanhos de vestuário (camisa, balaclava, boné)
CREATE TYPE tamanho_vestuario_enum AS ENUM (
    'PP', 'P', 'M', 'G', 'GG', 'XGG'
);


-- ─────────────────────────────────────────────────────────────────────────────
-- 2. FUNÇÃO AUXILIAR: updated_at automático via trigger
-- ─────────────────────────────────────────────────────────────────────────────

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ─────────────────────────────────────────────────────────────────────────────
-- 3. TABELA: temporadas
-- Controla as edições do evento e os temas visuais do frontend.
-- A Engine de Temas Dinâmicos do Next.js consome a temporada com is_ativa = TRUE.
-- ─────────────────────────────────────────────────────────────────────────────

CREATE TABLE temporadas (
    -- Identificação
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome                TEXT NOT NULL,                  -- Ex: "5ª Temporada Meninas na Pesca"
    edicao              SMALLINT NOT NULL,               -- Ex: 5
    slug                TEXT NOT NULL UNIQUE,            -- Ex: "5a-temporada-2025"

    -- Período do evento
    data_inicio_evento  DATE NOT NULL,
    data_fim_evento     DATE NOT NULL,

    -- Período de inscrições
    data_abertura_inscricoes  DATE NOT NULL,
    data_fechamento_inscricoes DATE NOT NULL,

    -- Vagas
    vagas_total         SMALLINT NOT NULL DEFAULT 30,
    vagas_disponiveis   SMALLINT NOT NULL DEFAULT 30,

    -- Financeiro da temporada
    valor_total         NUMERIC(10, 2) NOT NULL DEFAULT 3500.00,
    valor_entrada_pix   NUMERIC(10, 2) NOT NULL DEFAULT 800.00,
    valor_parcela       NUMERIC(10, 2) NOT NULL DEFAULT 540.00,
    qtd_parcelas        SMALLINT NOT NULL DEFAULT 5,

    -- Configurações visuais (consumidas pelo frontend em tempo de execução)
    -- Todos os valores devem ser CSS válidos: hex (#FF5733), hsl(20 80% 50%), rgb(...)
    cor_primaria        TEXT NOT NULL DEFAULT '#1a1a2e',
    cor_secundaria      TEXT NOT NULL DEFAULT '#e94560',
    cor_acento          TEXT NOT NULL DEFAULT '#f5a623',
    cor_fundo           TEXT NOT NULL DEFAULT '#0f0f1a',
    cor_texto           TEXT NOT NULL DEFAULT '#ffffff',
    logo_url            TEXT,
    banner_url          TEXT,

    -- Template do contrato (placeholders {{nome}}, {{cpf}}, etc.)
    contrato_template   TEXT,

    -- Controle
    is_ativa            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_datas_evento CHECK (data_fim_evento >= data_inicio_evento),
    CONSTRAINT chk_datas_inscricao CHECK (data_fechamento_inscricoes >= data_abertura_inscricoes),
    CONSTRAINT chk_vagas CHECK (vagas_disponiveis >= 0 AND vagas_disponiveis <= vagas_total),
    CONSTRAINT chk_valor_total CHECK (
        ABS((valor_entrada_pix + (valor_parcela * qtd_parcelas)) - valor_total) < 0.01
    )
);

CREATE TRIGGER trg_temporadas_updated_at
    BEFORE UPDATE ON temporadas
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_temporadas_is_ativa ON temporadas (is_ativa) WHERE is_ativa = TRUE;
CREATE INDEX idx_temporadas_slug ON temporadas (slug);

-- Apenas UMA temporada ativa por vez
CREATE UNIQUE INDEX idx_temporadas_unica_ativa ON temporadas (is_ativa)
    WHERE is_ativa = TRUE;

COMMENT ON TABLE temporadas IS
    'Edições do evento. Controla o tema visual dinâmico do frontend e as configurações financeiras.';
COMMENT ON COLUMN temporadas.is_ativa IS
    'Apenas UMA temporada pode ser ativa por vez (garantido por índice único parcial).';
COMMENT ON COLUMN temporadas.cor_primaria IS
    'Valor CSS válido: #e94560 ou hsl(350 78% 60%).';
COMMENT ON COLUMN temporadas.contrato_template IS
    'Template com placeholders: {{PARTICIPANTE_NOME}}, {{CPF_LAST4}}, {{DATA_ASSINATURA}}, etc.';


-- ─────────────────────────────────────────────────────────────────────────────
-- 4. TABELA: participantes
-- ─────────────────────────────────────────────────────────────────────────────

CREATE TABLE participantes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    temporada_id    UUID NOT NULL REFERENCES temporadas(id) ON DELETE RESTRICT,

    -- Dados pessoais
    nome_completo   TEXT NOT NULL,
    email           TEXT NOT NULL,
    whatsapp        TEXT NOT NULL,
    cpf_hash        TEXT NOT NULL,    -- SHA-256 do CPF puro (apenas dígitos)
    cpf_last4       CHAR(4) NOT NULL, -- Últimos 4 dígitos para exibição
    data_nascimento DATE,

    -- Contato de emergência
    emergencia_nome         TEXT NOT NULL,
    emergencia_telefone     TEXT NOT NULL,
    emergencia_parentesco   TEXT,

    -- Logística (vestuário)
    tamanho_camisa      tamanho_vestuario_enum NOT NULL,
    tamanho_balaclava   tamanho_vestuario_enum NOT NULL,
    tamanho_bone        tamanho_vestuario_enum NOT NULL,

    -- Status do ciclo de inscrição
    status  status_inscricao_enum NOT NULL DEFAULT 'pendente',

    -- Metadados do formulário
    ip_origem           INET,
    user_agent          TEXT,
    consentimento_lgpd  BOOLEAN NOT NULL DEFAULT FALSE,
    consentimento_at    TIMESTAMPTZ,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_whatsapp_formato CHECK (whatsapp ~ '^\d{10,11}$'),
    CONSTRAINT chk_cpf_last4 CHECK (cpf_last4 ~ '^\d{4}$'),
    CONSTRAINT chk_consentimento_lgpd CHECK (
        (consentimento_lgpd = FALSE) OR (consentimento_at IS NOT NULL)
    )
);

CREATE TRIGGER trg_participantes_updated_at
    BEFORE UPDATE ON participantes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_participantes_temporada  ON participantes (temporada_id);
CREATE INDEX idx_participantes_status     ON participantes (status);
CREATE INDEX idx_participantes_cpf_hash   ON participantes (cpf_hash);
CREATE INDEX idx_participantes_email      ON participantes (email);
CREATE INDEX idx_participantes_whatsapp   ON participantes (whatsapp);

-- Um CPF não pode se inscrever duas vezes na mesma temporada
CREATE UNIQUE INDEX idx_participantes_cpf_por_temporada
    ON participantes (cpf_hash, temporada_id);

COMMENT ON TABLE participantes IS
    'Inscritas no evento. CPF armazenado apenas como hash SHA-256 (conformidade LGPD).';
COMMENT ON COLUMN participantes.cpf_hash IS
    'SHA-256 do CPF puro (apenas dígitos, sem máscara). Gerado pelo backend Go.';
COMMENT ON COLUMN participantes.status IS
    'Ciclo: pendente → contrato_enviado → contrato_assinado → pago (ou cancelado).';


-- ─────────────────────────────────────────────────────────────────────────────
-- 5. TABELA: documentos_assinados
-- ─────────────────────────────────────────────────────────────────────────────

CREATE TABLE documentos_assinados (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    participante_id UUID NOT NULL REFERENCES participantes(id) ON DELETE RESTRICT,
    temporada_id    UUID NOT NULL REFERENCES temporadas(id) ON DELETE RESTRICT,

    -- Provider de assinatura (adapter pattern)
    provider    provider_assinatura_enum NOT NULL,

    -- Identificadores externos
    provider_doc_id     TEXT,   -- ID do documento no provider
    provider_signer_id  TEXT,   -- ID do signatário no provider

    -- URLs
    documento_url           TEXT, -- PDF base no Supabase Storage
    link_assinatura         TEXT, -- Link enviado à participante
    documento_assinado_url  TEXT, -- PDF final assinado

    -- Status
    status  status_documento_enum NOT NULL DEFAULT 'gerado',

    -- Conteúdo (snapshot do contrato com dados reais no momento da geração)
    conteudo_html   TEXT,

    -- Timestamps do ciclo
    enviado_em      TIMESTAMPTZ,
    visualizado_em  TIMESTAMPTZ,
    assinado_em     TIMESTAMPTZ,

    -- Auditoria do webhook
    webhook_payload     JSONB,
    webhook_recebido_em TIMESTAMPTZ,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_documentos_updated_at
    BEFORE UPDATE ON documentos_assinados
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_documentos_participante   ON documentos_assinados (participante_id);
CREATE INDEX idx_documentos_status         ON documentos_assinados (status);
CREATE INDEX idx_documentos_provider_id    ON documentos_assinados (provider_doc_id);
CREATE INDEX idx_documentos_temporada      ON documentos_assinados (temporada_id);

-- Uma participante tem no máximo um documento ativo por temporada
CREATE UNIQUE INDEX idx_documentos_unico_por_participante
    ON documentos_assinados (participante_id, temporada_id)
    WHERE status NOT IN ('cancelado');

COMMENT ON TABLE documentos_assinados IS
    'Ciclo de vida do contrato/termo de adesão. Suporta múltiplos providers via campo provider.';
COMMENT ON COLUMN documentos_assinados.webhook_payload IS
    'Payload JSON bruto do webhook do provider, preservado para auditoria e reprocessamento.';


-- ─────────────────────────────────────────────────────────────────────────────
-- 6. TABELA: transacoes
-- 1 PIX entrada (R$ 800) + 5 boletos (R$ 540 cada) = R$ 3.500 total por inscrita.
-- ─────────────────────────────────────────────────────────────────────────────

CREATE TABLE transacoes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    participante_id UUID NOT NULL REFERENCES participantes(id) ON DELETE RESTRICT,
    temporada_id    UUID NOT NULL REFERENCES temporadas(id) ON DELETE RESTRICT,

    -- Tipo e sequência
    tipo            tipo_transacao_enum NOT NULL,
    numero_parcela  SMALLINT,       -- NULL para entrada PIX; 1–5 para boletos
    descricao       TEXT NOT NULL,  -- Ex: "Entrada PIX — 5ª Temporada" | "Parcela 2/5"

    -- Valores
    valor           NUMERIC(10, 2) NOT NULL,
    data_vencimento DATE NOT NULL,

    -- Provider de pagamento (adapter pattern)
    provider            provider_pagamento_enum NOT NULL,
    provider_cobranca_id TEXT,   -- ID da cobrança no provider
    provider_customer_id TEXT,   -- ID do cliente no provider

    -- Campos PIX (apenas para tipo = 'entrada_pix')
    pix_qr_code_base64  TEXT,           -- QR Code em Base64
    pix_copia_cola      TEXT,           -- Código "Copia e Cola"
    pix_expiracao       TIMESTAMPTZ,    -- NOW() + 30 minutos

    -- Campos Boleto (apenas para tipo = 'parcela_boleto')
    boleto_url          TEXT,   -- URL do boleto
    boleto_codigo_barras TEXT,  -- Linha digitável

    -- Status
    status  status_transacao_enum NOT NULL DEFAULT 'pendente',
    pago_em TIMESTAMPTZ,  -- Preenchido via webhook quando status = 'pago'

    -- Auditoria do webhook
    webhook_payload     JSONB,
    webhook_recebido_em TIMESTAMPTZ,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_valor_positivo CHECK (valor > 0),
    CONSTRAINT chk_numero_parcela CHECK (
        (tipo = 'entrada_pix' AND numero_parcela IS NULL) OR
        (tipo = 'parcela_boleto' AND numero_parcela >= 1)
    ),
    CONSTRAINT chk_pix_dados CHECK (
        tipo = 'entrada_pix' OR (pix_qr_code_base64 IS NULL AND pix_copia_cola IS NULL)
    )
);

CREATE TRIGGER trg_transacoes_updated_at
    BEFORE UPDATE ON transacoes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_transacoes_participante ON transacoes (participante_id);
CREATE INDEX idx_transacoes_status       ON transacoes (status);
CREATE INDEX idx_transacoes_tipo         ON transacoes (tipo);
CREATE INDEX idx_transacoes_temporada    ON transacoes (temporada_id);
CREATE INDEX idx_transacoes_vencimento   ON transacoes (data_vencimento);
CREATE INDEX idx_transacoes_provider_id  ON transacoes (provider_cobranca_id);

-- Entrada PIX única por participante/temporada
CREATE UNIQUE INDEX idx_transacoes_entrada_unica
    ON transacoes (participante_id, temporada_id)
    WHERE tipo = 'entrada_pix' AND status != 'cancelado';

-- Unicidade de parcela por participante/temporada
CREATE UNIQUE INDEX idx_transacoes_parcela_unica
    ON transacoes (participante_id, temporada_id, numero_parcela)
    WHERE tipo = 'parcela_boleto';

COMMENT ON TABLE transacoes IS
    'Cobranças financeiras individuais. 1 PIX entrada (R$ 800) + 5 boletos (R$ 540 cada) = R$ 3.500/participante.';
COMMENT ON COLUMN transacoes.pix_expiracao IS
    'Expiração do QR Code PIX. O backend Go define como NOW() + 30 minutos.';
COMMENT ON COLUMN transacoes.webhook_payload IS
    'Payload JSON bruto do gateway, preservado para auditoria e reconciliação.';


-- ─────────────────────────────────────────────────────────────────────────────
-- 7. ROW LEVEL SECURITY (RLS)
-- ─────────────────────────────────────────────────────────────────────────────

ALTER TABLE temporadas            ENABLE ROW LEVEL SECURITY;
ALTER TABLE participantes         ENABLE ROW LEVEL SECURITY;
ALTER TABLE documentos_assinados  ENABLE ROW LEVEL SECURITY;
ALTER TABLE transacoes            ENABLE ROW LEVEL SECURITY;

-- ── temporadas: leitura pública (frontend busca o tema sem autenticação) ──
CREATE POLICY "temporadas_leitura_publica"
    ON temporadas FOR SELECT
    USING (TRUE);

CREATE POLICY "temporadas_escrita_service_role"
    ON temporadas FOR ALL
    USING (auth.role() = 'service_role')
    WITH CHECK (auth.role() = 'service_role');

-- ── participantes ──
CREATE POLICY "participantes_insercao_publica"
    ON participantes FOR INSERT
    WITH CHECK (TRUE);  -- Formulário público

CREATE POLICY "participantes_leitura_service_role"
    ON participantes FOR SELECT
    USING (auth.role() = 'service_role');

CREATE POLICY "participantes_atualizacao_service_role"
    ON participantes FOR UPDATE
    USING (auth.role() = 'service_role')
    WITH CHECK (auth.role() = 'service_role');

-- ── documentos_assinados: apenas service_role ──
CREATE POLICY "documentos_service_role"
    ON documentos_assinados FOR ALL
    USING (auth.role() = 'service_role')
    WITH CHECK (auth.role() = 'service_role');

-- ── transacoes: apenas service_role ──
CREATE POLICY "transacoes_service_role"
    ON transacoes FOR ALL
    USING (auth.role() = 'service_role')
    WITH CHECK (auth.role() = 'service_role');


-- ─────────────────────────────────────────────────────────────────────────────
-- 8. VIEWS UTILITÁRIAS
-- ─────────────────────────────────────────────────────────────────────────────

-- View pública: tema da temporada ativa (consumida pelo Next.js)
CREATE OR REPLACE VIEW vw_temporada_ativa AS
SELECT
    id,
    nome,
    edicao,
    slug,
    data_inicio_evento,
    data_fim_evento,
    data_abertura_inscricoes,
    data_fechamento_inscricoes,
    vagas_total,
    vagas_disponiveis,
    valor_total,
    valor_entrada_pix,
    valor_parcela,
    qtd_parcelas,
    cor_primaria,
    cor_secundaria,
    cor_acento,
    cor_fundo,
    cor_texto,
    logo_url,
    banner_url
FROM temporadas
WHERE is_ativa = TRUE
LIMIT 1;

COMMENT ON VIEW vw_temporada_ativa IS
    'Dados públicos da temporada ativa. Consumida pelo Next.js para injetar o tema visual em runtime.';


-- View administrativa: painel de inscrições (service_role)
CREATE OR REPLACE VIEW vw_painel_inscricoes AS
SELECT
    p.id                AS participante_id,
    p.nome_completo,
    p.email,
    p.whatsapp,
    p.cpf_last4,
    p.status            AS status_inscricao,
    p.tamanho_camisa,
    p.tamanho_balaclava,
    p.tamanho_bone,
    t.nome              AS temporada,
    d.status            AS status_documento,
    d.provider          AS provider_assinatura,
    d.assinado_em,
    COALESCE(tx_entrada.status::TEXT, 'nao_gerado') AS status_pix_entrada,
    tx_entrada.valor        AS valor_entrada,
    tx_entrada.pago_em      AS entrada_paga_em,
    COUNT(tx_parcelas.id)   AS parcelas_geradas,
    SUM(CASE WHEN tx_parcelas.status = 'pago' THEN tx_parcelas.valor ELSE 0 END) AS total_pago_parcelas,
    p.created_at        AS inscrito_em
FROM participantes p
JOIN temporadas t ON t.id = p.temporada_id
LEFT JOIN documentos_assinados d
    ON d.participante_id = p.id AND d.status != 'cancelado'
LEFT JOIN transacoes tx_entrada
    ON tx_entrada.participante_id = p.id
    AND tx_entrada.tipo = 'entrada_pix'
    AND tx_entrada.status != 'cancelado'
LEFT JOIN transacoes tx_parcelas
    ON tx_parcelas.participante_id = p.id
    AND tx_parcelas.tipo = 'parcela_boleto'
GROUP BY
    p.id, p.nome_completo, p.email, p.whatsapp, p.cpf_last4, p.status,
    p.tamanho_camisa, p.tamanho_balaclava, p.tamanho_bone,
    t.nome, d.status, d.provider, d.assinado_em,
    tx_entrada.status, tx_entrada.valor, tx_entrada.pago_em, p.created_at;

COMMENT ON VIEW vw_painel_inscricoes IS
    'Visão consolidada para o painel administrativo. Acesso restrito ao service_role.';


-- ─────────────────────────────────────────────────────────────────────────────
-- 9. DADOS DE SEMENTE — 5ª Temporada
-- ─────────────────────────────────────────────────────────────────────────────

INSERT INTO temporadas (
    nome,
    edicao,
    slug,
    data_inicio_evento,
    data_fim_evento,
    data_abertura_inscricoes,
    data_fechamento_inscricoes,
    vagas_total,
    vagas_disponiveis,
    valor_total,
    valor_entrada_pix,
    valor_parcela,
    qtd_parcelas,
    cor_primaria,
    cor_secundaria,
    cor_acento,
    cor_fundo,
    cor_texto,
    logo_url,
    banner_url,
    contrato_template,
    is_ativa
) VALUES (
    '5ª Temporada Meninas na Pesca',
    5,
    '5a-temporada-2025',
    '2025-09-01',
    '2025-09-07',
    '2025-07-01',
    '2025-08-20',
    30,
    30,
    3500.00,
    800.00,
    540.00,
    5,
    -- Paleta: verde-água profundo / dourado (identidade da 5ª Temporada)
    'hsl(168 72% 40%)',
    'hsl(45 95% 55%)',
    'hsl(168 80% 70%)',
    'hsl(220 20% 10%)',
    'hsl(0 0% 97%)',
    NULL,  -- logo_url: definir após upload no Supabase Storage
    NULL,  -- banner_url: definir após upload
    -- Template do contrato com placeholders para o backend Go substituir
    E'# TERMO DE ADESÃO — {{TEMPORADA_NOME}}\n\n'
    E'**Participante:** {{PARTICIPANTE_NOME}}\n'
    E'**CPF:** ***.***.**{{CPF_LAST4}}\n'
    E'**E-mail:** {{PARTICIPANTE_EMAIL}}\n'
    E'**WhatsApp:** {{PARTICIPANTE_WHATSAPP}}\n\n'
    E'## 1. Objeto\n'
    E'O presente termo regula a participação de {{PARTICIPANTE_NOME}} no evento '
    E'**{{TEMPORADA_NOME}}**, promovido pela Associação Feminina de Pesca Esportiva '
    E'de Rondônia (AFPERO), a realizar-se no período de {{DATA_INICIO}} a {{DATA_FIM}}.\n\n'
    E'## 2. Investimento e Forma de Pagamento\n'
    E'O valor total de participação é de **R$ {{VALOR_TOTAL}}**, dividido da seguinte forma:\n'
    E'- **Entrada via PIX:** R$ {{VALOR_ENTRADA}} (vencimento: {{DATA_VENCIMENTO_ENTRADA}})\n'
    E'- **{{QTD_PARCELAS}} parcelas mensais via boleto:** R$ {{VALOR_PARCELA}} cada\n\n'
    E'O não pagamento da entrada até o vencimento implicará no cancelamento automático desta inscrição.\n\n'
    E'## 3. Obrigações da Participante\n'
    E'A participante se compromete a: (i) efetuar os pagamentos nos prazos estabelecidos; '
    E'(ii) seguir as normas de segurança e conduta do evento; '
    E'(iii) informar condições de saúde relevantes à organização.\n\n'
    E'## 4. Cancelamento\n'
    E'Em caso de desistência, a AFPERO reterá 30% do valor já pago a título de taxa administrativa. '
    E'Cancelamentos com menos de 30 dias do evento não são reembolsáveis.\n\n'
    E'## 5. Uso de Imagem\n'
    E'A participante autoriza o uso de imagens e vídeos captados durante o evento para fins de '
    E'divulgação nas redes sociais e materiais institucionais da AFPERO, sem ônus financeiro.\n\n'
    E'## 6. Responsabilidade\n'
    E'A AFPERO não se responsabiliza por eventuais acidentes decorrentes de negligência '
    E'da participante ou de situações de força maior.\n\n'
    E'**Data de assinatura:** {{DATA_ASSINATURA}}',
    TRUE
);
