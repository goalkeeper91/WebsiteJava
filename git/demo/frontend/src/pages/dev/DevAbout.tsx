import Seo from "../../components/Seo";

type Job = {
  period: string;
  title: string;
  company?: string;
  bullets: string[];
};

// CV content sourced directly from the user's own Lebenslauf.pdf - not
// invented. Deliberately omits home address and phone number even though
// both are in the source PDF: neither belongs on a public marketing page
// (the Impressum page already covers the legally required business
// address, and the Contact page's form is the safer public entry point
// than a scraped-and-spammed raw email/phone).
const experience: Job[] = [
  {
    period: "07/2025 – heute",
    title: "Fullstack Developer (Eigenprojekt)",
    bullets: [
      "Entwicklung einer Personal-Brand- & Community-Plattform (goalkeeper91.de).",
      "Tech-Stack: Java (Spring Boot), React, PostgreSQL.",
      "Automatisierung: Entwicklung von Twitch- & Discord-Bots; laufendes Refactoring auf Python zur Performance-Optimierung.",
      "Deployment: Containerisierung mit Docker & Linux-Server-Administration.",
    ],
  },
  {
    period: "02/2024 – 07/2025",
    title: "Webentwickler",
    company: "Travello GmbH",
    bullets: [
      "Weiterentwicklung komplexer Webanwendungen mit PHP (Symfony) & React.",
      "Implementierung von Microservices (Event Driven Design).",
      "Optimierung von SQL-Abfragen und Docker-Umgebungen.",
    ],
  },
  {
    period: "02/2022 – 02/2024",
    title: "Fokusphase: Professionalisierung",
    bullets: [
      "Vertiefung moderner Software-Architekturen & Coding-Best-Practices.",
      "Qualifizierung bei Salo + Partner (Kaufmännische Grundlagen & EDV).",
    ],
  },
  {
    period: "12/2021 – 02/2022",
    title: "Webentwickler",
    company: "Hans John Versicherungsmakler GmbH",
    bullets: [
      "Backend-Entwicklung für Versicherungslösungen (PHP/Symfony).",
      "NoSQL-Datenmanagement mit MongoDB.",
    ],
  },
  {
    period: "02/2012 – 10/2019",
    title: "Rettungssanitäter",
    company: "Hamburg & Salzhausen",
    bullets: [
      "Medizinische Erstversorgung und Krisenmanagement.",
      "Kernkompetenzen, die direkt in die Softwareentwicklung einfließen: Stressresistenz, Teamarbeit, schnelle Problemlösung unter Druck.",
    ],
  },
];

const education: Job[] = [
  {
    period: "02/2026 – 05/2026",
    title: "Certified AI Developer (Vollzeit)",
    company: "velpTEC GmbH - Schwerpunkt Python & Prozessautomatisierung",
    bullets: [
      "Python Ecosystem: fundierte Entwicklung inkl. automatisierter Workflows.",
      "Prozessautomatisierung: Design & Implementierung KI-gestützter Automatisierungslösungen (Digital Automations).",
      "KI & Prompt Engineering: Anwendung von LLMs für datenintensive Aufgaben.",
      "Projektarbeit: skalierbare Lösungen & API-Integrationen.",
    ],
  },
  {
    period: "04/2020 – 02/2022",
    title: "Fachinformatiker Anwendungsentwicklung",
    company: "CBM Projektmanagement, Hamburg (Umschulung) - Abschluss: IHK-Zeugnis",
    bullets: [],
  },
];

const certificates = [
  "EXIN Agile Scrum Foundation",
  "Certified AI Developer (velpTEC)",
  "AI Engineer Zertifikat",
  "Digital Automations mit KI: Prozessoptimierung",
  "Prompt Engineering",
  "Python 3 Zertifikat",
  "GitHub/GitLab Zertifikat",
];

const techStack = {
  "Programmierung": ["Java", "PHP", "Python", "JavaScript (React)", "C#", "SQL"],
  "Frameworks": ["Spring Boot", "Symfony", "Laravel", ".NET Core", "Redis"],
  "Datenbanken": ["PostgreSQL", "MongoDB", "MySQL", "Oracle"],
  "DevOps & Tools": ["Docker", "Git", "n8n (Workflow-Automation)"],
};

const DevAbout = () => {
  return (
    <section className="relative w-full min-h-screen bg-slate-950 text-white py-16">
      <Seo
        title="Über mich - Marcel Turlach"
        description="Marcel Turlach: Fullstack-Entwickler mit Fokus auf Backend & KI-Automatisierung. Werdegang, Tech-Stack und Zertifikate."
        path="/about"
      />
      <div className="relative z-10 max-w-4xl mx-auto space-y-16 px-4">

        <div className="text-center">
          <img
            src="/images/marcel-turlach-headshot.jpg"
            alt="Marcel Turlach"
            width={160}
            height={160}
            loading="lazy"
            className="w-40 h-40 rounded-full object-cover border-4 border-goalyBlue mx-auto mb-6"
          />
          <h1 className="text-5xl font-bold mb-4 text-white">Marcel Turlach</h1>
          <p className="text-xl text-goalyBlue font-semibold mb-4">Fullstack Developer | Backend & AI Automation</p>
          <p className="text-lg text-gray-300 max-w-2xl mx-auto">
            Lösungsorientierter Fullstack-Entwickler mit fundierter Erfahrung in moderner Web-Architektur (PHP) und
            aktueller Spezialisierung auf Python sowie KI-gestützte Prozessautomatisierung. Durch meinen Hintergrund
            im Rettungsdienst bringe ich hohe Belastbarkeit, Teamgeist und ein strukturiertes Vorgehen in kritischen
            Situationen in die Softwareentwicklung ein.
          </p>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-6 text-goalyBlue">Berufserfahrung</h2>
          <div className="space-y-6">
            {experience.map((job) => (
              <div key={job.period + job.title} className="bg-slate-800/40 p-5 rounded-xl shadow-lg">
                <div className="flex flex-wrap items-baseline justify-between gap-2 mb-2">
                  <h3 className="text-lg font-bold text-white">
                    {job.title}{job.company && <span className="text-gray-400 font-normal"> · {job.company}</span>}
                  </h3>
                  <span className="text-sm text-gray-400 whitespace-nowrap">{job.period}</span>
                </div>
                <ul className="list-disc list-inside text-sm text-gray-300 space-y-1">
                  {job.bullets.map((b) => <li key={b}>{b}</li>)}
                </ul>
              </div>
            ))}
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-6 text-goalyBlue">Ausbildung & Weiterbildung</h2>
          <div className="space-y-6">
            {education.map((edu) => (
              <div key={edu.period + edu.title} className="bg-slate-800/40 p-5 rounded-xl shadow-lg">
                <div className="flex flex-wrap items-baseline justify-between gap-2 mb-2">
                  <h3 className="text-lg font-bold text-white">{edu.title}</h3>
                  <span className="text-sm text-gray-400 whitespace-nowrap">{edu.period}</span>
                </div>
                {edu.company && <p className="text-sm text-gray-400 mb-2">{edu.company}</p>}
                {edu.bullets.length > 0 && (
                  <ul className="list-disc list-inside text-sm text-gray-300 space-y-1">
                    {edu.bullets.map((b) => <li key={b}>{b}</li>)}
                  </ul>
                )}
              </div>
            ))}
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-6 text-goalyBlue text-center">Tech-Stack</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 text-sm text-gray-300">
            {Object.entries(techStack).map(([category, items]) => (
              <div key={category} className="bg-slate-800/40 p-4 rounded-xl shadow-lg">
                <h3 className="text-lg font-bold mb-2 text-white">{category}</h3>
                <ul className="list-disc list-inside">
                  {items.map((item) => <li key={item}>{item}</li>)}
                </ul>
              </div>
            ))}
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-6 text-goalyBlue text-center">Zertifikate</h2>
          <div className="flex flex-wrap justify-center gap-3">
            {certificates.map((cert) => (
              <span key={cert} className="bg-slate-800/60 border border-goalyBlue/40 text-sm text-gray-200 px-4 py-2 rounded-full">
                {cert}
              </span>
            ))}
          </div>
        </div>

        <div className="text-center text-sm text-gray-400">
          Sprachen: Deutsch (Muttersprache), Englisch (Technical English)
        </div>
      </div>
    </section>
  );
};

export default DevAbout;
