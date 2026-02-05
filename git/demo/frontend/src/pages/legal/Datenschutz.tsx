import React from "react";

const Datenschutz: React.FC = () => {
  return (
    <section className="w-full min-h-screen bg-slate-950 flex items-center justify-center py-10 px-4">
      <div className="max-w-4xl bg-slate-900 text-white p-8 rounded-lg shadow-lg overflow-auto">

        <h1 className="text-3xl font-bold mb-8 text-goalyBlue">Datenschutzerklärung</h1>

        <h2 className="text-2xl font-semibold mb-4">1. Datenschutz auf einen Blick</h2>

        <h3 className="text-xl font-semibold mt-4 mb-2">Allgemeine Hinweise</h3>
        <p>
          Die folgenden Hinweise geben einen einfachen Überblick darüber, was mit Ihren personenbezogenen Daten passiert, wenn Sie diese Website besuchen. Personenbezogene Daten sind alle Daten, mit denen Sie persönlich identifiziert werden können. Ausführliche Informationen zum Thema Datenschutz entnehmen Sie der vollständigen Datenschutzerklärung.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Datenerfassung auf dieser Website</h3>
        <p className="mb-2">
          <strong>Wer ist verantwortlich?</strong>
          Die Datenverarbeitung erfolgt durch den Websitebetreiber. Kontaktdaten finden Sie im Abschnitt „Hinweis zur verantwortlichen Stelle“.
        </p>
        <p className="mb-2">
          <strong>Wie erfassen wir Ihre Daten?</strong>
          Ihre Daten werden einerseits durch Eingaben (z. B. Kontaktformular) erhoben, andererseits automatisch durch unsere IT-Systeme beim Besuch der Website (z. B. Browser, Betriebssystem, Uhrzeit).
        </p>
        <p className="mb-2">
          <strong>Wofür nutzen wir Ihre Daten?</strong>
          - Zur fehlerfreien Bereitstellung der Website
          - Zur Analyse des Nutzerverhaltens
          - Zur Abwicklung von Vertragsangeboten, Bestellungen oder sonstigen Anfragen
        </p>
        <p>
          <strong>Ihre Rechte:</strong> jederzeit Auskunft, Berichtigung, Löschung, Einschränkung, Widerruf von Einwilligungen, Beschwerderecht bei der zuständigen Aufsichtsbehörde.
        </p>

        <h2 className="text-2xl font-semibold mb-4">2. Hosting</h2>
        <p>
          Diese Website wird extern gehostet. Die personenbezogenen Daten werden auf den Servern des Hosters gespeichert, z. B. IP-Adressen, Kontaktdaten, Vertragsdaten, Websitezugriffe.
        </p>
        <p className="mb-2">
          Hosting erfolgt zur Vertragserfüllung und im Interesse einer sicheren, schnellen und effizienten Bereitstellung der Website.
        </p>
        <p className="mb-2">
          <strong>Hoster:</strong> Contabo GmbH, Aschauer Straße 32a, 81549 München, Deutschland
        </p>
        <p>
          Ein Vertrag zur Auftragsverarbeitung (AVV) stellt sicher, dass personenbezogene Daten nur nach Weisung und DSGVO-konform verarbeitet werden.
        </p>

        <h2 className="text-2xl font-semibold mb-4">3. Allgemeine Hinweise und Pflichtinformationen</h2>

        <h3 className="text-xl font-semibold mt-4 mb-2">Datenschutz</h3>
        <p>
          Wir behandeln Ihre Daten vertraulich und gemäß den gesetzlichen Datenschutzvorschriften. Bitte beachten Sie, dass die Datenübertragung im Internet Sicherheitslücken aufweisen kann.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Hinweis zur verantwortlichen Stelle</h3>
        <p>
          <strong>Goalkeeper91 c/o NextlevelNation</strong><br />
          Inhaber: Christian Steinbach<br />
          Stettener Weg 2, 89584 Ehingen (Donau)<br />
          Telefon: [Telefonnummer]<br />
          E-Mail: info@goalkeeper91.de
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Speicherdauer</h3>
        <p>
          Personenbezogene Daten verbleiben bei uns, bis der Zweck der Verarbeitung entfällt oder gesetzliche Vorgaben eine längere Speicherung erfordern.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Rechtsgrundlagen der Datenverarbeitung</h3>
        <p>
          Verarbeitung erfolgt auf Basis von Einwilligung, Vertragserfüllung, rechtlicher Verpflichtung oder berechtigtem Interesse gemäß Art. 6 DSGVO. Spezielle Datenkategorien erfolgen auf Grundlage von Art. 9 DSGVO.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Empfänger von Daten</h3>
        <p>
          Daten werden nur an externe Stellen weitergegeben, wenn dies gesetzlich erlaubt oder zur Vertragserfüllung erforderlich ist. Auftragsverarbeiter erhalten Daten nur nach AVV.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Widerruf & Widerspruch</h3>
        <p>
          Bereits erteilte Einwilligungen können jederzeit widerrufen werden. Sie haben das Recht auf Widerspruch gegen Verarbeitung zu besonderen Fällen oder Direktwerbung.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Beschwerderecht</h3>
        <p>
          Bei Verstößen gegen DSGVO können Sie sich bei der zuständigen Aufsichtsbehörde beschweren.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Recht auf Datenübertragbarkeit</h3>
        <p>
          Daten, die automatisiert verarbeitet werden, können Ihnen in maschinenlesbarem Format übermittelt werden.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Auskunft, Berichtigung, Löschung, Einschränkung</h3>
        <p>
          Sie haben das Recht auf Auskunft, Berichtigung, Löschung oder Einschränkung Ihrer personenbezogenen Daten.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">SSL- bzw. TLS-Verschlüsselung</h3>
        <p>
          Die Website nutzt SSL/TLS zur sicheren Übertragung vertraulicher Inhalte.
        </p>

        <h2 className="text-2xl font-semibold mb-4">4. Datenerfassung auf dieser Website</h2>

        <h3 className="text-xl font-semibold mt-4 mb-2">Cookies und Local Storage</h3>
        <p>
          Wir nutzen technisch notwendige Cookies und den Local Storage Ihres Browsers, um Ihren Login-Status sicherzustellen. Ohne diese Speicherung kann die Dashboard-Funktionalität nicht bereitgestellt werden.
        </p>

        <h3 className="text-xl font-semibold mt-4 mb-2">Anfragen per E-Mail, Telefon oder Fax</h3>
        <p>
          Kontaktanfragen werden zum Zwecke der Bearbeitung gespeichert und nur mit Einwilligung weitergegeben.
        </p>

        <h2 className="text-2xl font-semibold mb-4">5. Plugins und Tools</h2>

        <h3 className="text-xl font-semibold mt-4 mb-2">YouTube mit erweitertem Datenschutz</h3>
        <p>
          Videos von YouTube werden im erweiterten Datenschutzmodus eingebunden. Betreiber: Google Ireland Limited, Dublin, Irland.
        </p>

        <h2 className="text-2xl font-semibold mb-4">6. Twitch OAuth und API-Nutzung</h2>
        <p className="mb-2">
          Wir bieten die Anmeldung über den Dienst „Twitch“ an (Twitch Interactive, Inc., USA).
        </p>
        <p className="mb-2">
          <strong>Datenkategorien & Speicherung:</strong> Beim Login speichern wir dauerhaft in unserer Datenbank: Twitch-User-ID, Anzeigename, Profilbild-URL und einen verschlüsselten Access-Token.
        </p>
        <p className="mb-2">
          <strong>Zweck:</strong> Synchronisation von Kanal-Events (Subs, Bits, Follows) für den Timer sowie Speicherung Ihrer Dashboard-Konfiguration.
        </p>
        <p>
          <strong>Rechtsgrundlage:</strong> Einwilligung (Art. 6 Abs. 1 lit. a DSGVO) und Vertragserfüllung (Art. 6 Abs. 1 lit. b DSGVO). Zur Löschung Ihrer Daten senden Sie bitte eine E-Mail an info@goalkeeper91.de.
        </p>

      <p className="text-sm mt-10 text-gray-400">
        Quelle: <a href="https://www.e-recht24.de" className="underline text-goalyBlue" target="_blank" rel="noopener noreferrer">eRecht24</a> & Ergänzungen für Twitch API.
      </p>
     </div>
    </section>
  );
};

export default Datenschutz;