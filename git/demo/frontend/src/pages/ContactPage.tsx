import ContactForm from "../components/ContactForm";

const ContactPage = () => {
  return (
    <section className="relative w-full min-h-screen py-16 px-6 bg-gray-900 text-white">
      <div className="absolute inset-0 bg-gradient-to-br from-purple-500 via-pink-500 to-yellow-500 opacity-20 blur-3xl z-0" />
      <div className="relative z-10 max-w-3xl mx-auto text-center">
        <h1 className="text-4xl font-bold mb-4">Kontakt aufnehmen</h1>
        <p className="text-lg text-gray-300 mb-12">
          Schreib mir eine Nachricht für Freelance-Projekte oder Kooperationen. Ich melde mich schnellstmöglich zurück!
        </p>

        <ContactForm />
      </div>
    </section>
  );
};

export default ContactPage;
