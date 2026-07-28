// lib/paddleSdk.ts
//
// Wraps @paddle/paddle-js's initializePaddle() for the public pricing page
// (Pricing.tsx / usePricePreview.ts). Deliberately separate from
// lib/paddle.ts (the hand-rolled script-tag loader used by the existing
// logged-in SubscriptionDashboard.tsx from SaaS Phase 1) - that flow already
// works and isn't touched here to avoid regressing it.

import { initializePaddle, type Paddle, type Environments } from "@paddle/paddle-js";

let paddleInstance: Promise<Paddle | undefined> | null = null;

// Never silently defaults the environment - an unset VITE_PADDLE_ENVIRONMENT
// throws immediately rather than risking a build accidentally running
// against the wrong Paddle account (sandbox vs production).
function requireEnvironment(): Environments {
  const env = import.meta.env.VITE_PADDLE_ENVIRONMENT;
  if (env !== "sandbox" && env !== "production") {
    throw new Error(
      `VITE_PADDLE_ENVIRONMENT muss "sandbox" oder "production" sein, ist aber ${JSON.stringify(env)}`
    );
  }
  return env;
}

export function getPaddle(): Promise<Paddle | undefined> {
  if (paddleInstance) return paddleInstance;

  const token = import.meta.env.VITE_PADDLE_CLIENT_TOKEN;
  if (!token) {
    throw new Error("VITE_PADDLE_CLIENT_TOKEN ist nicht konfiguriert");
  }

  paddleInstance = initializePaddle({
    token,
    environment: requireEnvironment(),
  });

  return paddleInstance;
}
