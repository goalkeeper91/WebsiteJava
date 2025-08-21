import FollowMe from '../components/socials/FollowMe';

const About = () => {
  return (
    <section className="relative w-full min-h-screen bg-black text-white py-16 px-6">
      <div className="absolute inset-0 bg-gradient-to-br from-purple-700 via-pink-500 to-yellow-400 opacity-20 blur-3xl z-0" />

      <div className="relative z-10 max-w-5xl mx-auto space-y-16">

        <div className="text-center">
          <h1 className="text-5xl font-bold mb-4 text-white">Über Goalkeeper91</h1>
          <p className="text-lg text-gray-300">Willkommen in meiner Welt des Gamings, Streamings & Entertainments!</p>
        </div>

        <div className="flex flex-col md:flex-row items-center md:items-start gap-10">
          <img
            src="/images/goalkeeper_logo.png"
            alt="Profil"
            className="w-40 h-40 rounded-full object-cover border-4 border-goalyBlue"
          />
          <div>
            <h2 className="text-2xl font-semibold mb-2 text-goalyBlue">Wer bin ich?</h2>
            <p className="text-gray-300 leading-relaxed">
              Hey! Ich bin Goalkeeper91 – leidenschaftlicher Gamer, Entertainer und Livestreamer. Seit 2018
              streame ich (manchmal auch weniger) regelmäßig auf Twitch und teile meine Highlights auf YouTube & TikTok.
              Bei mir findest du Gameplay, coole Community-Events und jede Menge Spaß!
            </p>
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-6 text-goalyBlue text-center">Mein Setup</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 text-sm text-gray-300">
            <div className="bg-slate-800/40 p-4 rounded-xl shadow-lg">
              <h3 className="text-lg font-bold mb-2 text-white">🎮 Hardware</h3>
              <ul className="list-disc list-inside">
                <li>CPU: AMD Ryzen 7 5700X</li>
                <li>GPU: Gainward GeForce RTX 4070 Ti Super</li>
                <li>RAM: 32GB (2x 16384MB) G.Skill RipJaws V</li>
                <li>Motherboard: Asus ROG Strix X470-F Gaming</li>
                <li>Elgato Stream Deck</li>
              </ul>
            </div>
            <div className="bg-slate-800/40 p-4 rounded-xl shadow-lg">
              <h3 className="text-lg font-bold mb-2 text-white">📷 Streaming Gear</h3>
              <ul className="list-disc list-inside">
                <li>Cam: Logitech C922</li>
                <li>Mic: Shure MV7+</li>
                <li>Headset: Beyerdynamic MMX 100</li>
                <li>Audio: GoXLR Mini</li>
                <li>Beleuchtung: Key Light</li>
              </ul>
            </div>
          </div>
        </div>

        <div className="flex justify-center mb-10">
            <img
                src="/images/streaming setup_and_me.jpg"
                alt="Goalkeeper91 und sein Streaming-Setup"
                className="rounded-2xl shadow-xl w-full max-w-md object-cover"
                />
        </div>

        <FollowMe />
      </div>
    </section>
  );
};

export default About;

