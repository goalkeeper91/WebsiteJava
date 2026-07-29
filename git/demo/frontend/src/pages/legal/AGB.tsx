import React from "react";

const AGB: React.FC = () => {
  return (
    <section className="w-full min-h-screen bg-slate-950 flex items-center justify-center py-10 px-4">
        <div className="max-w-4xl bg-slate-900 text-white p-8 rounded-lg shadow-lg overflow-auto">
          <h1 className="text-3xl font-bold mb-6">Nutzungsbedingungen (Terms of Service)</h1>

          <p className="mb-4">
            Willkommen auf unserer Plattform! Mit der Nutzung dieser Website und der angebotenen Bot-Dienste
            erklären Sie sich mit den folgenden Bedingungen einverstanden. Bitte lesen Sie diese sorgfältig.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">1. Zugang und Nutzerkonto</h2>
          <p className="mb-4">
            Der Zugang zu unseren Diensten erfolgt über die Anmeldung mit einem Twitch-Account (Self-Service).
            Jeder Twitch-Nutzer kann sich eigenständig registrieren und für seinen eigenen Kanal ein Nutzerkonto
            anlegen; eine Freischaltung durch uns ist nicht erforderlich. Ein Kanalinhaber kann darüber hinaus
            weitere Twitch-Accounts (z. B. seine Moderatoren) als Team-Mitglieder einladen, damit diese seinen
            Kanal im Dashboard mitverwalten können. Unberechtigter Zugang oder der Versuch, sich Zugriff auf ein
            fremdes Konto zu verschaffen, ist untersagt und kann zur Sperrung des Accounts führen.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">2. Nutzung der Bots und Dienste</h2>
          <p className="mb-4">
            Unsere Bots und Dashboard-Funktionen dürfen nur im Rahmen der von uns bereitgestellten Funktionen
            genutzt werden. Jeglicher Versuch, die Dienste zu manipulieren, automatisierte Anfragen außerhalb des
            vorgesehenen Zwecks zu senden oder die Dienste anderweitig zu missbrauchen, ist streng verboten.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">3. Tarife, Bezahlung und Kündigung</h2>
          <p className="mb-4">
            Ein Teil unserer Funktionen (z. B. erweiterte Automatisierungen) ist kostenpflichtigen Tarifen
            vorbehalten, die auf unserer Preisseite beschrieben sind. Die Zahlungsabwicklung erfolgt über unseren
            Zahlungsdienstleister Paddle.com Market Limited ("Paddle"), der als Merchant of Record auch rechtlich
            als Verkäufer auftritt, die Umsatzsteuer abführt und Rechnungen ausstellt. Kostenpflichtige Tarife
            verlängern sich automatisch um den gewählten Abrechnungszeitraum (monatlich oder jährlich), sofern sie
            nicht vorher gekündigt werden. Eine Kündigung ist jederzeit zum Ende des laufenden Abrechnungszeitraums
            über das Paddle-Kundenportal (verlinkt im Dashboard) oder durch Herabstufen auf den kostenlosen Tarif
            möglich. Informationen zu Ihrem Widerrufsrecht bei Vertragsschluss finden Sie in unserer{" "}
            <a href="/legal/widerruf" className="underline text-goalyBlue">Widerrufsbelehrung</a>.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">4. Haftung</h2>
          <p className="mb-4">
            Die Nutzung unserer Dienste erfolgt auf eigene Verantwortung. Wir übernehmen keine Haftung für
            Schäden, Ausfälle oder Datenverlust, die direkt oder indirekt durch die Nutzung der Dienste entstehen.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">5. Datenschutz</h2>
          <p className="mb-4">
            Die Verarbeitung Ihrer Daten erfolgt ausschließlich gemäß unserer{" "}
            <a href="/legal/datenschutz" className="underline text-goalyBlue">Datenschutzerklärung</a>. Dies
            umfasst insbesondere Twitch-ID, Login-Zeitpunkt, Account-Informationen sowie bei kostenpflichtigen
            Tarifen die zur Abrechnung über Paddle notwendigen Daten.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">6. Änderungen der Nutzungsbedingungen</h2>
          <p className="mb-4">
            Wir behalten uns das Recht vor, diese Nutzungsbedingungen jederzeit zu ändern. Änderungen werden
            auf der Website veröffentlicht, und die weitere Nutzung der Dienste gilt als Zustimmung zu den neuen Bedingungen.
          </p>

          <p className="mt-6">
            Sollten einzelne Bestimmungen dieser Nutzungsbedingungen unwirksam sein, bleibt die Wirksamkeit der
            übrigen Bestimmungen unberührt.
          </p>
        </div>
    </section>
  );
};

export default AGB;
