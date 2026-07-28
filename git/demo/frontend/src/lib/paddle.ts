// lib/paddle.ts
//
// Lazily loads Paddle.js (Paddle Billing's client-side checkout SDK) and
// initializes it once. Loaded on demand rather than in index.html since
// only the subscription page ever needs it.
//
// customData passed into Checkout.open() is copied by Paddle onto the
// resulting transaction and then the subscription itself, and propagates to
// every future renewal transaction - this is how the backend's webhook
// handler (PaddleService.resolveTwitchUserID) maps an incoming Paddle event
// back to a local user without needing a pre-created customer mapping.

declare global {
  interface Window {
    Paddle?: {
      Environment: {
        set: (env: "sandbox" | "production") => void;
      };
      Initialize: (opts: { token: string; eventCallback?: (event: { name: string }) => void }) => void;
      Checkout: {
        open: (opts: {
          items: { priceId: string; quantity: number }[];
          customData?: Record<string, string>;
        }) => void;
      };
    };
  }
}

let loadPromise: Promise<void> | null = null;

function loadScript(): Promise<void> {
  if (loadPromise) return loadPromise;

  loadPromise = new Promise((resolve, reject) => {
    if (window.Paddle) {
      resolve();
      return;
    }
    const script = document.createElement("script");
    script.src = "https://cdn.paddle.com/paddle/v2/paddle.js";
    script.async = true;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error("Paddle.js konnte nicht geladen werden"));
    document.head.appendChild(script);
  });

  return loadPromise;
}

let initialized = false;
let onCheckoutCompleted: (() => void) | null = null;

// setCheckoutCompletedHandler lets the subscription page react once the
// overlay checkout finishes - the actual tier change still only happens
// once Paddle's webhook lands (this is best-effort UX so the page can
// re-poll/reload sooner rather than the user needing to refresh manually).
export function setCheckoutCompletedHandler(handler: (() => void) | null) {
  onCheckoutCompleted = handler;
}

async function ensurePaddleInitialized(): Promise<void> {
  await loadScript();
  if (initialized) return;

  const token = import.meta.env.VITE_PADDLE_CLIENT_TOKEN;
  if (!token) {
    throw new Error("VITE_PADDLE_CLIENT_TOKEN ist nicht konfiguriert");
  }

  const environment = import.meta.env.VITE_PADDLE_ENVIRONMENT || "production";
  if (environment === "sandbox") {
    window.Paddle?.Environment.set("sandbox");
  }

  window.Paddle?.Initialize({
    token,
    eventCallback: (event) => {
      if (event.name === "checkout.completed") {
        onCheckoutCompleted?.();
      }
    },
  });
  initialized = true;
}

// getPaddlePriceID maps a tier + billing cycle to the Paddle price_id
// configured for it - mirrors the backend's PaddleConfig.TierForPriceID,
// just in the opposite direction (tier -> price_id instead of price_id ->
// tier). Price IDs aren't secret - Paddle.js needs them client-side to open
// a checkout - so they're plain build-time env vars, same as the client token.
export function getPaddlePriceID(tier: "pro" | "premium", cycle: "monthly" | "yearly"): string | undefined {
  const key = `VITE_PADDLE_PRICE_${tier.toUpperCase()}_${cycle.toUpperCase()}`;
  return (import.meta.env as Record<string, string | undefined>)[key];
}

export async function openPaddleCheckout(priceId: string, twitchUserId: string): Promise<void> {
  await ensurePaddleInitialized();
  window.Paddle?.Checkout.open({
    items: [{ priceId, quantity: 1 }],
    customData: { twitch_user_id: twitchUserId },
  });
}
