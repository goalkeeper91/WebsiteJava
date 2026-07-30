// features/cs2/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type {
  CS2CasterSettings,
  CS2CasterSettingsUpdateRequest,
  CS2LiveStatus,
  CS2Note,
  CS2NoteCreateRequest,
  CS2NoteUpdateRequest,
} from './types';

const BASE_URL = '/api/dashboard/cs2';

export async function getCS2Settings(): Promise<CS2CasterSettings> {
  const res = await apiFetch(`${BASE_URL}/settings`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der CS2-Einstellungen');
  }

  return res.json();
}

export async function updateCS2Settings(data: CS2CasterSettingsUpdateRequest): Promise<CS2CasterSettings> {
  const res = await apiFetch(`${BASE_URL}/settings`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Speichern der CS2-Einstellungen');
  }

  return res.json();
}

export async function getCS2LiveStatus(): Promise<CS2LiveStatus> {
  const res = await apiFetch(`${BASE_URL}/live-status`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden des Live-Status');
  }

  return res.json();
}

export async function getCS2Notes(): Promise<CS2Note[]> {
  const res = await apiFetch(`${BASE_URL}/notes`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Notizen');
  }

  return res.json();
}

export async function createCS2Note(data: CS2NoteCreateRequest): Promise<CS2Note> {
  const res = await apiFetch(`${BASE_URL}/notes`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Anlegen der Notiz');
  }

  return res.json();
}

export async function updateCS2Note(id: number, data: CS2NoteUpdateRequest): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/notes/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Aktualisieren der Notiz');
  }
}

export async function deleteCS2Note(id: number): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/notes/${id}`, {
    method: 'DELETE',
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Löschen der Notiz');
  }
}
