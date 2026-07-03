// Tipos TypeScript espelhando o schema do banco de dados
export interface Temporada {
  id: string;
  nome: string;
  edicao: number;
  slug: string;
  data_inicio_evento: string;
  data_fim_evento: string;
  data_abertura_inscricoes: string;
  data_fechamento_inscricoes: string;
  vagas_total: number;
  vagas_disponiveis: number;
  valor_total: number;
  valor_entrada_pix: number;
  valor_parcela: number;
  qtd_parcelas: number;
  // Tema visual
  cor_primaria: string;
  cor_secundaria: string;
  cor_acento: string;
  cor_fundo: string;
  cor_texto: string;
  logo_url: string | null;
  banner_url: string | null;
}

export interface CriarInscricaoPayload {
  temporada_id: string;
  nome_completo: string;
  email: string;
  whatsapp: string;
  cpf: string;
  data_nascimento?: string;
  emergencia_nome: string;
  emergencia_telefone: string;
  emergencia_parentesco?: string;
  tamanho_camisa: string;
  tamanho_balaclava: string;
  tamanho_bone: string;
  consentimento_lgpd: boolean;
}

export interface CriarInscricaoResult {
  participante_id: string;
  contrato_html: string;
}

export interface AceitarContratoResult {
  transacao_id: string;
  pix_qr_code_base64: string;
  pix_copia_cola: string;
  pix_expiracao: string;
  valor_entrada: number;
}

export type TamanhoVestuario = 'PP' | 'P' | 'M' | 'G' | 'GG' | 'XGG';

export type StepInscricao = 1 | 2 | 3;
