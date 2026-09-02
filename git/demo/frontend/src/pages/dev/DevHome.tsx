import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { FaCheckCircle, FaLayerGroup, FaRobot, FaAmbulance } from "react-icons/fa";
import Seo from "../../components/Seo";

// Landing page for dev.goalkeeper91.de - the extracted software-dev-
// business storefront (see lib/devDomain.ts). Deliberately carries no
// streamer/gaming framing at all, unlike the main domain's Hero - a
// visitor here is evaluating whether to hire a developer, not watching a
// stream.

const trustBar = [
  "Kostenloses Erstgespräch",
  "Kein Agentur-Overhead",
  "Von Prototyp bis Produktivbetrieb",
];

// Short teaser version of the value props on the About page - just enough
// to earn a click through to /about, not a duplicate of the full pitch.
const teaser = [
  { icon: <FaLayerGroup size={24} />, title: "Breiter Stack", text: "Java, Python, PHP, React, C# - das passende Werkzeug fürs Problem." },
  { icon: <FaRobot size={24} />, title: "KI-Automatisierung", text: "Certified AI Developer, Fokus auf Prozessautomatisierung mit LLMs." },
  { icon: <FaAmbulance size={24} />, title: "Ruhe unter Druck", text: "Rettungsdienst-Hintergrund statt Bootcamp-Gelassenheit." },
];

const DevHome = () => {
  return (
    <div className="relative w-full bg-slate-950 text-white overflow-hidden">
      <Seo
        title="Marcel Turlach - Fullstack-Entwicklung & KI-Automatisierung"
        description="Individuelle Softwarelösungen, Backend-Architektur und KI-gestützte Prozessautomatisierung - Freelance-Entwicklung von Marcel Turlach."
        path="/"
      />

      {/* Same decorative glow as About/Services/Portfolio, so the whole
          storefront reads as one designed site instead of four separate
          pages bolted together. */}
      <div className="pointer-events-none absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[900px] rounded-full bg-goalyBlue/10 blur-3xl" />

      {/* ============ HERO ============ */}
      <section className="relative z-10 min-h-[85vh] flex items-center justify-center py-16">
        <div className="text-center w-full px-4 sm:px-6 md:px-8">
          <motion.h1
            className="text-2xl sm:text-4xl md:text-5xl lg:text-6xl font-extrabold mb-6 leading-snug hyphens-auto"
            initial={{ opacity: 0, y: -30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
          >
            Fullstack-Entwicklung & KI-Automatisierung
          </motion.h1>

          <motion.p
            className="text-base sm:text-lg md:text-xl lg:text-2xl mb-8 text-white hyphens-auto max-w-3xl mx-auto"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.3 }}
          >
            Ich bin Marcel Turlach – Fullstack-Entwickler mit Fokus auf{" "}
            <span className="text-indigo-400 font-semibold">Backend-Architektur</span>,{" "}
            <span className="text-indigo-400 font-semibold">Web-Anwendungen</span> und{" "}
            <span className="text-indigo-400 font-semibold">KI-gestützte Prozessautomatisierung</span>.
            Von Prototyp bis Produktivbetrieb – maßgeschneidert statt von der Stange.
          </motion.p>

          <motion.div
            className="flex flex-wrap justify-center gap-4 mb-8"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
          >
            <Link
              to="/contact"
              className="w-full sm:w-auto px-6 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition text-center"
            >
              Kostenloses Erstgespräch
            </Link>
            <Link
              to="/portfolio"
              className="w-full sm:w-auto px-6 py-3 border-2 border-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition text-center"
            >
              Portfolio ansehen
            </Link>
            <Link
              to="/services"
              className="w-full sm:w-auto px-6 py-3 border-2 border-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition text-center"
            >
              Leistungen
            </Link>
          </motion.div>

          <motion.div
            className="flex flex-wrap justify-center gap-x-6 gap-y-2 text-sm text-gray-400"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.9 }}
          >
            {trustBar.map((item) => (
              <span key={item} className="flex items-center gap-2">
                <FaCheckCircle className="text-goalyCyan" /> {item}
              </span>
            ))}
          </motion.div>
        </div>
      </section>

      {/* ============ WARUM MIT MIR (Teaser) ============ */}
      <section className="relative z-10 bg-slate-900 py-16">
        <div className="max-w-5xl mx-auto px-4">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-8">
            {teaser.map((t, i) => (
              <motion.div
                key={t.title}
                className="bg-slate-800/60 rounded-2xl p-6"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: "-50px" }}
                transition={{ delay: i * 0.1 }}
              >
                <div className="text-goalyBlue mb-3">{t.icon}</div>
                <h3 className="font-bold mb-1">{t.title}</h3>
                <p className="text-sm text-gray-400">{t.text}</p>
              </motion.div>
            ))}
          </div>
          <div className="text-center">
            <Link to="/about" className="text-goalyCyan font-semibold hover:underline">
              Mehr über mich →
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
};

export default DevHome;
