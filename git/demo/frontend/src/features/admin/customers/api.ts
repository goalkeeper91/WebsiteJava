// features/admin/customers/api.ts
import type { AdminCustomer, AdminCustomerStats, PaginatedResponse } from './types';

const CUSTOMERS_URL = '/api/admin/customers';

export async function getAdminCustomers(
  page = 1,
  pageSize = 20
): Promise<PaginatedResponse<AdminCustomer>> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  const res = await fetch(`${CUSTOMERS_URL}?${params}`, { credentials: 'include' });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Kundenliste');
  }

  return res.json();
}

export async function getAdminCustomerStats(): Promise<AdminCustomerStats> {
  const res = await fetch(`${CUSTOMERS_URL}/stats`, { credentials: 'include' });

  if (!res.ok) {
    throw new Error('Fehler beim Laden der Kundenstatistik');
  }

  return res.json();
}

export async function setAdminCustomerTier(twitchId: string, tierId: string): Promise<void> {
  const res = await fetch(`${CUSTOMERS_URL}/${twitchId}/tier`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ tier_id: tierId }),
  });

  if (!res.ok) {
    const error = await res.text();
    throw new Error(error || 'Fehler beim Ändern des Tarifs');
  }
}
