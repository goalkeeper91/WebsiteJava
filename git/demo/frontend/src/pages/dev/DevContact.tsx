import { motion } from "framer-motion";
import { FaCheckCircle } from "react-icons/fa";
import Seo from "../../components/Seo";
import ContactForm from "../../components/ContactForm";

const steps = [
  { step: "1", title: "Erstgespräch", text: "Kostenlos und unverbindlich - wir klären, worum es geht." },
  { step: "2", title: "Individuelles Angebot", text: "Kein Pauschalpreis - passend zu Umfang und Budget deines Projekts." },
  { step: "3", title: "Umsetzung", text: "Direkter Draht zu mir, kein Agentur-Overhead dazwischen." },
];

const DevContact = () => {
  return (
    <div className="relative w-full bg-slate-950 text-white overflow-hidden">
      <Seo
        title="Kontakt"
        description="Du suchst einen Entwickler für dein Softwareprojekt? Schreib mir eine Nachricht - kostenloses Erstgespräch, keine Pauschalpreise, individuelles Angebot."
        path="/contact"
      />

      {/* Same decorative glow as the rest of the storefront. */}
      <div className="pointer-events-none absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[900px] rounded-full bg-goalyBlue/10 blur-3xl" />

      <section className="relative z-10 max-w-5xl mx-auto px-4 py-20 grid grid-cols-1 md:grid-cols-2 gap-12 items-start">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
        >
          <h1 className="text-4xl font-bold mb-4">Lass uns sprechen</h1>
          <p className="text-lg text-gray-300 mb-10">
            Egal ob kleine Automatisierung oder komplexe Softwarelösung - schreib mir eine Nachricht,
            ich melde mich schnellstmöglich zurück.
          </p>

          <div className="space-y-6 mb-10">
            {steps.map((s) => (
              <div key={s.step} className="flex gap-4">
                <span className="shrink-0 w-9 h-9 rounded-full bg-goalyBlue/20 border border-goalyBlue text-goalyCyan font-bold flex items-center justify-center">
                  {s.step}
                </span>
                <div>
                  <h3 className="font-semibold">{s.title}</h3>
                  <p className="text-sm text-gray-400">{s.text}</p>
                </div>
              </div>
            ))}
          </div>

          <div className="flex items-center gap-2 text-sm text-gray-400">
            <FaCheckCircle className="text-goalyCyan" /> Antwort in der Regel innerhalb von 1-2 Werktagen
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6, delay: 0.15 }}
        >
          <ContactForm />
        </motion.div>
      </section>
    </div>
  );
};

export default DevContact;
