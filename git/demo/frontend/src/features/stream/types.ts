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
  avgViewers: number;
}

export type ActivityType = "FOLLOW" | "SUB" | "RAID" | "CHEER" | "HOST";

export interface Activity {
  type: ActivityType;
  username: string;
  displayName: string;
  timestamp: number;

  // Optional fields je nach Type
  viewers?: number;     // for RAID
  bits?: number;        // for CHEER
  tier?: string;        // for SUB
  message?: string;     // for CHEER
}