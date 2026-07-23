import type { SubathonState, UpdateSubathonSettingsRequest } from './types';

const BASE_URL = '/api/subathon';

export async function getSubathonState(): Promise<SubathonState> {
  const res = await fetch(`${BASE_URL}/state`, { credentials: 'include' });
  if (!res.ok) throw new Error('Fehler beim Laden des Subathon-Timers');
  return res.json();
}

export async function startSubathon(): Promise<SubathonState> {
  const res = await fetch(`${BASE_URL}/start`, { method: 'POST', credentials: 'include' });
  if (!res.ok) throw new Error('Fehler beim Starten des Timers');
  return res.json();
}

export async function pauseSubathon(): Promise<SubathonState> {
  const res = await fetch(`${BASE_URL}/pause`, { method: 'POST', credentials: 'include' });
  if (!res.ok) throw new Error('Fehler beim Pausieren des Timers');
  return res.json();
}

export async function resetSubathon(): Promise<SubathonState> {
  const res = await fetch(`${BASE_URL}/reset`, { method: 'POST', credentials: 'include' });
  if (!res.ok) throw new Error('Fehler beim Zuruecksetzen des Timers');
  return res.json();
}

export async function updateSubathonSettings(
  data: UpdateSubathonSettingsRequest
): Promise<SubathonState> {
  const res = await fetch(`${BASE_URL}/settings`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Fehler beim Speichern der Einstellungen');
  return res.json();
}
