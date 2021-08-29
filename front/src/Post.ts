export type Post = {
  id: number;
  uploaded_at?: string;
  user: string;
  comment?: string;
  delete_password?: string;
  url?: string;
  email: string;
  title: string;
  kind: string;
  count: number;
  limit_count?: number;
  require_message?: boolean;
  limit_date?: string;
};
