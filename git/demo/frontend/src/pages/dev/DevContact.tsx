import Seo from "../../components/Seo";
import ContactForm from "../../components/ContactForm";

const DevContact = () => {
  return (
    <section className="relative w-full min-h-screen py-16 px-6 bg-slate-950 text-white">
      <Seo
        title="Kontakt"
        description="Du suchst einen Entwickler für dein Softwareprojekt? Schreib mir eine Nachricht - kostenloses Erstgespräch, keine Pauschalpreise, individuelles Angebot."
        path="/contact"
      />
      <div className="relative z-10 max-w-3xl mx-auto text-center">
        <h1 className="text-4xl font-bold mb-4">Kontakt aufnehmen</h1>
        <p className="text-lg text-gray-300 mb-12">
          Schreib mir eine Nachricht für dein Projekt - egal ob kleine Automatisierung oder komplexe
          Softwarelösung. Kostenloses Erstgespräch, ich melde mich schnellstmöglich zurück.
        </p>

        <ContactForm />
      </div>
    </section>
  );
};

export default DevContact;
