// features/automod/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type { AutomodSettings, UpdateAutomodSettingsRequest, AutomodEvent, PaginatedResponse } from './types';

const BASE_URL = '/api/dashboard/automod';

export async function getAutomodSettings(): Promise<AutomodSettings> {
  const res = await apiFetch(BASE_URL);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Automod-Einstellungen');
  }

  return res.json();
}

export async function updateAutomodSettings(data: UpdateAutomodSettingsRequest): Promise<AutomodSettings> {
  const res = await apiFetch(BASE_URL, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Speichern');
  }

  return res.json();
}

export async function getAutomodEvents(page = 1, pageSize = 20): Promise<PaginatedResponse<AutomodEvent>> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const res = await apiFetch(`${BASE_URL}/events?${params}`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Automod-Ereignisse');
  }

  return res.json();
}
