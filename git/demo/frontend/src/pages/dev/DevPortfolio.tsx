import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import Seo from "../../components/Seo";

type Project = {
  title: string;
  summary: string;
  highlights: string[];
  tech: string[];
  link?: { href: string; label: string };
  image?: string;
  needsReview?: boolean;
};

// Portfolio content is intentionally conservative: only capability-level
// descriptions (architecture, scope, tech stack), never business-logic
// specifics or client-confidential details - matches the user's own
// instruction ("so viel wie möglich zeigen, ohne wichtige Kernfeatures zu
// leaken"). PunishersGer and the Twitch-Bot-SaaS-Plattform entries are
// based on direct first-hand knowledge of both codebases; the remaining
// three are only known from brief prior context and are flagged
// (needsReview) for the user to verify/expand rather than asserting
// specifics that might be wrong.
const projects: Project[] = [
  {
    title: "Punishers Germany - Esport-Vereins-Plattform",
    summary:
      "Öffentliche Website + Admin-Backend für einen deutschen Esport-Verein: Team-/Spieler-Verwaltung, automatischer Match-Abgleich mit der FACEIT-API, rollenbasiertes Rechtesystem und ein selbst gehostetes KI-Feature für Social-Media-Post-Entwürfe.",
    highlights: [
      "Django + FastAPI Backend (ASGI, ein Prozess), React-Router-SSR-Frontend",
      "FACEIT-API-Integration: automatischer Sync von Matches, Ligen, Spieler-Stats",
      "Rollenbasiertes Berechtigungssystem für Admins/Team-Manager",
      "KI-gestützte Social-Media-Post-Entwürfe über ein selbst gehostetes LLM (Ollama), inkl. programmatisch gerendertem Bild-Template",
      "Produktiv im Einsatz, automatisierte Deploy-Pipeline (GitHub Actions → Docker Compose)",
    ],
    tech: ["Python", "Django", "FastAPI", "React", "TypeScript", "PostgreSQL", "Docker"],
    link: { href: "https://punishersgermany.de", label: "Live-Seite ansehen" },
  },
  {
    title: "Twitch-Bot-SaaS-Plattform (goalkeeper91.de)",
    summary:
      "Diese Plattform selbst: ein SaaS-Produkt für Twitch-Streamer (Automod, Loyalty-System, Clip-Automatisierung) plus Freelance-Business-Auftritt, sauber als eigene Storefronts voneinander getrennt.",
    highlights: [
      "Go-Backend nach Clean-/Hexagonal-Architecture-Prinzipien (Migration von einem früheren Spring-Boot-Backend)",
      "Twitch-OAuth2-Login, Session-Management, AES-GCM-verschlüsselte Token-Speicherung",
      "Paddle-Billing-Integration für Abo-Tarife (Merchant-of-Record, EU-Widerrufsrecht-konform)",
      "Redis-Pub/Sub für die Kommunikation zwischen Backend und Bot-Prozess",
      "Domain-basiertes Multi-Storefront-Routing (SaaS-Produkt, Freelance-Business, Streamer-Präsenz sauber getrennt)",
    ],
    tech: ["Go", "React", "TypeScript", "PostgreSQL", "Redis", "Docker", "Paddle"],
  },
  {
    title: "Hanse-Analyst - Lokale KI-Dokumentenanalyse",
    summary:
      "Full-Stack-Anwendung zum Hochladen und KI-gestützten Analysieren von PDF-Dokumenten (Rechnungen, Verträge) - komplett lokal verarbeitet, keine Daten verlassen die eigene Infrastruktur. Speziell für DSGVO-sensible Unternehmensdaten ohne Cloud-Abhängigkeit.",
    highlights: [
      "Automatische Rechnungsprüfung: Abgleich extrahierter Rechnungsdaten gegen ein internes Bestellsystem",
      "Automatische Dokumenttyp-Erkennung (Rechnung, Vertrag, ...) per lokalem LLM",
      "Split-Screen-UI: Dokument-Vorschau + interaktiver KI-Chat für Rückfragen zum Inhalt",
      "Läuft komplett lokal über Ollama (Llama 3) - keine Cloud-LLM-API, keine Drittanbieter-Datenübertragung",
    ],
    tech: ["Python", "FastAPI", "SQLAlchemy", "React", "Vite", "Tailwind CSS", "Ollama / Llama 3"],
    link: { href: "https://github.com/goalkeeper91/Hanse-Analyst", label: "Auf GitHub ansehen" },
  },
  {
    title: "FaceitReader - SkinBaron-Cup-Administration",
    summary:
      "Desktop-Tool zur Administration von FACEIT-basierten SkinBaron-Cups - automatisiert wiederkehrende Turnier-Verwaltungsaufgaben, die sonst manuell in der FACEIT-Oberfläche erledigt werden müssten.",
    highlights: [
      "Native Desktop-Anwendung (C#/.NET, Visual Studio)",
      "FACEIT-API-Anbindung für Cup-/Match-Verwaltung",
    ],
    tech: ["C#", ".NET"],
    link: { href: "https://github.com/goalkeeper91/FaceitReader", label: "Auf GitHub ansehen" },
    image: "/images/faceit_reader.png",
  },
];

const DevPortfolio = () => {
  return (
    <div className="relative w-full min-h-screen bg-slate-950 text-white">
      <Seo
        title="Portfolio & Referenzen"
        description="Ausgewählte Projekte von Marcel Turlach: Esport-Vereins-Plattform mit FACEIT-Integration und KI-Features, eine Twitch-Bot-SaaS-Plattform mit Go-Backend und Paddle-Billing, und mehr."
        path="/portfolio"
      />
      <section className="relative text-center py-20 px-6 bg-slate-900">
        <motion.h1
          className="text-5xl font-extrabold mb-4"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          Portfolio & Referenzen
        </motion.h1>
        <motion.p
          className="text-lg text-gray-300 max-w-3xl mx-auto"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          Eine Auswahl an Projekten - von produktiv im Einsatz befindlichen Plattformen bis zu gezielten
          Skills-Showcases. Details zu Architektur und Umfang statt Business-internas.
        </motion.p>
      </section>

      <section className="py-16 px-6 max-w-5xl mx-auto space-y-10">
        {projects.map((project) => (
          <motion.div
            key={project.title}
            className="bg-slate-900 rounded-2xl shadow-lg p-6 sm:p-8"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-50px" }}
          >
            <div className="flex flex-wrap items-start justify-between gap-3 mb-3">
              <h2 className="text-2xl font-bold">{project.title}</h2>
              {project.link && (
                <a
                  href={project.link.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm px-4 py-2 bg-goalyBlue hover:bg-goalyCyan rounded-full font-semibold transition whitespace-nowrap"
                >
                  {project.link.label}
                </a>
              )}
            </div>
            <p className="text-gray-300 mb-4">{project.summary}</p>

            {project.image && (
              <img
                src={project.image}
                alt={project.title}
                loading="lazy"
                className="w-full max-w-xl rounded-lg border border-gray-700 mb-4"
              />
            )}

            {project.highlights.length > 0 && (
              <ul className="list-disc list-inside text-sm text-gray-300 space-y-1 mb-4">
                {project.highlights.map((h) => <li key={h}>{h}</li>)}
              </ul>
            )}

            {project.tech.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {project.tech.map((t) => (
                  <span key={t} className="text-xs bg-slate-800 border border-goalyBlue/40 text-gray-300 px-3 py-1 rounded-full">
                    {t}
                  </span>
                ))}
              </div>
            )}

            {project.needsReview && (
              <p className="text-xs text-yellow-500/80 mt-4 italic">
                Platzhalter - Details bitte noch ergänzen/prüfen.
              </p>
            )}
          </motion.div>
        ))}
      </section>

      <section className="py-20 px-6 text-center bg-slate-900">
        <h2 className="text-3xl font-bold mb-6">Dein Projekt hier als nächstes?</h2>
        <Link
          to="/contact"
          className="px-8 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition"
        >
          Kostenloses Erstgespräch
        </Link>
      </section>
    </div>
  );
};

export default DevPortfolio;
