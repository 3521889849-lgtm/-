export interface MerchantApplication {
  id?: string;
  merchant_name: string;
  contact_person: string;
  contact_phone: string;
  status?: 'pending' | 'approved' | 'rejected' | 'draft';
  created_at?: string;
}

export interface QualificationFile {
  id?: string;
  file_name: string;
  file_url: string;
  file_type?: string;
}

export interface AuditConfig {
  id: string;
  name: string;
  steps: AuditStep[];
}

export interface AuditStep {
  name: string;
  role: string;
}

export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data: T;
}
