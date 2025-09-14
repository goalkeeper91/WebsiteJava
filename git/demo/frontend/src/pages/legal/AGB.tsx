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
            Der Zugang zu unseren Bot-Diensten erfolgt ausschließlich über die Anmeldung mit einem Twitch-Account.
            Es ist nur Moderatoren („Mods“) gestattet, sich einzuloggen. Wir überprüfen die Berechtigung automatisch
            über Twitch-Datenbankabgleiche. Unberechtigter Zugang ist untersagt und kann zur Sperrung des Accounts führen.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">2. Nutzung der Bots</h2>
          <p className="mb-4">
            Unsere Bots dürfen nur im Rahmen der von uns bereitgestellten Funktionen genutzt werden. Jeglicher
            Versuch, die Bots zu manipulieren, automatisierte Anfragen außerhalb des vorgesehenen Zwecks zu senden
            oder die Dienste anderweitig zu missbrauchen, ist streng verboten.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">3. Haftung</h2>
          <p className="mb-4">
            Die Nutzung der Bot-Dienste erfolgt auf eigene Verantwortung. Wir übernehmen keine Haftung für
            Schäden, Ausfälle oder Datenverlust, die direkt oder indirekt durch die Nutzung der Bots entstehen.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">4. Datenschutz</h2>
          <p className="mb-4">
            Die Verarbeitung Ihrer Daten erfolgt ausschließlich gemäß unserer Datenschutzerklärung. Dies umfasst
            insbesondere Twitch-ID, Login-Zeitpunkt und Account-Informationen, um sicherzustellen, dass nur
            berechtigte Moderatoren Zugriff haben.
          </p>

          <h2 className="text-2xl font-semibold mt-6 mb-2">5. Änderungen der Nutzungsbedingungen</h2>
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
