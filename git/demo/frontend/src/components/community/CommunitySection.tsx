import { motion } from "framer-motion";
import { FaDiscord } from "react-icons/fa";

const CommunitySection = () => {
  return (
    <section className="w-screen bg-gradient-to-br from-goalyBlue to-goalyCyan py-16 px-6">
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
          className="mb-8 text-lg"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
        >
          Tritt unserem Discord bei, tausche dich mit anderen aus und sei immer up to date.
        </motion.p>

        <motion.a
          href="https://discord.gg/deinserverlink"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-3 px-6 py-3 bg-white text-goalyBlue font-semibold rounded-full shadow-lg hover:bg-gray-100 transition"
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
