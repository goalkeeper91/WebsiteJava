import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import Seo from "../../components/Seo";

type Visual =
  | { type: "image"; src: string }
  | { type: "icon"; icon: React.ReactNode; gradient: string };

type Project = {
  title: string;
  category: string;
  summary: string;
  highlights: string[];
  tech: string[];
  link?: { href: string; label: string };
  visual: Visual;
};

// Portfolio content is intentionally conservative: only capability-level
// descriptions (architecture, scope, tech stack), never business-logic
// specifics or client-confidential details - matches the user's own
// instruction ("so viel wie möglich zeigen, ohne wichtige Kernfeatures zu
// leaken"). All four visuals are real: Punishers Germany's brand logo, a
// real Twitch Bot dashboard screenshot, and fresh screenshots taken of
// Hanse-Analyst (its own dev server, browser DOM) and FaceitReader (the
// actual installed desktop app, OS-level window capture) - deliberately
// NOT the very first FaceitReader screenshot that was here, which
// happened to catch an in-app error dialog. The `icon` variant of Visual
// below is kept as a documented fallback pattern (styled gradient+icon
// panel) for a future project that doesn't have a clean screenshot yet,
// even though nothing currently uses it.
const projects: Project[] = [
  {
    title: "Punishers Germany",
    category: "🌐 Website & Plattform",
    summary:
      "Öffentliche Website + Admin-Backend für einen deutschen Esport-Verein: Team-/Spieler-Verwaltung, automatischer Match-Abgleich mit der FACEIT-API, rollenbasiertes Rechtesystem und ein selbst gehostetes KI-Feature für Social-Media-Post-Entwürfe.",
    highlights: [
      "Django + FastAPI Backend (ASGI, ein Prozess), React-Router-SSR-Frontend",
      "FACEIT-API-Integration: automatischer Sync von Matches, Ligen, Spieler-Stats",
      "KI-gestützte Social-Media-Post-Entwürfe über ein selbst gehostetes LLM (Ollama)",
      "Produktiv im Einsatz, automatisierte Deploy-Pipeline (GitHub Actions → Docker Compose)",
    ],
    tech: ["Python", "Django", "FastAPI", "React", "TypeScript", "PostgreSQL", "Docker"],
    link: { href: "https://punishersgermany.de", label: "Live-Seite ansehen" },
    visual: { type: "image", src: "/images/punishers-germany-logo.png" },
  },
  {
    title: "Twitch-Bot-SaaS-Plattform",
    category: "⚙️ SaaS-Produkt",
    summary:
      "Diese Plattform selbst: ein SaaS-Produkt für Twitch-Streamer (Automod, Loyalty-System, Clip-Automatisierung) plus Freelance-Business-Auftritt, sauber als eigene Storefronts voneinander getrennt.",
    highlights: [
      "Go-Backend nach Clean-/Hexagonal-Architecture-Prinzipien",
      "Twitch-OAuth2-Login, AES-GCM-verschlüsselte Token-Speicherung",
      "Paddle-Billing-Integration für Abo-Tarife (Merchant-of-Record)",
      "Domain-basiertes Multi-Storefront-Routing (dieselbe Seite, die du gerade siehst)",
    ],
    tech: ["Go", "React", "TypeScript", "PostgreSQL", "Redis", "Docker", "Paddle"],
    visual: { type: "image", src: "/images/Twitch Bot Overview.png" },
  },
  {
    title: "Hanse-Analyst",
    category: "🧠 KI-Automatisierung",
    summary:
      "Full-Stack-Anwendung zum Hochladen und KI-gestützten Analysieren von PDF-Dokumenten (Rechnungen, Verträge) - komplett lokal verarbeitet, keine Daten verlassen die eigene Infrastruktur.",
    highlights: [
      "Automatische Rechnungsprüfung gegen ein internes Bestellsystem",
      "Automatische Dokumenttyp-Erkennung per lokalem LLM",
      "Split-Screen-UI: Dokument-Vorschau + interaktiver KI-Chat",
      "Läuft komplett lokal über Ollama (Llama 3) - keine Cloud-LLM-API",
    ],
    tech: ["Python", "FastAPI", "SQLAlchemy", "React", "Vite", "Ollama / Llama 3"],
    link: { href: "https://github.com/goalkeeper91/Hanse-Analyst", label: "Auf GitHub ansehen" },
    visual: { type: "image", src: "/images/hanse-analyst-screenshot.png" },
  },
  {
    title: "FaceitReader",
    category: "🖥️ Desktop-Tool",
    summary:
      "Desktop-Tool zur Administration von FACEIT-basierten SkinBaron-Cups - automatisiert wiederkehrende Turnier-Verwaltungsaufgaben, die sonst manuell in der FACEIT-Oberfläche erledigt werden müssten.",
    highlights: [
      "Native Desktop-Anwendung (C#/.NET, Visual Studio)",
      "FACEIT-API-Anbindung für Cup-/Match-Verwaltung",
    ],
    tech: ["C#", ".NET"],
    link: { href: "https://github.com/goalkeeper91/FaceitReader", label: "Auf GitHub ansehen" },
    visual: { type: "image", src: "/images/faceitreader-screenshot.png" },
  },
];

const DevPortfolio = () => {
  return (
    <div className="relative w-full bg-slate-950 text-white overflow-hidden">
      <Seo
        title="Portfolio & Referenzen"
        description="Ausgewählte Projekte von Marcel Turlach: Esport-Vereins-Plattform mit FACEIT-Integration und KI-Features, eine Twitch-Bot-SaaS-Plattform mit Go-Backend und Paddle-Billing, und mehr."
        path="/portfolio"
      />

      {/* Same decorative glow as the rest of the storefront. */}
      <div className="pointer-events-none absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[900px] rounded-full bg-goalyBlue/10 blur-3xl" />

      <section className="relative z-10 text-center py-20 px-6">
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

      <section className="relative z-10 py-16 px-6 max-w-6xl mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
          {projects.map((project, i) => (
            <motion.div
              key={project.title}
              className="bg-slate-900 rounded-2xl shadow-lg overflow-hidden flex flex-col"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-50px" }}
              transition={{ delay: i * 0.1 }}
              whileHover={{ y: -4 }}
            >
              {/* Visual header - real screenshot/logo where one exists and
                  looks clean, otherwise a styled icon panel instead of
                  nothing (or worse, a misleading screenshot). */}
              <div className="relative h-44 w-full overflow-hidden">
                {project.visual.type === "image" ? (
                  <img
                    src={project.visual.src}
                    alt={project.title}
                    loading="lazy"
                    className="w-full h-full object-cover object-top"
                  />
                ) : (
                  <div className={`w-full h-full flex items-center justify-center bg-gradient-to-br ${project.visual.gradient} text-white/90`}>
                    {project.visual.icon}
                  </div>
                )}
                <span className="absolute top-3 left-3 bg-slate-950/80 backdrop-blur text-xs font-semibold px-3 py-1.5 rounded-full">
                  {project.category}
                </span>
              </div>

              <div className="p-6 flex flex-col grow">
                <div className="flex flex-wrap items-center justify-between gap-2 mb-3">
                  <h2 className="text-xl font-bold">{project.title}</h2>
                  {project.link && (
                    <a
                      href={project.link.href}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-xs px-3 py-1.5 bg-goalyBlue hover:bg-goalyCyan rounded-full font-semibold transition whitespace-nowrap"
                    >
                      {project.link.label}
                    </a>
                  )}
                </div>
                <p className="text-gray-300 text-sm mb-4">{project.summary}</p>

                <ul className="text-sm text-gray-400 space-y-1.5 mb-4">
                  {project.highlights.map((h) => (
                    <li key={h} className="flex gap-2">
                      <span className="text-goalyCyan">›</span> {h}
                    </li>
                  ))}
                </ul>

                <div className="flex flex-wrap gap-2 mt-auto pt-2">
                  {project.tech.map((t) => (
                    <span key={t} className="text-xs bg-slate-800 border border-goalyBlue/30 text-gray-300 px-2.5 py-1 rounded-full">
                      {t}
                    </span>
                  ))}
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </section>

      <section className="relative z-10 bg-slate-900 py-20 px-6 text-center">
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
