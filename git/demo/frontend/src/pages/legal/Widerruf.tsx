import React from "react";
import Seo from "../../components/Seo";

const Widerruf: React.FC = () => {
  return (
    <section className="w-full min-h-screen bg-slate-950 flex items-center justify-center py-10 px-4">
      <Seo
        title="Widerrufsbelehrung"
        description="Widerrufsbelehrung für kostenpflichtige Tarife bei Goalkeeper91: Fristen, Muster-Widerrufsformular und Hinweise zum vorzeitigen Erlöschen des Widerrufsrechts."
        path="/legal/widerruf"
      />
      <div className="max-w-4xl bg-slate-900 text-white p-8 rounded-lg shadow-lg overflow-auto">
        <h1 className="text-3xl font-bold mb-6">Widerrufsbelehrung</h1>

        <p className="mb-4">
          Diese Belehrung gilt für den Abschluss kostenpflichtiger Tarife (digitale Dienstleistungen) über unsere
          Plattform. Vertragspartner für die Zahlungsabwicklung ist Paddle.com Market Limited, 30 Old Bailey,
          London, EC4M 7AU, Vereinigtes Königreich ("Paddle"), das als Merchant of Record auch rechtlich als
          Verkäufer auftritt.
        </p>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Widerrufsrecht</h2>
        <p className="mb-4">
          Sie haben das Recht, binnen vierzehn Tagen ohne Angabe von Gründen diesen Vertrag zu widerrufen. Die
          Widerrufsfrist beträgt vierzehn Tage ab dem Tag des Vertragsschlusses.
        </p>
        <p className="mb-4">
          Um Ihr Widerrufsrecht auszuüben, müssen Sie uns (Goalkeeper91 c/o NextlevelNation, Inhaber: Christian
          Steinbach, Stettener Weg 2, 89584 Ehingen (Donau), E-Mail: info@goalkeeper91.de) mittels einer eindeutigen
          Erklärung (z. B. per E-Mail) über Ihren Entschluss, diesen Vertrag zu widerrufen, informieren. Sie können
          dafür das unten stehende Muster-Widerrufsformular verwenden, das jedoch nicht vorgeschrieben ist.
        </p>
        <p className="mb-4">
          Zur Wahrung der Widerrufsfrist reicht es aus, dass Sie die Mitteilung über die Ausübung des
          Widerrufsrechts vor Ablauf der Widerrufsfrist absenden.
        </p>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Folgen des Widerrufs</h2>
        <p className="mb-4">
          Wenn Sie diesen Vertrag widerrufen, haben wir Ihnen alle Zahlungen, die wir von Ihnen erhalten haben,
          unverzüglich und spätestens binnen vierzehn Tagen ab dem Tag zurückzuzahlen, an dem die Mitteilung über
          Ihren Widerruf dieses Vertrags bei uns eingegangen ist. Für diese Rückzahlung verwenden wir dasselbe
          Zahlungsmittel, das Sie bei der ursprünglichen Transaktion eingesetzt haben, es sei denn, mit Ihnen wurde
          ausdrücklich etwas anderes vereinbart; in keinem Fall werden Ihnen wegen dieser Rückzahlung Entgelte
          berechnet.
        </p>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Vorzeitiges Erlöschen des Widerrufsrechts</h2>
        <p className="mb-4">
          Ihr Widerrufsrecht erlischt vorzeitig, wenn wir mit der Ausführung des Vertrags erst begonnen haben,
          nachdem Sie ausdrücklich zugestimmt haben, dass wir vor Ablauf der Widerrufsfrist mit der Ausführung des
          Vertrags beginnen, und Sie Ihre Kenntnis davon bestätigt haben, dass Sie durch Ihre Zustimmung mit
          Beginn der Ausführung des Vertrags Ihr Widerrufsrecht verlieren. Da unsere kostenpflichtigen Tarife den
          sofortigen Zugriff auf digitale Funktionen gewähren, holen wir diese Zustimmung direkt auf unserer
          Preisseite über eine gesondert zu bestätigende Checkbox ein, bevor der Bezahlvorgang gestartet wird.
        </p>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Muster-Widerrufsformular</h2>
        <p className="mb-2">
          (Wenn Sie den Vertrag widerrufen wollen, füllen Sie bitte dieses Formular aus und senden Sie es zurück.)
        </p>
        <p className="mb-1">An: Goalkeeper91 c/o NextlevelNation, Stettener Weg 2, 89584 Ehingen (Donau), E-Mail: info@goalkeeper91.de</p>
        <p className="mb-1">Hiermit widerrufe(n) ich/wir den von mir/uns abgeschlossenen Vertrag über die Erbringung der folgenden Dienstleistung:</p>
        <p className="mb-1">Bestellt am / erhalten am:</p>
        <p className="mb-1">Name des/der Verbraucher(s):</p>
        <p className="mb-1">Anschrift des/der Verbraucher(s):</p>
        <p className="mb-1">Unterschrift des/der Verbraucher(s) (nur bei Mitteilung auf Papier)</p>
        <p className="mb-4">Datum:</p>

        <h2 className="text-2xl font-semibold mt-6 mb-2">Kündigung laufender Abonnements</h2>
        <p>
          Unabhängig vom gesetzlichen Widerrufsrecht können Sie ein laufendes kostenpflichtiges Abonnement
          jederzeit zum Ende des aktuellen Abrechnungszeitraums über unsere{" "}
          <a href="/vertrag-kuendigen" className="underline text-goalyBlue">Kündigungsseite</a> (auch über den
          Link "Kündigen" in der Fußzeile jeder Seite erreichbar) oder das Paddle-Kundenportal kündigen bzw. auf
          den kostenlosen Tarif herabstufen - siehe unsere{" "}
          <a href="/legal/agb" className="underline text-goalyBlue">Nutzungsbedingungen</a>.
        </p>
      </div>
    </section>
  );
};

export default Widerruf;
