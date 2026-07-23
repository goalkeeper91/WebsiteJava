// features/automod/types.ts

export interface AutomodSettings {
  user_twitch_id: string;
  enabled: boolean;
  blocked_words: string[];
  link_filter_enabled: boolean;
  allowed_domains: string[];
  created_at: string;
  updated_at: string;
}

export interface UpdateAutomodSettingsRequest {
  enabled?: boolean;
  blocked_words?: string[];
  link_filter_enabled?: boolean;
  allowed_domains?: string[];
}
