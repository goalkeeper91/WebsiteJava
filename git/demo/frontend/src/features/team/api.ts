// features/team/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type { TeamMember, ManagedChannel } from './types';

const BASE_URL = '/api/dashboard/team';

export async function getTeamMembers(): Promise<TeamMember[]> {
  const res = await apiFetch(`${BASE_URL}/members`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Team-Mitglieder');
  }

  return res.json();
}

export async function inviteTeamMember(login: string): Promise<TeamMember> {
  const res = await apiFetch(`${BASE_URL}/invite`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login }),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Einladen');
  }

  return res.json();
}

export async function removeTeamMember(memberTwitchId: string): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/members/${memberTwitchId}`, {
    method: 'DELETE',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Entfernen des Team-Mitglieds');
  }
}

export async function leaveTeam(ownerTwitchId: string): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/leave/${ownerTwitchId}`, {
    method: 'DELETE',
  });

  if (!res.ok) {
    throw new Error('Fehler beim Verlassen des Teams');
  }
}

export async function getManagedChannels(): Promise<ManagedChannel[]> {
  const res = await apiFetch(`${BASE_URL}/managed-channels`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der verwaltbaren Kanäle');
  }

  return res.json();
}
