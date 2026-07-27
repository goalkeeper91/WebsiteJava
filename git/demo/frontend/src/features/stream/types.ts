export interface StreamInfo {
  broadcasterId: string;
  broadcasterName: string;
  title: string;
  gameId: string;
  gameName: string;
}

export interface UpdateStreamInfoRequest {
  title?: string;
  gameId?: string;
}

export interface LiveStream {
  isLive: boolean;
  viewerCount: number;
  startedAt?: string;
  title?: string;
  gameName?: string;
  thumbnailUrl?: string;
}

export interface Category {
  id: string;
  name: string;
  boxArtUrl?: string;
}

export interface DashboardStats {
  isLive: boolean;
  currentViewers: number;
  followerCount: number;
  subscriberCount: number;
  uptime?: string;

  // Statistiken
  followsToday: number;
  subsThisWeek: number;
  bitsToday: number;
  avgViewers: number;
}

export interface CommercialRequest {
  length: 30 | 60 | 90 | 120 | 150 | 180;
}

export interface CommercialResult {
  length: number;
  message: string;
  retryAfter: number;
}

// Muss exakt zu Go's domain.ActivityType passen (internal/domain/stream_activity.go).
export type ActivityType = "FOLLOW" | "SUBSCRIBE" | "RAID" | "CHEER" | "GIFT_SUB" | "HOSTING" | "RESUBSCRIBE";

// Feldnamen und Typen müssen exakt zu Go's domain.StreamActivity JSON passen
// (internal/domain/stream_activity.go) - snake_case + ISO-Timestamp-String,
// wie in jedem anderen Feature dieser App (z.B. Giveaway.started_at).
export interface Activity {
  id: number;
  type: ActivityType;
  username: string;
  display_name: string;
  timestamp: string;

  // Optional fields je nach Type
  viewers?: number; // for RAID
  bits?: number; // for CHEER
  tier?: string; // for SUBSCRIBE/GIFT_SUB
  message?: string; // for CHEER
}