// features/scheduledMessages/types.ts

export interface ScheduledMessage {
  id: number;
  channel_id: string;
  message: string | null;
  command_id?: number;
  interval_seconds: number;
  enabled: boolean;
  only_when_live: boolean;
  next_send_at: string;
  last_sent_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateScheduledMessageRequest {
  message?: string;
  command_id?: number;
  interval_seconds: number;
}

export interface UpdateScheduledMessageRequest {
  message?: string;
  interval_seconds?: number;
  enabled?: boolean;
  only_when_live?: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}
