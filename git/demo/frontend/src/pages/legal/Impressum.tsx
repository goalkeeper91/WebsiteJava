import React from "react";
import { FaEnvelope } from "react-icons/fa";

const Impressum: React.FC = () => {
  const emailUser = "info";
  const emailDomain = "goalkeeper91.de";

  const handleEmailClick = () => {
    window.location.href = `mailto:${emailUser}@${emailDomain}`;
  };

  return (
    <section className="relative w-full min-h-screen bg-black text-white py-16 px-6">
      {/* Hintergrund-Gradient für coolen Effekt */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-700 via-pink-500 to-yellow-400 opacity-20 blur-3xl z-0" />

      <div className="relative z-10 max-w-4xl mx-auto space-y-10">
        {/* Titel */}
        <div className="text-center">
          <h1 className="text-5xl font-bold mb-4 text-white">Impressum</h1>
          <p className="text-lg text-gray-300">
            Gesetzliche Anbieterkennzeichnung nach § 5 TMG
          </p>
        </div>

        {/* Inhalt in Karten-Style */}
        <div className="bg-slate-800/40 p-6 rounded-xl shadow-lg space-y-6">
          <div>
            <h2 className="text-2xl font-semibold mb-2 text-goalyBlue">Angaben zum Anbieter</h2>
            <p className="text-gray-300 leading-relaxed">
              Goalkeeper91 <br />
              c/o NextlevelNation <br />
              Inhaber: Christian Steinbach <br />
              Stettener Weg 2 <br />
              89584 Ehingen (Donau)
            </p>
          </div>

          <div>
            <h2 className="text-2xl font-semibold mb-2 text-goalyBlue">Verantwortlich für den Inhalt</h2>
            <p className="text-gray-300">Marcel Turlach</p>
          </div>

          <div>
            <h2 className="text-2xl font-semibold mb-2 text-goalyBlue">Kontakt</h2>
            <p
              onClick={handleEmailClick}
              className="flex items-center gap-2 cursor-pointer text-gray-300 hover:text-goalyBlue transition text-lg"
            >
              <FaEnvelope className="text-goalyBlue" />
              <span className="font-mono">info</span>
              <span className="text-goalyBlue">@</span>
              <span className="font-mono">goalkeeper91.de</span>
            </p>

          </div>
        </div>
      </div>
    </section>
  );
};

export default Impressum;
