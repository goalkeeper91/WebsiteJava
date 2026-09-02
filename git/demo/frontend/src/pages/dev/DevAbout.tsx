import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { FaLayerGroup, FaRobot, FaAmbulance, FaBolt } from "react-icons/fa";
import Seo from "../../components/Seo";

type ValueProp = {
  icon: React.ReactNode;
  title: string;
  text: string;
};

// Four real, grounded reasons to hire him - not resume bullets restated,
// the actual pitch. The paramedic angle is a genuine differentiator (see
// experience below), not a gimmick - leading with it here instead of
// burying it as the last chronological entry is the whole point of this
// rewrite: a client skims this section, not a CV.
const valueProps: ValueProp[] = [
  {
    icon: <FaLayerGroup size={28} />,
    title: "Breiter, pragmatischer Stack",
    text: "Java, Python, PHP, React, C# - ich wähle das Werkzeug, das zum Problem passt, nicht umgekehrt.",
  },
  {
    icon: <FaRobot size={28} />,
    title: "KI-Automatisierung als Spezialisierung",
    text: "Certified AI Developer (velpTEC), Fokus auf Prozessautomatisierung mit LLMs - nicht nur Buzzword, sondern aktuelle Vollzeit-Weiterbildung.",
  },
  {
    icon: <FaAmbulance size={28} />,
    title: "Ruhe unter Druck",
    text: "7 Jahre Rettungssanitäter vor der Softwareentwicklung. Krisenmanagement, Priorisieren, nicht in Panik geraten - das bringe ich in jedes Projekt mit.",
  },
  {
    icon: <FaBolt size={28} />,
    title: "Direkter Draht, echte Verantwortung",
    text: "Kein Agentur-Overhead zwischen dir und der Umsetzung - von Prototyp bis Produktivbetrieb aus einer Hand.",
  },
];

type Job = { period: string; title: string; company?: string };

// Condensed on purpose - the full chronological detail (bullet-level)
// lives in the value props / portfolio instead. This is a scan-in-five-
// seconds timeline, not the CV itself.
const experience: Job[] = [
  { period: "07/2025 – heute", title: "Fullstack Developer", company: "Eigenprojekt (goalkeeper91.de)" },
  { period: "02/2024 – 07/2025", title: "Webentwickler", company: "Travello GmbH" },
  { period: "02/2022 – 02/2024", title: "Professionalisierung", company: "Vertiefung Software-Architektur" },
  { period: "12/2021 – 02/2022", title: "Webentwickler", company: "Hans John Versicherungsmakler GmbH" },
  { period: "02/2012 – 10/2019", title: "Rettungssanitäter", company: "Hamburg & Salzhausen" },
];

const education: Job[] = [
  { period: "02/2026 – 05/2026", title: "Certified AI Developer", company: "velpTEC GmbH (Vollzeit)" },
  { period: "04/2020 – 02/2022", title: "Fachinformatiker Anwendungsentwicklung", company: "IHK-Zeugnis" },
];

const techStack = [
  "Java", "Python", "PHP", "JavaScript / React", "C#", "SQL",
  "Spring Boot", "Symfony", "Laravel", ".NET Core", "Redis",
  "PostgreSQL", "MongoDB", "MySQL", "Oracle",
  "Docker", "Git", "n8n",
];

const certificates = [
  "Certified AI Developer (velpTEC)",
  "AI Engineer Zertifikat",
  "Digital Automations mit KI",
  "Prompt Engineering",
  "Python 3 Zertifikat",
  "EXIN Agile Scrum Foundation",
  "GitHub/GitLab Zertifikat",
];

const fadeUp = {
  initial: { opacity: 0, y: 20 },
  whileInView: { opacity: 1, y: 0 },
  viewport: { once: true, margin: "-50px" },
};

const DevAbout = () => {
  return (
    <div className="relative w-full bg-slate-950 text-white overflow-hidden">
      <Seo
        title="Über mich - Marcel Turlach"
        description="Marcel Turlach: Fullstack-Entwickler mit Fokus auf Backend & KI-Automatisierung. Warum mit mir arbeiten, Werdegang, Tech-Stack und Zertifikate."
        path="/about"
      />

      {/* Decorative background glow, matching the site's brand accent -
          purely so this page reads as a designed landing section instead
          of a plain document, same spirit as Hero.tsx's radial backdrop. */}
      <div className="pointer-events-none absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[900px] rounded-full bg-goalyBlue/10 blur-3xl" />

      {/* ============ HERO ============ */}
      <section className="relative z-10 max-w-5xl mx-auto px-4 pt-20 pb-16 flex flex-col md:flex-row items-center gap-12">
        <motion.div
          className="relative shrink-0"
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.6 }}
        >
          <div className="absolute inset-0 rounded-full bg-gradient-to-tr from-goalyBlue to-goalyCyan blur-2xl opacity-40" />
          <img
            src="/images/marcel-turlach-headshot.jpg"
            alt="Marcel Turlach"
            width={240}
            height={240}
            loading="lazy"
            className="relative w-52 h-52 md:w-60 md:h-60 rounded-full object-cover border-4 border-goalyBlue shadow-2xl"
          />
          <span className="absolute -bottom-2 left-1/2 -translate-x-1/2 whitespace-nowrap bg-slate-900 border border-goalyBlue/50 text-xs font-semibold px-3 py-1.5 rounded-full shadow-lg">
            🚑 Ex-Rettungssanitäter → Fullstack-Entwickler
          </span>
        </motion.div>

        <motion.div
          className="text-center md:text-left"
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6, delay: 0.15 }}
        >
          <h1 className="text-4xl md:text-5xl font-extrabold mb-3">Marcel Turlach</h1>
          <p className="text-xl text-goalyBlue font-semibold mb-5">
            Fullstack-Entwicklung, die einfach funktioniert
          </p>
          <p className="text-lg text-gray-300 max-w-xl mb-8">
            Ich baue Software für Menschen, die eine echte Lösung brauchen, keine Buzzword-Präsentation.
            Vielseitiger Stack, aktuelle Spezialisierung auf KI-gestützte Automatisierung, und eine
            Gelassenheit unter Druck, die ich nicht im Bootcamp gelernt habe.
          </p>
          <div className="flex flex-wrap justify-center md:justify-start gap-4">
            <Link
              to="/contact"
              className="px-6 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition"
            >
              Kostenloses Erstgespräch
            </Link>
            <Link
              to="/portfolio"
              className="px-6 py-3 border-2 border-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition"
            >
              Portfolio ansehen
            </Link>
          </div>
        </motion.div>
      </section>

      {/* ============ WARUM MIT MIR ARBEITEN ============ */}
      <section className="relative z-10 max-w-6xl mx-auto px-4 pb-20">
        <motion.h2 className="text-3xl font-bold text-center mb-10" {...fadeUp}>
          Warum mit mir arbeiten?
        </motion.h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          {valueProps.map((v, i) => (
            <motion.div
              key={v.title}
              className="bg-slate-900 rounded-2xl p-6 shadow-lg"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-50px" }}
              transition={{ delay: i * 0.08 }}
            >
              <div className="text-goalyBlue mb-3">{v.icon}</div>
              <h3 className="text-lg font-bold mb-2">{v.title}</h3>
              <p className="text-gray-300 text-sm">{v.text}</p>
            </motion.div>
          ))}
        </div>
      </section>

      {/* ============ WERDEGANG (kompakt) ============ */}
      <section className="relative z-10 bg-slate-900 py-16">
        <div className="max-w-4xl mx-auto px-4">
          <motion.h2 className="text-2xl font-bold text-center mb-10" {...fadeUp}>
            Werdegang
          </motion.h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
            <div>
              <h3 className="text-sm font-semibold uppercase tracking-wide text-goalyBlue mb-4">Berufserfahrung</h3>
              <div className="space-y-3">
                {experience.map((job) => (
                  <div key={job.period + job.title} className="flex justify-between gap-3 text-sm border-b border-slate-800 pb-2">
                    <div>
                      <span className="font-semibold text-white">{job.title}</span>
                      {job.company && <span className="text-gray-400"> · {job.company}</span>}
                    </div>
                    <span className="text-gray-500 whitespace-nowrap">{job.period}</span>
                  </div>
                ))}
              </div>
            </div>
            <div>
              <h3 className="text-sm font-semibold uppercase tracking-wide text-goalyBlue mb-4">Ausbildung & Weiterbildung</h3>
              <div className="space-y-3">
                {education.map((edu) => (
                  <div key={edu.period + edu.title} className="flex justify-between gap-3 text-sm border-b border-slate-800 pb-2">
                    <div>
                      <span className="font-semibold text-white">{edu.title}</span>
                      {edu.company && <span className="text-gray-400"> · {edu.company}</span>}
                    </div>
                    <span className="text-gray-500 whitespace-nowrap">{edu.period}</span>
                  </div>
                ))}
              </div>
              <p className="text-xs text-gray-500 mt-4">Sprachen: Deutsch (Muttersprache), Englisch</p>
            </div>
          </div>
        </div>
      </section>

      {/* ============ SKILLS & ZERTIFIZIERUNGEN ============ */}
      <section className="relative z-10 max-w-4xl mx-auto px-4 py-16">
        <motion.h2 className="text-2xl font-bold text-center mb-8" {...fadeUp}>
          Tech-Stack
        </motion.h2>
        <motion.div className="flex flex-wrap justify-center gap-2 mb-14" {...fadeUp}>
          {techStack.map((item) => (
            <span key={item} className="bg-slate-800 border border-goalyBlue/30 text-sm text-gray-200 px-3.5 py-1.5 rounded-full">
              {item}
            </span>
          ))}
        </motion.div>

        <motion.h2 className="text-2xl font-bold text-center mb-8" {...fadeUp}>
          Zertifikate
        </motion.h2>
        <motion.div className="flex flex-wrap justify-center gap-2" {...fadeUp}>
          {certificates.map((cert) => (
            <span key={cert} className="bg-slate-800/60 border border-goalyCyan/30 text-sm text-gray-200 px-3.5 py-1.5 rounded-full">
              {cert}
            </span>
          ))}
        </motion.div>
      </section>

      {/* ============ CTA ============ */}
      <section className="relative z-10 bg-slate-900 py-20 px-6 text-center">
        <motion.h2 className="text-3xl font-bold mb-6" {...fadeUp}>
          Klingt nach deinem Projekt?
        </motion.h2>
        <p className="text-gray-300 mb-8 max-w-xl mx-auto">
          Lass uns in einem kostenlosen Erstgespräch klären, ob und wie ich dir helfen kann.
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

export default DevAbout;
