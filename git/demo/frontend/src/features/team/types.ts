// features/team/types.ts

export interface TeamMember {
  id: number;
  owner_twitch_id: string;
  member_twitch_id: string;
  member_login: string;
  created_at: string;
}

export interface ManagedChannel {
  owner_twitch_id: string;
  owner_login: string;
}
