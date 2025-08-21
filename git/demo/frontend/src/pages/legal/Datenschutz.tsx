import React from "react";

const Datenschutz: React.FC = () => {
  return (
    <div className="max-w-3xl mx-auto p-4 text-gray-900 bg-white rounded shadow">
      <h1 className="text-2xl font-bold mb-4">Datenschutzerklärung</h1>

      <p className="mb-4">
        Hier steht deine Datenschutzerklärung. Dieser Text dient als Platzhalter
        und sollte später durch eine vollständige Datenschutzerklärung ersetzt werden.
      </p>

      <p>
        <strong>Beispiel:</strong> Wir nehmen den Schutz Ihrer persönlichen Daten sehr ernst.
        Personenbezogene Daten werden vertraulich behandelt und entsprechend der gesetzlichen
        Datenschutzvorschriften verarbeitet.
      </p>
    </div>
  );
};

export default Datenschutz;
