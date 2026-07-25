// features/giveaways/api.ts
import type { Giveaway, GiveawayStatusResponse, PaginatedResponse } from './types';

const BASE_URL = '/api/dashboard/giveaways';

export async function getGiveawayStatus(): Promise<GiveawayStatusResponse> {
  const res = await fetch(`${BASE_URL}/status`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Laden des Giveaway-Status');
  }

  return res.json();
}

export async function getGiveawayHistory(page = 1, pageSize = 20): Promise<PaginatedResponse<Giveaway>> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const res = await fetch(`${BASE_URL}/history?${params}`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Giveaway-Historie');
  }

  return res.json();
}
