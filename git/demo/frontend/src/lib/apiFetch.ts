// lib/apiFetch.ts
//
// Thin wrapper around fetch() used by every dashboard feature's api.ts.
// Its one job: inject the X-Acting-As header when the logged-in user is
// currently managing another channel's dashboard (Team-Zugriff feature),
// and fall back to "acting as myself" automatically if the backend ever
// returns 403 (e.g. access was just revoked mid-session).
//
// State lives here as a module-level singleton rather than in React
// context, since api.ts functions are plain functions (not components) and
// can't use hooks. TeamContext syncs its React state into this module via
// setActingAsChannel() and listens for revocation via setForbiddenHandler().

let currentActingAs: string | null = null;
let onForbidden: (() => void) | null = null;

export function setActingAsChannel(twitchId: string | null) {
  currentActingAs = twitchId;
}

export function setForbiddenHandler(handler: (() => void) | null) {
  onForbidden = handler;
}

export async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
  const headers = new Headers(options.headers);
  if (currentActingAs) {
    headers.set('X-Acting-As', currentActingAs);
  }

  const res = await fetch(url, { ...options, credentials: 'include', headers });

  if (res.status === 403 && currentActingAs) {
    currentActingAs = null;
    onForbidden?.();
  }

  return res;
}
