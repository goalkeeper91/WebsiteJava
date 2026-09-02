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
      icon: <FaCode size={40} className="text-indigo-400 mb-4" />,
    },
    {
      title: "KI-gestützte Prozessautomatisierung",
      description:
        "Design und Implementierung KI-gestützter Automatisierungslösungen (Digital Automations, LLM-/Prompt-Engineering, API-Integrationen) für datenintensive oder wiederkehrende Aufgaben.",
      icon: <FaRobot size={40} className="text-green-400 mb-4" />,
    },
    {
      title: "Bots & Integrationen",
      description:
        "Bots, Workflow-Automatisierungen und Schnittstellen-Integrationen (Discord, Twitch, interne Tools) - Prozesse automatisieren statt manuell wiederholen.",
      icon: <FaCogs size={40} className="text-pink-400 mb-4" />,
    },
  ];

  return (
    <div className="relative w-full min-h-screen bg-slate-950 text-white">
      <Seo
        title="Leistungen: Softwareentwicklung & KI-Automatisierung"
        description="Individuelle Softwarelösungen, KI-gestützte Prozessautomatisierung und Bot-/API-Integrationen - Freelance-Entwicklung von Marcel Turlach."
        path="/services"
      />
      <section className="relative text-center py-20 px-6 bg-slate-900">
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

      <section className="py-16 px-6 max-w-6xl mx-auto">
        <h2 className="sr-only">Leistungen im Überblick</h2>
        <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-3">
          {services.map((service) => (
            <motion.div
              key={service.title}
              className="bg-slate-900 rounded-2xl shadow-lg p-6 flex flex-col items-center text-center"
              whileHover={{ scale: 1.05 }}
            >
              {service.icon}
              <h3 className="text-2xl font-bold mb-2">{service.title}</h3>
              <p className="text-gray-300">{service.description}</p>
            </motion.div>
          ))}
        </div>
      </section>

      <section className="py-20 px-6 text-center bg-slate-900">
        <motion.h2
          className="text-3xl font-bold mb-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
        >
          Bereit, dein Projekt zu starten?
        </motion.h2>
        <p className="text-gray-100 mb-8">
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
