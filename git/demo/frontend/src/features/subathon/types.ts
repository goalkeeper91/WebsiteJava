export interface SubathonEventLogEntry {
  time: string;
  text: string;
}

export interface SubathonState {
  userId: string;
  timeRemaining: number;
  isRunning: boolean;
  targetTimestamp?: number;
  totalSubs: number;
  totalBits: number;
  totalEvents: number;
  initialTime: number;
  subTime: number;
  bitsTime: number;
  eventLog: SubathonEventLogEntry[];
}

export interface UpdateSubathonSettingsRequest {
  initialTime?: number;
  subTime?: number;
  bitsTime?: number;
}
