// features/commands/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type { ChatCommand, CreateCommandRequest, UpdateCommandRequest } from './types';

const BASE_URL = '/api/dashboard/commands';

export async function getCommands(
  page = 1,
  pageSize = 20,
  search?: string,
  enabled?: boolean
): Promise<{ data: ChatCommand[]; total: number; page: number; pageSize: number; totalPages: number }> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  if (search) params.append('search', search);
  if (enabled !== undefined) params.append('enabled', enabled.toString());

  const res = await apiFetch(`${BASE_URL}?${params}`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Commands');
  }

  return res.json();
}

export async function getCommand(id: number): Promise<ChatCommand> {
  const res = await apiFetch(`${BASE_URL}/${id}`);

  if (!res.ok) {
    throw new Error('Command nicht gefunden');
  }

  return res.json();
}

export async function createCommand(data: CreateCommandRequest): Promise<ChatCommand> {
  const res = await apiFetch(BASE_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Erstellen');
  }

  return res.json();
}

export async function updateCommand(id: number, data: UpdateCommandRequest): Promise<ChatCommand> {
  const res = await apiFetch(`${BASE_URL}/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Aktualisieren');
  }

  return res.json();
}

export async function deleteCommand(id: number): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/${id}`, {
    method: 'DELETE',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Löschen');
  }
}

export async function toggleCommand(id: number, enabled: boolean): Promise<ChatCommand> {
  const res = await apiFetch(`${BASE_URL}/${id}/toggle`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ enabled }),
  });

  if (!res.ok) {
    throw new Error('Fehler beim Aktivieren/Deaktivieren');
  }

  return res.json();
}
