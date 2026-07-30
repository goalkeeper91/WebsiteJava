// features/cs2/types.ts

export interface CS2CasterSettings {
  user_twitch_id: string;
  gsi_token: string;
  predictions_enabled: boolean;
  multikill_announce_enabled: boolean;
  map_end_announce_enabled: boolean;
  title_update_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface CS2CasterSettingsUpdateRequest {
  predictions_enabled?: boolean;
  multikill_announce_enabled?: boolean;
  map_end_announce_enabled?: boolean;
  title_update_enabled?: boolean;
}

export type CS2NoteSubjectType = "team" | "player";

export interface CS2Note {
  id: number;
  user_twitch_id: string;
  subject_type: CS2NoteSubjectType;
  subject_name: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CS2NoteCreateRequest {
  subject_type: CS2NoteSubjectType;
  subject_name: string;
  content: string;
}

export interface CS2NoteUpdateRequest {
  subject_name?: string;
  content?: string;
}

export interface CS2LiveStatus {
  active: boolean;
  observed_player_name?: string;
  team_ct_name?: string;
  team_t_name?: string;
  score_ct: number;
  score_t: number;
  map_name?: string;
  team_ct_players?: string[];
  team_t_players?: string[];
}
