// features/commands/types.ts

export interface ChatCommand {
  id: number;
  channel_id: string;
  trigger: string;
  response: string;
  cooldown: number;
  enabled: boolean;
  command_type: "simple" | "advanced";
  created_at: string;
  updated_at: string;
}

export interface CreateCommandRequest {
  trigger: string;
  response: string;
  cooldown: number;
}

export interface UpdateCommandRequest {
  trigger?: string;
  response?: string;
  cooldown?: number;
  enabled?: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}