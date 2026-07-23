// features/automod/api.ts
import type { AutomodSettings, UpdateAutomodSettingsRequest } from './types';

const BASE_URL = '/api/dashboard/automod';

export async function getAutomodSettings(): Promise<AutomodSettings> {
  const res = await fetch(BASE_URL, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Automod-Einstellungen');
  }

  return res.json();
}

export async function updateAutomodSettings(data: UpdateAutomodSettingsRequest): Promise<AutomodSettings> {
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
