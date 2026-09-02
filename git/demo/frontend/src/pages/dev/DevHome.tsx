import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import Seo from "../../components/Seo";

// Landing page for dev.goalkeeper91.de - the extracted software-dev-
// business storefront (see lib/devDomain.ts). Deliberately carries no
// streamer/gaming framing at all, unlike the main domain's Hero - a
// visitor here is evaluating whether to hire a developer, not watching a
// stream.
const DevHome = () => {
  return (
    <section className="relative w-full min-h-screen py-7 bg-slate-950 flex items-center justify-center">
      <Seo
        title="Marcel Turlach - Fullstack-Entwicklung & KI-Automatisierung"
        description="Individuelle Softwarelösungen, Backend-Architektur und KI-gestützte Prozessautomatisierung - Freelance-Entwicklung von Marcel Turlach."
        path="/"
      />
      <div className="relative z-20 text-center text-white w-full px-4 sm:px-6 md:px-8">
        <motion.h1
          className="text-2xl sm:text-4xl md:text-5xl lg:text-6xl font-extrabold mb-6 leading-snug hyphens-auto"
          initial={{ opacity: 0, y: -30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
        >
          Fullstack-Entwicklung & KI-Automatisierung
        </motion.h1>

        <motion.p
          className="text-base sm:text-lg md:text-xl lg:text-2xl mb-10 text-white hyphens-auto max-w-3xl mx-auto"
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
          className="flex flex-wrap justify-center gap-4"
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
      </div>
    </section>
  );
};

export default DevHome;
