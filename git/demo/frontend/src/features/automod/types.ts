// features/automod/types.ts

export interface AutomodSettings {
  user_twitch_id: string;
  enabled: boolean;
  blocked_words: string[];
  link_filter_enabled: boolean;
  allowed_domains: string[];
  caps_filter_enabled: boolean;
  symbol_filter_enabled: boolean;
  emote_filter_enabled: boolean;
  emote_threshold: number;
  repetition_filter_enabled: boolean;
  exempt_vips: boolean;
  exempt_regulars: boolean;
  exempt_users: string[];
  created_at: string;
  updated_at: string;
}

export interface UpdateAutomodSettingsRequest {
  enabled?: boolean;
  blocked_words?: string[];
  link_filter_enabled?: boolean;
  allowed_domains?: string[];
  caps_filter_enabled?: boolean;
  symbol_filter_enabled?: boolean;
  emote_filter_enabled?: boolean;
  emote_threshold?: number;
  repetition_filter_enabled?: boolean;
  exempt_vips?: boolean;
  exempt_regulars?: boolean;
  exempt_users?: string[];
}

export interface AutomodEvent {
  id: number;
  user_twitch_id: string;
  offender_twitch_id: string;
  offender_name: string;
  reason: string;
  message_excerpt?: string;
  timeout_seconds: number;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}
