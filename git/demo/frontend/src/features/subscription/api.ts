// features/subscription/api.ts
import { apiFetch } from '../../lib/apiFetch';
import type { SubscriptionTier, UserSubscription } from './types';

const BASE_URL = '/api/subscription';

export async function getSubscription(): Promise<UserSubscription> {
  const res = await apiFetch(BASE_URL);

  if (!res.ok) {
    throw new Error('Fehler beim Laden des Abos');
  }

  return res.json();
}

export async function getTiers(): Promise<SubscriptionTier[]> {
  const res = await apiFetch(`${BASE_URL}/tiers`);

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Tarife');
  }

  return res.json();
}

// downgradeToFree is the only tier change this app still performs directly -
// any paid tier requires a real Paddle checkout (see openCheckout below), the
// backend rejects a direct POST to a paid tier_id.
export async function downgradeToFree(): Promise<UserSubscription> {
  const res = await apiFetch(`${BASE_URL}/upgrade`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tier_id: 'free' }),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Wechsel auf Free');
  }

  return res.json();
}

export async function cancelSubscription(): Promise<void> {
  const res = await apiFetch(`${BASE_URL}/cancel`, { method: 'POST' });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Kündigen');
  }
}

// getPortalLink mints a fresh, single-use Paddle Customer Portal URL - only
// available once a paid tier has actually been purchased through Paddle
// (i.e. a paddle_customer_id exists on the subscription).
export async function getPortalLink(): Promise<string> {
  const res = await apiFetch(`${BASE_URL}/portal-link`);

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Laden des Kundenportal-Links');
  }

  const data: { url: string } = await res.json();
  return data.url;
}
