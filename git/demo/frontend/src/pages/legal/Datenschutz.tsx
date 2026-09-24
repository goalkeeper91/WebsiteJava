import React from "react";
import Seo from "../../components/Seo";

const Datenschutz: React.FC = () => {
  return (
    <section className="w-full min-h-screen bg-slate-950 flex items-center justify-center py-10 px-4">
      <Seo
        title="Datenschutzerklärung"
        description="Datenschutzerklärung von Goalkeeper91: welche Daten bei Twitch-/Discord-Login, der optionalen TikTok-Verknüpfung, Nutzung des Dashboards und der Bezahlung über Paddle verarbeitet werden."
        path="/legal/datenschutz"
      />
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

        <h3 className="text-xl font-semibold mt-4 mb-2">Übermittlung in Drittländer</h3>
        <p>
          Einige der unten genannten Dienstleister (u. a. Twitch, Discord, Paddle, TikTok und Composio) haben ihren
          Sitz oder betreiben Server ganz oder teilweise außerhalb der EU bzw. des EWR, insbesondere in den USA und
          im Vereinigten Königreich. Eine Übermittlung personenbezogener Daten dorthin erfolgt nur, wenn hierfür
          eine Grundlage nach Art. 44 ff. DSGVO besteht, etwa ein Angemessenheitsbeschluss der EU-Kommission (z. B.
          das EU-US Data Privacy Framework) oder EU-Standardvertragsklauseln des jeweiligen Anbieters, oder wenn Sie
          ausdrücklich eingewilligt haben (Art. 49 Abs. 1 lit. a DSGVO), etwa indem Sie freiwillig ein Konto eines
          solchen Anbieters mit uns verknüpfen.
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
        <p className="mt-2">
          Wenn Sie auf der Seite „/tiktok“ freiwillig ein TikTok-Konto verbinden (siehe Abschnitt 9), setzen wir
          zusätzlich ein technisch notwendiges Cookie mit dem Namen <code>tiktok_demo_uid</code>. Es enthält eine
          zufällig erzeugte Kennung (keine Angaben zu Ihrer Person), ist für Skripte nicht lesbar (HttpOnly), läuft
          nach 24 Stunden ab und dient allein dazu, Ihren Browser der von Ihnen gestarteten Verbindung zuzuordnen.
          Rechtsgrundlage: § 25 Abs. 2 Nr. 2 TDDDG sowie Art. 6 Abs. 1 lit. f DSGVO.
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

        <h2 className="text-2xl font-semibold mb-4">6. Zahlungsabwicklung über Paddle</h2>
        <p className="mb-2">
          Für kostenpflichtige Tarife nutzen wir den Zahlungsdienstleister Paddle.com Market Limited, 30 Old
          Bailey, London, EC4M 7AU, Vereinigtes Königreich ("Paddle"). Paddle agiert als Merchant of Record, tritt
          also rechtlich selbst als Verkäufer auf und übernimmt in diesem Zusammenhang eigenständige
          datenschutzrechtliche Verantwortung (z. B. für Rechnungsstellung, Umsatzsteuer und Zahlungsabwicklung).
        </p>
        <p className="mb-2">
          <strong>Datenkategorien:</strong> Name, E-Mail-Adresse, Rechnungsanschrift (falls angegeben) sowie
          Zahlungs- und Transaktionsdaten. Ihre Zahlungsdaten (z. B. Kartendaten) werden ausschließlich von Paddle
          verarbeitet und laufen nicht über unsere eigenen Server.
        </p>
        <p className="mb-2">
          <strong>Zweck:</strong> Abwicklung des Kaufs kostenpflichtiger Tarife, Verwaltung von Abonnements
          (Verlängerung, Kündigung, Tarifwechsel) über das Paddle-Kundenportal.
        </p>
        <p>
          <strong>Rechtsgrundlage:</strong> Vertragserfüllung (Art. 6 Abs. 1 lit. b DSGVO). Weitere Informationen
          finden Sie in Paddles eigener{" "}
          <a href="https://www.paddle.com/legal/privacy" className="underline text-goalyBlue" target="_blank" rel="noopener noreferrer">
            Datenschutzerklärung
          </a>.
        </p>

        <h2 className="text-2xl font-semibold mb-4 mt-6">7. Twitch OAuth und API-Nutzung</h2>
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

        <h2 className="text-2xl font-semibold mb-4 mt-6">8. Discord OAuth und Bot-Integration</h2>
        <p className="mb-2">
          Optional können Sie Ihren Discord-Server mit unserem Bot verbinden (Discord Inc., 444 De Haro Street,
          San Francisco, USA).
        </p>
        <p className="mb-2">
          <strong>Datenkategorien & Speicherung:</strong> Bei der Verbindung Ihres Discord-Accounts speichern wir
          dauerhaft in unserer Datenbank Ihre Discord-User-ID, Ihren Discord-Nutzernamen sowie einen verschlüsselten
          Zugriffs- und Refresh-Token. Fügen Sie den Bot einem Discord-Server hinzu, speichern wir zusätzlich die
          Server-ID (Guild-ID), den Servernamen, das Server-Icon, die Mitgliederzahl sowie die Discord-ID des
          Server-Inhabers zur Zuordnung der Berechtigung, den Server im Dashboard zu verwalten.
        </p>
        <p className="mb-2">
          <strong>Zweck:</strong> Bereitstellung der Bot-Funktionen auf Ihrem Discord-Server (z. B. Benachrichtigungen,
          Join-to-Create-Sprachkanäle) sowie Verwaltung dieser Einstellungen über Ihr Dashboard.
        </p>
        <p>
          <strong>Rechtsgrundlage:</strong> Einwilligung (Art. 6 Abs. 1 lit. a DSGVO) und Vertragserfüllung (Art. 6
          Abs. 1 lit. b DSGVO). Zur Löschung Ihrer Daten senden Sie bitte eine E-Mail an info@goalkeeper91.de.
        </p>

        <h2 className="text-2xl font-semibold mb-4 mt-6">9. TikTok-Verknüpfung (Login Kit / Content Posting API) über Composio</h2>
        <p className="mb-2">
          Auf der Seite „/tiktok“ können Sie optional ein TikTok-Konto mit uns verbinden. Die Verbindung erfolgt
          ausschließlich auf Ihre Veranlassung über das „Login Kit“ von TikTok (OAuth 2.0); Ihr TikTok-Passwort
          erhalten wir zu keinem Zeitpunkt. Für Nutzer im EWR ist TikTok Technology Limited (Irland) der Anbieter von
          TikTok, deren Umgang mit Ihren Daten in der{" "}
          <a href="https://www.tiktok.com/legal/page/eea/privacy-policy/en" className="underline text-goalyBlue" target="_blank" rel="noopener noreferrer">
            Datenschutzerklärung von TikTok
          </a>{" "}
          beschrieben ist.
        </p>
        <p className="mb-2">
          <strong>Angeforderte Berechtigungen:</strong> <code>user.info.basic</code> (Anzeigename und Profilbild
          Ihres TikTok-Kontos) sowie <code>video.upload</code> (Übermittlung eines Videos als Entwurf in den
          TikTok-Posteingang des verbundenen Kontos). Über diese Berechtigung wird nichts veröffentlicht; ein Entwurf
          wird erst durch Sie selbst in der TikTok-App bearbeitet und gepostet.
        </p>
        <p className="mb-2">
          <strong>Datenkategorien & Speicherung:</strong> Anzeigename und Profilbild-URL werden beim Aufruf der Seite
          von TikTok abgefragt und nur zur Anzeige verwendet; wir speichern sie nicht in unserer eigenen Datenbank.
          Ebenso speichern wir keine Inhalte Ihres TikTok-Kontos. Die Autorisierung selbst (verschlüsselte
          Zugriffs- und Aktualisierungs-Token, die TikTok-Kennung „open_id“, Verbindungsstatus) wird bei unserem
          technischen Dienstleister Composio (composio.dev) gespeichert, bis Sie die Verbindung trennen. Bei einem
          Video-Upload werden das Video und die zugehörigen Angaben (z. B. Beschreibungstext) an TikTok übermittelt.
          Zusätzlich fallen kurzzeitig technische Server-Zugriffsdaten (u. a. IP-Adresse, Zeitpunkt) beim Hoster an
          (siehe Abschnitt 2) sowie das in Abschnitt 4 beschriebene Cookie.
        </p>
        <p className="mb-2">
          <strong>Composio als Dienstleister:</strong> Composio vermittelt die Anmeldung bei TikTok und die
          API-Aufrufe in unserem Auftrag; der Zugriffsschlüssel für Composio liegt ausschließlich auf unserem Server
          und nie in Ihrem Browser. Hierzu gelten die Auftragsverarbeitungsbedingungen von Composio. Weitere
          Informationen finden Sie in der{" "}
          <a href="https://composio.dev/privacy" className="underline text-goalyBlue" target="_blank" rel="noopener noreferrer">
            Datenschutzerklärung von Composio
          </a>
          . Zur möglichen Übermittlung in Drittländer siehe Abschnitt 3.
        </p>
        <p className="mb-2">
          <strong>Zweck:</strong> Anzeige des verbundenen Kontos und Bereitstellung der Möglichkeit, Videos als
          Entwurf an das verbundene TikTok-Konto zu übermitteln.
        </p>
        <p className="mb-2">
          <strong>Rechtsgrundlage:</strong> Ihre Einwilligung (Art. 6 Abs. 1 lit. a DSGVO), die Sie durch Klick auf
          „Connect with TikTok“ und die Zustimmung auf dem TikTok-Bildschirm erteilen. Sie können sie jederzeit mit
          Wirkung für die Zukunft widerrufen.
        </p>
        <p>
          <strong>Widerruf und Löschung:</strong> Über die Schaltfläche „Disconnect“ auf der Seite „/tiktok“ trennen
          Sie die Verbindung; dabei wird die Autorisierung bei Composio gelöscht und der Zugriff bei TikTok
          widerrufen. Alternativ können Sie den Zugriff in den Einstellungen Ihres TikTok-Kontos entziehen oder uns
          per E-Mail an info@goalkeeper91.de um Löschung bitten.
        </p>

        <h2 className="text-2xl font-semibold mb-4 mt-6">10. Google Tag Manager & Cookiebot</h2>
        <p className="mb-2">
          Wir nutzen den Google Tag Manager (Google Ireland Limited, Gordon House, Barrow Street, Dublin 4, Irland),
          um in die Website eingebundene Dienste zentral zu verwalten. Der Tag Manager selbst setzt keine Cookies
          und erhebt keine personenbezogenen Daten, sondern sorgt lediglich für das Auslösen anderer Tags.
        </p>
        <p>
          Die Einwilligungsverwaltung für nicht technisch notwendige Cookies erfolgt über Cookiebot (Cybot A/S,
          Havnegade 39, 1058 Kopenhagen, Dänemark). Details zu den eingesetzten Cookie-Kategorien finden Sie in
          unserer <a href="/legal/cookies" className="underline text-goalyBlue">Cookie-Richtlinie</a>.
        </p>

      <p className="text-sm mt-10 text-gray-400">
        Quelle: <a href="https://www.e-recht24.de" className="underline text-goalyBlue" target="_blank" rel="noopener noreferrer">eRecht24</a> & Ergänzungen für Twitch/Discord/TikTok API.
      </p>
      <p className="text-sm mt-2 text-gray-400">Stand: 24.09.2026</p>
     </div>
    </section>
  );
};

export default Datenschutz;