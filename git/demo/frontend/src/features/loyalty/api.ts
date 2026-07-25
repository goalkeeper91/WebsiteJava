// features/loyalty/api.ts
import type { LoyaltySettings, UpdateLoyaltySettingsRequest, LeaderboardEntry, PaginatedResponse } from './types';

const BASE_URL = '/api/dashboard/loyalty';

export async function getLoyaltySettings(): Promise<LoyaltySettings> {
  const res = await fetch(BASE_URL, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Loyalty-Einstellungen');
  }

  return res.json();
}

export async function updateLoyaltySettings(data: UpdateLoyaltySettingsRequest): Promise<LoyaltySettings> {
  const res = await fetch(BASE_URL, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Speichern');
  }

  return res.json();
}

export async function getLoyaltyLeaderboard(page = 1, pageSize = 20): Promise<PaginatedResponse<LeaderboardEntry>> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const res = await fetch(`${BASE_URL}/leaderboard?${params}`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Loyalty-Bestenliste');
  }

  return res.json();
}
