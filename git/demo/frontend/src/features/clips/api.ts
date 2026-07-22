// features/clips/api.ts
import type {
  AutomationSettings,
  UpdateAutomationSettingsRequest,
  ClipLog,
  ClipLogStats,
  PaginatedResponse,
} from './types';

const SETTINGS_URL = '/api/automation/settings';
const CLIPS_URL = '/api/clips';

export async function getAutomationSettings(): Promise<AutomationSettings> {
  const res = await fetch(SETTINGS_URL, { credentials: 'include' });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Laden der Automation Settings');
  }

  return res.json();
}

export async function updateAutomationSettings(
  data: UpdateAutomationSettingsRequest
): Promise<AutomationSettings> {
  const res = await fetch(SETTINGS_URL, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Speichern der Automation Settings');
  }

  return res.json();
}

export async function getClips(
  page = 1,
  pageSize = 20
): Promise<PaginatedResponse<ClipLog>> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const res = await fetch(`${CLIPS_URL}?${params}`, { credentials: 'include' });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Clips');
  }

  return res.json();
}

export async function getClipStats(): Promise<ClipLogStats> {
  const res = await fetch(`${CLIPS_URL}/stats`, { credentials: 'include' });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Clip-Statistiken');
  }

  return res.json();
}
