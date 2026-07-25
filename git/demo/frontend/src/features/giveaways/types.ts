// features/giveaways/types.ts

export interface Giveaway {
  id: number;
  user_twitch_id: string;
  status: "open" | "closed";
  keyword: string;
  sub_bonus: boolean;
  winner_twitch_id?: string;
  winner_login?: string;
  started_at: string;
  ended_at?: string;
  created_at: string;
  updated_at: string;
  entry_count: number;
}

export interface GiveawayStatusResponse {
  giveaway: Giveaway | null;
  entry_count: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface StartGiveawayRequest {
  keyword: string;
  sub_bonus: boolean;
}
