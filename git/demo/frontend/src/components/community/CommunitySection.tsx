import { useEffect, useState } from 'react';
import { motion } from "framer-motion";
import { FaDiscord } from "react-icons/fa";

const CommunitySection = () => {
    const [memberCount, setMemberCount] = useState<number | null>(null);

    useEffect(() => {
        fetch("https://discord.com/api/guilds/624607110860374016/widget.json")
        .then((res) => res.json())
        .then((data) => setMemberCount(data.presence_count))
        .catch(console.error);
    }, []);

  return (
    <section className="relative z-10 w-full bg-slate-900/80 py-16 px-6">
      <div className="max-w-4xl mx-auto text-center text-white">
        <motion.h2
          className="text-4xl font-bold mb-4"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          Werde Teil der Community
        </motion.h2>

        <motion.p
          className="mb-8 text-lg text-gray-300"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
        >
          <strong>{memberCount ?? '...'}</strong> Mitglieder sind dir nicht genug?<br />
          Tritt unserem Discord bei, tausch dich aus & bleib immer auf dem Laufenden.
        </motion.p>

        <motion.a
          href="https://discord.gg/XE8sW56"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-3 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-full shadow-md transition"
          whileHover={{ scale: 1.05 }}
        >
          <FaDiscord size={24} />
          Jetzt beitreten
        </motion.a>
      </div>
    </section>
  );
};

export default CommunitySection;
