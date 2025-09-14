import { useState } from "react";

const ContactForm = () => {
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    phone: "",
    subject: "",
    message: "",
  });

  const [status, setStatus] = useState<null | "loading" | "success" | "error">(null);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [consentGiven, setConsentGiven] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setStatus("loading");
    setErrors({});

    // Frontend-Validierung für DSGVO
    if (!consentGiven) {
      setErrors({ consent: "Du musst der Speicherung deiner Daten zustimmen." });
      setStatus("error");
      return;
    }

    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...formData, consentGiven }),
      });

      if (res.status === 400) {
        const data = await res.json();
        setErrors(data);
        setStatus("error");
      } else if (res.ok) {
        setStatus("success");
        setFormData({ name: "", email: "", phone: "", subject: "", message: "" });
        setConsentGiven(false);
      } else {
        setStatus("error");
      }
    } catch (err) {
      setStatus("error");
    }
  };

  return (
    <div className="w-full py-10 px-4">
      <h1 className="text-3xl font-bold mb-6 text-center text-white">Kontakt</h1>

      <form onSubmit={handleSubmit} className="space-y-4 bg-gray-900/80 p-8 rounded-2xl shadow-2xl backdrop-blur-md">
        <input
          name="name"
          type="text"
          placeholder="Dein Name"
          value={formData.name}
          onChange={handleChange}
          className="w-full p-4 rounded-xl bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-goalyBlue transition"
        />
        {errors.name && <p className="text-red-500">{errors.name}</p>}

        <input
          name="email"
          type="email"
          placeholder="Deine Email"
          value={formData.email}
          onChange={handleChange}
          className="w-full p-4 rounded-xl bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-goalyBlue transition"
        />
        {errors.email && <p className="text-red-500">{errors.email}</p>}

        <input
          name="phone"
          type="text"
          placeholder="Telefon (optional)"
          value={formData.phone}
          onChange={handleChange}
          className="w-full p-4 rounded-xl bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-goalyBlue transition"
        />

        <input
          name="subject"
          type="text"
          placeholder="Betreff"
          value={formData.subject}
          onChange={handleChange}
          className="w-full p-4 rounded-xl bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-goalyBlue transition"
        />
        {errors.subject && <p className="text-red-500">{errors.subject}</p>}

        <textarea
          name="message"
          placeholder="Nachricht"
          value={formData.message}
          onChange={handleChange}
          rows={5}
          className="w-full p-4 rounded-xl bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-goalyBlue transition"
        />
        {errors.message && <p className="text-red-500">{errors.message}</p>}

        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={consentGiven}
            onChange={() => setConsentGiven(prev => !prev)}
          />
          <span>Ich stimme der Speicherung meiner Daten zu</span>
        </label>
        {errors.consent && <p className="text-red-500">{errors.consent}</p>}

        <button
          type="submit"
          disabled={status === "loading"}
          className="w-full py-4 bg-goalyBlue hover:bg-goalyCyan text-white rounded-2xl font-semibold transition-all shadow-lg hover:shadow-xl"
        >
          {status === "loading" ? "Sende..." : "Absenden"}
        </button>

        {status === "success" && <p className="text-green-400 mt-2 text-center">✅ Danke! Deine Anfrage wurde gesendet.</p>}
        {status === "error" && !errors.consent && <p className="text-red-500 mt-2 text-center">❌ Es ist ein Fehler aufgetreten.</p>}
      </form>
    </div>
  );
};

export default ContactForm;
