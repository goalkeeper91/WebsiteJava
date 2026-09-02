// dev.goalkeeper91.de is the extracted software-development-business
// storefront - carries ONLY freelance/dev content (Home, Services,
// Portfolio, About/CV, Contact), no streamer/gaming/community content at
// all. Mirrors botDomain.ts's pattern exactly (same reasoning: a single
// constant deciding the restricted hostname, used both for route-gating in
// App.tsx and nav trimming in Navlinks.tsx/Footer.tsx), started as a
// subdomain of the main site but written host-name-driven and content-
// self-contained specifically so it can move to its own independent domain
// later without a rewrite - just update this constant (and, when that day
// comes, point the new domain at the same frontend_dist build).
export const DEV_STOREFRONT_HOSTNAME = "dev.goalkeeper91.de";

export function isDevStorefront(): boolean {
  return typeof window !== "undefined" && window.location.hostname === DEV_STOREFRONT_HOSTNAME;
}
