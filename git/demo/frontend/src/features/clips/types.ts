// features/clips/types.ts

export type AITone = "neutral" | "casual" | "professional" | "energetic" | "funny";
export type AIStyle = "standard" | "short" | "detailed" | "viral" | "storytelling";
export type PostingFrequency = "realtime" | "hourly" | "daily" | "weekly" | "manual";

// v1: nur Link-Sharing-Ziele (discord/twitter) sind auswählbar.
// tiktok/youtube/instagram sind für natives Re-Upload vorgesehen (spätere Phase).
export type Platform = "tiktok" | "youtube" | "instagram" | "discord" | "twitter";

export const LINK_SHARE_PLATFORMS: Platform[] = ["discord", "twitter"];
export const NATIVE_UPLOAD_PLATFORMS: Platform[] = ["tiktok", "youtube", "instagram"];

export interface AutomationSettings {
  user_twitch_id: string;
  is_enabled: boolean;
  ai_tone: AITone;
  ai_style: AIStyle;
  use_hashtags: boolean;
  max_hashtags: number;
  min_clip_duration: number;
  max_clip_duration: number;
  min_view_count: number;
  only_verified_clips: boolean;
  target_platforms: Platform[];
  auto_post_enabled: boolean;
  posting_frequency: PostingFrequency;
  preferred_posting_times?: string[];
  blocked_words?: string[];
  allowed_games?: string[];
  notify_on_post: boolean;
  notify_on_error: boolean;
  created_at: string;
  updated_at: string;
}

export type UpdateAutomationSettingsRequest = Partial<
  Pick<
    AutomationSettings,
    | "is_enabled"
    | "ai_tone"
    | "ai_style"
    | "use_hashtags"
    | "max_hashtags"
    | "min_clip_duration"
    | "max_clip_duration"
    | "min_view_count"
    | "only_verified_clips"
    | "target_platforms"
    | "auto_post_enabled"
    | "posting_frequency"
    | "blocked_words"
    | "allowed_games"
    | "notify_on_post"
    | "notify_on_error"
  >
>;

export type ClipStatus = "pending" | "processing" | "completed" | "failed" | "cancelled";

export interface PlatformPost {
  platform: Platform;
  post_id: string;
  url: string;
  posted_at: string;
}

export interface ClipLog {
  id: string;
  user_twitch_id: string;
  clip_id: string;
  clip_url: string;
  clip_title?: string;
  broadcaster_name?: string;
  game_name?: string;
  duration_seconds: number;
  view_count: number;
  status: ClipStatus;
  ai_caption?: string;
  ai_hashtags?: string[];
  posted_platforms: PlatformPost[];
  error_message?: string;
  retry_count: number;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface ClipLogStats {
  total_clips: number;
  pending_clips: number;
  processing_clips: number;
  completed_clips: number;
  failed_clips: number;
  success_rate: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}
