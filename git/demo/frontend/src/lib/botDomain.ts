// bot.goalkeeper91.de is a Paddle-facing storefront that carries ONLY the
// Twitch Bot SaaS product - no freelance/software-development-services
// content anywhere on it. Paddle rejected the main domain for exactly that
// combination ("Human Services/IT Services or Software Development
// Services" is categorically excluded from their Acceptable Use Policy,
// and their review scans the whole domain, not just the checkout page) -
// this constant is the single place that decides which hostname is the
// restricted storefront, used both for route-gating (App.tsx) and nav
// trimming (Navlinks.tsx/Footer.tsx).
export const BOT_STOREFRONT_HOSTNAME = "bot.goalkeeper91.de";

export function isBotStorefront(): boolean {
  return typeof window !== "undefined" && window.location.hostname === BOT_STOREFRONT_HOSTNAME;
}
