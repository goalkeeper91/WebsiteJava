import Seo from '../components/Seo';
import FollowMe from '../components/socials/FollowMe';

const About = () => {
  return (
    <section className="relative w-full min-h-screen bg-slate-950 text-white py-16">
      <Seo
        title="Über mich"
        description="Marcel alias Goalkeeper91: Full-Stack Softwareentwickler und Twitch-Streamer. Tech-Stack, Hintergrund und wie Softwareentwicklung und Live-Entertainment bei mir zusammenkommen."
        path="/about"
      />
      <div className="relative z-10 max-w-5xl mx-auto space-y-16">

        {/* Hero */}
        <div className="text-center">
          <h1 className="text-5xl font-bold mb-4 text-white">Über Goalkeeper91</h1>
          <p className="text-lg text-gray-300">
            Willkommen in meiner Welt aus Gaming, Livestreaming und Softwareentwicklung!
            Ich kombiniere kreative Projekte mit technischem Know-how – für einzigartige Lösungen.
          </p>
        </div>

        {/* Wer bin ich? */}
        <div className="flex flex-col md:flex-row items-center md:items-start gap-10">
          <img
            src="/images/goalkeeper_logo.png"
            alt="Goalkeeper91 Logo"
            width={160}
            height={160}
            loading="lazy"
            className="w-40 h-40 rounded-full object-cover border-4 border-goalyBlue"
          />
          <div>
            <h2 className="text-2xl font-semibold mb-2 text-goalyBlue">Wer bin ich?</h2>
            <p className="text-gray-300 leading-relaxed">
              Hey! Ich bin Goalkeeper91 – leidenschaftlicher Gamer, Livestreamer und Softwareentwickler.
              Seit 2018 streame ich regelmäßig auf Twitch, entwickle Tools für meine Community und realisiere individuelle Softwareprojekte.
              Dabei kombiniere ich Kreativität und Technik, um sowohl interaktive Erlebnisse als auch effiziente Lösungen für echte Probleme zu schaffen.
            </p>
          </div>
        </div>

        {/* Setup / Tech Stack */}
        <div>
          <h2 className="text-2xl font-semibold mb-6 text-goalyBlue text-center">Mein Setup & Tech Stack</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 text-sm text-gray-300">
            <div className="bg-slate-800/40 p-4 rounded-xl shadow-lg">
              <h3 className="text-lg font-bold mb-2 text-white">🎮 Hardware & Streaming</h3>
              <ul className="list-disc list-inside">
                <li>CPU: AMD Ryzen 7 5700X</li>
                <li>GPU: Gainward GeForce RTX 4070 Ti Super</li>
                <li>RAM: 32GB G.Skill RipJaws V</li>
                <li>Cam: Logitech C922</li>
                <li>Mic: Shure MV7+</li>
              </ul>
            </div>
            <div className="bg-slate-800/40 p-4 rounded-xl shadow-lg">
              <h3 className="text-lg font-bold mb-2 text-white">💻 Entwicklung & Tools</h3>
              <ul className="list-disc list-inside">
                <li>Frontend: React, Flutter</li>
                <li>Backend: Spring Boot, Laravel, GO, Python</li>
                <li>Datenbanken: MySQL, PostgreSQL</li>
                <li>DevOps & Automatisierung: Docker, CI/CD Pipelines</li>
                <li>Tools: VSCode, GitHub, Figma, N8N</li>
              </ul>
            </div>
          </div>
          <div className="flex justify-center mt-10">
            <img
              src="/images/streaming setup_and_me.jpg"
              alt="Goalkeeper91 und sein Streaming-Setup"
              loading="lazy"
              width={1000}
              height={1333}
              className="rounded-2xl shadow-xl w-full max-w-md object-cover"
            />
          </div>
        </div>

        {/* Socials */}
        <div className="flex justify-center mt-10">
          <FollowMe />
        </div>
      </div>
    </section>
  );
};

export default About;
