// features/admin/customers/types.ts

export type TierId = "free" | "pro" | "premium";

export interface AdminCustomer {
  twitch_id: string;
  username: string;
  email?: string;
  is_admin: boolean;
  created_at: string;
  tier_id: string;
  tier_name: string;
  status: string;
  is_active: boolean;
  billing_cycle?: string;
  expires_at?: string;
  price_monthly: number;
  price_yearly: number;
  paddle_customer_id?: string;
  paddle_subscription_id?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ClipDetectorStatus {
  active_channels: { twitch_user_id: string; login: string }[];
  max_concurrent: number;
  updated_at: string;
}

export interface AdminCustomerStats {
  mrr: number;
  active_customers: number;
  clip_detector: ClipDetectorStatus | null;
}
