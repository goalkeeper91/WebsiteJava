// features/loyalty/types.ts

export interface LoyaltySettings {
  user_twitch_id: string;
  enabled: boolean;
  points_name: string;
  points_per_interval: number;
  interval_minutes: number;
  regulars_threshold: number;
  next_accrual_at: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateLoyaltySettingsRequest {
  enabled?: boolean;
  points_name?: string;
  points_per_interval?: number;
  interval_minutes?: number;
  regulars_threshold?: number;
}

export interface LeaderboardEntry {
  user_twitch_id: string;
  viewer_twitch_id: string;
  viewer_login: string;
  points: number;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}
