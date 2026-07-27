// features/giveaways/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type { Giveaway, GiveawayStatusResponse, PaginatedResponse, StartGiveawayRequest } from './types';

const BASE_URL = '/api/dashboard/giveaways';

export async function getGiveawayStatus(): Promise<GiveawayStatusResponse> {
  const res = await apiFetch(`${BASE_URL}/status`);

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

  const res = await apiFetch(`${BASE_URL}/history?${params}`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Giveaway-Historie');
  }

  return res.json();
}

export async function startGiveaway(data: StartGiveawayRequest): Promise<Giveaway> {
  const res = await apiFetch(`${BASE_URL}/start`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Starten des Giveaways');
  }

  return res.json();
}

export async function drawGiveawayWinner(): Promise<Giveaway> {
  const res = await apiFetch(`${BASE_URL}/draw`, {
    method: 'POST',
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Ziehen des Gewinners');
  }

  return res.json();
}

export async function cancelGiveaway(): Promise<Giveaway> {
  const res = await apiFetch(`${BASE_URL}/cancel`, {
    method: 'POST',
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Abbrechen des Giveaways');
  }

  return res.json();
}
