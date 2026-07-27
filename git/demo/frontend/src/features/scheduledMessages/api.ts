// features/scheduledMessages/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type {
  ScheduledMessage,
  CreateScheduledMessageRequest,
  UpdateScheduledMessageRequest,
  PaginatedResponse,
} from './types';

const BASE_URL = '/api/dashboard/scheduled-messages';

export async function getScheduledMessages(
  page = 1,
  pageSize = 20
): Promise<PaginatedResponse<ScheduledMessage>> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const res = await apiFetch(`${BASE_URL}?${params}`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der automatisierten Nachrichten');
  }

  return res.json();
}

export async function createScheduledMessage(
  data: CreateScheduledMessageRequest
): Promise<ScheduledMessage> {
  const res = await apiFetch(BASE_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Erstellen');
  }

  return res.json();
}

export async function updateScheduledMessage(
  id: number,
  data: UpdateScheduledMessageRequest
): Promise<ScheduledMessage> {
  const res = await apiFetch(`${BASE_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Aktualisieren');
  }

  return res.json();
}

export async function deleteScheduledMessage(id: number): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/${id}`, {
    method: 'DELETE',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Löschen');
  }
}

export async function toggleScheduledMessage(id: number, enabled: boolean): Promise<ScheduledMessage> {
  const res = await apiFetch(`${BASE_URL}/${id}/toggle`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ enabled }),
  });

  if (!res.ok) {
    throw new Error('Fehler beim Aktivieren/Deaktivieren');
  }

  return res.json();
}
