import React, { useEffect, useRef } from "react";
import Seo from "../../components/Seo";

const COOKIEBOT_ID = "1273f5fe-7c66-466a-b0bc-7a15c4b60657";

const Cookies: React.FC = () => {
  const declarationRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = declarationRef.current;
    if (!container) return;

    // Cookiebot's auto-generated, live-scanned cookie list - a plain <script>
    // tag in JSX never executes, so it's appended imperatively here instead.
    const script = document.createElement("script");
    script.id = "CookieDeclaration";
    script.src = `https://consent.cookiebot.com/${COOKIEBOT_ID}/cd.js`;
    script.type = "text/javascript";
    script.async = true;
    container.appendChild(script);

    return () => {
      while (container.firstChild) {
        container.removeChild(container.firstChild);
      }
    };
  }, []);

  return (
    <section className="w-full min-h-screen bg-slate-950 flex items-center justify-center py-10 px-4">
      <Seo
        title="Cookie-Richtlinie"
        description="Welche Cookies Goalkeeper91 einsetzt und wie du deine Einwilligung über Cookiebot verwaltest."
        path="/legal/cookies"
      />
      <div className="max-w-4xl bg-slate-900 text-white p-8 rounded-lg shadow-lg overflow-auto">
        <h1 className="text-3xl font-bold mb-6">Cookie-Richtlinie</h1>

        <p className="mb-4">
          Wir nutzen Cookiebot, um Ihnen die Kontrolle über nicht technisch notwendige Cookies zu geben. Beim
          ersten Besuch dieser Website erscheint ein Consent-Banner, in dem Sie auswählen können, welche
          Cookie-Kategorien Sie zulassen möchten. Ihre Auswahl können Sie jederzeit über den Cookiebot-Banner
          erneut anpassen.
        </p>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Kategorien</h2>
        <ul className="list-disc list-inside mb-4 space-y-1">
          <li>
            <strong>Notwendig:</strong> technisch erforderliche Cookies, u. a. zur Aufrechterhaltung Ihres
            Login-Status im Dashboard und zur Speicherung Ihrer Cookie-Einwilligung selbst. Ohne diese Cookies
            funktioniert die Plattform nicht.
          </li>
          <li>
            <strong>Präferenzen:</strong> speichern Einstellungen wie z. B. Ihr bevorzugtes Layout im Dashboard.
          </li>
          <li>
            <strong>Statistik:</strong> helfen uns zu verstehen, wie die Website genutzt wird, damit wir sie
            verbessern können.
          </li>
          <li>
            <strong>Marketing:</strong> werden nur genutzt, wenn eingebettete Inhalte (z. B. YouTube-Videos) dies
            erfordern, und nur mit Ihrer Einwilligung geladen.
          </li>
        </ul>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Aktuelle Cookie-Übersicht</h2>
        <p className="mb-4 text-gray-400 text-sm">
          Die folgende Liste wird automatisch von Cookiebot anhand eines regelmäßigen Scans dieser Website erstellt
          und aktuell gehalten.
        </p>
        <div ref={declarationRef} />

        <p className="text-sm mt-10 text-gray-400">
          Weitere Informationen zur Datenverarbeitung finden Sie in unserer{" "}
          <a href="/legal/datenschutz" className="underline text-goalyBlue">Datenschutzerklärung</a>.
        </p>
      </div>
    </section>
  );
};

export default Cookies;
