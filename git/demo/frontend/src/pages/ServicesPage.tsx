import { motion } from "framer-motion";
import { useState } from "react";
import { Link } from 'react-router-dom';
import { FaCode, FaCogs, FaStream, FaGithub } from "react-icons/fa";
import Seo from "../components/Seo";

type Service = {
  title: string;
  description: string;
  icon: React.ReactNode;
  image?: { src: string; alt: string } | null;
  images?: { src: string; alt: string }[];
};

const ServicesPage = () => {
  const [lightboxImage, setLightboxImage] = useState<string | null>(null);

  const services: Service[] = [
    {
      title: "Individuelle Softwarelösungen",
      description:
        "Von Prototyp bis fertiges Produkt – maßgeschneiderte Anwendungen für deine Anforderungen. Egal ob Web, Mobile oder Desktop.",
      icon: <FaCode size={40} className="text-indigo-400 mb-4" />,
      image: { src: "/images/faceit_reader.png", alt: "Individuelle Softwarelösung: Faceit-Reader-Tool" },
    },
    {
      title: "Twitch Bot & Automatisierung",
      description:
        "Ich entwickle Twitch-Bots, Discord-Integrationen und SaaS-Automatisierungen (Automod, Clip-Erstellung, Chat-Commands), die Prozesse automatisieren – für Twitch, Discord oder deine internen Abläufe.",
      icon: <FaCogs size={40} className="text-green-400 mb-4" />,
      image: { src: "/images/Twitch Bot Overview.png", alt: "Übersicht des Twitch-Bot-Dashboards" },
    },
    {
      title: "Streaming & Community Tech",
      description:
        "Als Streamer entwickle ich Tools & Overlays, die deine Community stärker binden und neue Features ins Stream-Erlebnis bringen.",
      icon: <FaStream size={40} className="text-pink-400 mb-4" />,
      images: [
        { src: "/images/discord-bot.png", alt: "Discord-Bot Overlay Screenshot" },
        { src: "/images/twitch-bot.png", alt: "Twitch-Bot Screenshot" },
      ],
    },
  ];

  return (
    <div className="relative w-full min-h-screen bg-slate-950 text-white">
      <Seo
        title="Services: Softwareentwicklung & Twitch Bot"
        description="Individuelle Softwarelösungen, Twitch-Bot-Automatisierung (Automod, Clip-Erstellung, Chat-Commands) und Streaming-Tools - Freelance-Entwicklung von Marcel alias Goalkeeper91."
        path="/services"
      />
      {/* Hero Section */}
      <section className="relative text-center py-20 px-6 bg-slate-900">
        <motion.h1
          className="text-5xl font-extrabold mb-4"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          Meine Services
        </motion.h1>
        <motion.p
          className="text-lg text-gray-300 max-w-3xl mx-auto mb-8"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          Ich bin kein „Alleskönner“ – und das ist gut so. Statt alles ein bisschen zu machen,
          konzentriere ich mich darauf, Probleme zu verstehen und die passende Lösung zu entwickeln.
        </motion.p>
      </section>

      {/* Services Grid */}
      <section className="py-16 px-6 max-w-6xl mx-auto">
        <h2 className="sr-only">Meine Leistungen im Überblick</h2>
        <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-3">
        {services.map((service: Service) => (
          <motion.div
            key={service.title}
            className="bg-slate-900 rounded-2xl shadow-lg p-6 flex flex-col items-center text-center"
            whileHover={{ scale: 1.05 }}
          >
            {service.icon}
            <h3 className="text-2xl font-bold mb-2">{service.title}</h3>
            <p className="text-gray-300 mb-4">{service.description}</p>

            {/* Platzhalter, wenn kein Bild */}
            {service.image === null && (
              <div className="w-full max-w-xl mx-auto rounded-lg border border-gray-700 h-40 flex items-center justify-center mb-3">
                <span className="text-gray-500">[ Screenshot / Projektbild ]</span>
              </div>
            )}

            {/* Falls ein einzelnes Bild */}
            {service.image && (
              <img
                src={service.image.src}
                alt={service.image.alt}
                loading="lazy"
                className="w-full max-w-xl mx-auto rounded-lg border border-gray-700 cursor-pointer hover:scale-105 transition-transform mb-3"
                onClick={() => setLightboxImage(service.image!.src)}
              />
            )}

            {/* Falls mehrere Bilder */}
            {service.images &&
              service.images.map((img) => (
                <img
                  key={img.src}
                  src={img.src}
                  alt={img.alt}
                  loading="lazy"
                  className="w-full max-w-xl mx-auto rounded-lg border border-gray-700 cursor-pointer hover:scale-105 transition-transform mb-3"
                  onClick={() => setLightboxImage(img.src)}
                />
              ))}
          </motion.div>
        ))}
        </div>
      </section>

      {/* Lightbox */}
      {lightboxImage && (
        <div
          className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
          onClick={() => setLightboxImage(null)}
        >
          <img
            src={lightboxImage}
            alt="Vollbild"
            className="max-h-full max-w-full rounded-lg shadow-lg"
          />
        </div>
      )}

      {/* GitHub / Referenzen */}
      <section className="py-16 px-6 text-center bg-slate-900">
        <motion.h2
          className="text-3xl font-bold mb-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
        >
          Open Source & Referenzen
        </motion.h2>
        <p className="text-gray-400 mb-8">
          Einen Teil meiner Projekte findest du direkt auf GitHub – darunter Flutter-Apps,
          React-Frontends, Spring Boot Backends & mehr.
        </p>
        <a
          href="https://github.com/goalkeeper91"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-3 px-6 py-3 bg-goalyBlue hover:bg-goalyCyan text-white font-semibold rounded-full shadow-md transition"
        >
          <FaGithub size={24} />
          Zu meinem GitHub
        </a>
      </section>

      {/* Call to Action */}
      <section className="py-20 px-6 text-center bg-slate-950">
        <motion.h2
          className="text-3xl font-bold mb-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
        >
          Bereit, dein Projekt zu starten?
        </motion.h2>
        <p className="text-gray-100 mb-8">
          Lass uns über deine Ideen sprechen – ob kleine Automatisierung oder komplexe Softwarelösung.
        </p>
        <Link
          to="/contact"
          className="px-8 py-3 bg-goalyBlue hover:bg-goalyCyan rounded-lg text-white font-semibold transition"
        >
          Kontakt aufnehmen
        </Link>
      </section>
    </div>
  );
};

export default ServicesPage;
