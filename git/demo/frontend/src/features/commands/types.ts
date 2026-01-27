export interface ChatCommand {
  id: number;
  trigger: string;
  response: string;
  cooldown: number;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PageResponse<T> {
  content: T[];
  totalElements: number;
  totalPages: number;
  number: number;
}
