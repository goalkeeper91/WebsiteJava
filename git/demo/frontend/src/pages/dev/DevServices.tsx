import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { FaCode, FaRobot, FaCogs } from "react-icons/fa";
import Seo from "../../components/Seo";

type Service = {
  title: string;
  description: string;
  icon: React.ReactNode;
};

// Adapted from the old main-domain ServicesPage.tsx, with the personal
// "als Streamer entwickle ich..." framing dropped - a client here is
// evaluating capability, not stream personality. KI-gestützte
// Prozessautomatisierung gets its own card (was folded into "Twitch Bot &
// Automatisierung" before) since it's now the CV's own headline
// specialization (velpTEC Certified AI Developer, see DevAbout.tsx).
const DevServices = () => {
  const services: Service[] = [
    {
      title: "Individuelle Softwarelösungen",
      description:
        "Von Prototyp bis Produktivbetrieb – maßgeschneiderte Web-, Mobile- oder Desktop-Anwendungen für deine Anforderungen, kein Baukasten von der Stange.",
      icon: <FaCode size={32} />,
    },
    {
      title: "KI-gestützte Prozessautomatisierung",
      description:
        "Design und Implementierung KI-gestützter Automatisierungslösungen (Digital Automations, LLM-/Prompt-Engineering, API-Integrationen) für datenintensive oder wiederkehrende Aufgaben.",
      icon: <FaRobot size={32} />,
    },
    {
      title: "Bots & Integrationen",
      description:
        "Bots, Workflow-Automatisierungen und Schnittstellen-Integrationen (Discord, Twitch, interne Tools) - Prozesse automatisieren statt manuell wiederholen.",
      icon: <FaCogs size={32} />,
    },
  ];

  return (
    <div className="relative w-full bg-slate-950 text-white overflow-hidden">
      <Seo
        title="Leistungen: Softwareentwicklung & KI-Automatisierung"
        description="Individuelle Softwarelösungen, KI-gestützte Prozessautomatisierung und Bot-/API-Integrationen - Freelance-Entwicklung von Marcel Turlach."
        path="/services"
      />

      {/* Same decorative glow as About/Home/Portfolio for a consistent look. */}
      <div className="pointer-events-none absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[900px] rounded-full bg-goalyBlue/10 blur-3xl" />

      <section className="relative z-10 text-center py-20 px-6">
        <motion.h1
          className="text-5xl font-extrabold mb-4"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          Meine Leistungen
        </motion.h1>
        <motion.p
          className="text-lg text-gray-300 max-w-3xl mx-auto mb-8"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          Ich bin kein „Alleskönner" – und das ist gut so. Statt alles ein bisschen zu machen,
          konzentriere ich mich darauf, Probleme zu verstehen und die passende Lösung zu entwickeln.
        </motion.p>
      </section>

      <section className="relative z-10 py-16 px-6 max-w-6xl mx-auto">
        <h2 className="sr-only">Leistungen im Überblick</h2>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {services.map((service, i) => (
            <motion.div
              key={service.title}
              className="bg-slate-900 rounded-2xl shadow-lg p-6"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-50px" }}
              transition={{ delay: i * 0.1 }}
              whileHover={{ scale: 1.03 }}
            >
              <div className="text-goalyBlue mb-3">{service.icon}</div>
              <h3 className="text-xl font-bold mb-2">{service.title}</h3>
              <p className="text-gray-300 text-sm">{service.description}</p>
            </motion.div>
          ))}
        </div>
      </section>

      <section className="relative z-10 bg-slate-900 py-20 px-6 text-center">
        <motion.h2
          className="text-3xl font-bold mb-6"
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
        >
          Bereit, dein Projekt zu starten?
        </motion.h2>
        <p className="text-gray-300 mb-8 max-w-xl mx-auto">
          Lass uns über deine Ideen sprechen – ob kleine Automatisierung oder komplexe Softwarelösung.
          Keine Pauschalpreise, jedes Projekt bekommt nach einem kurzen Erstgespräch ein individuelles Angebot.
        </p>
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

export default DevServices;
