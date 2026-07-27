import { useState } from "react";
import { startCommercial } from "../../features/stream/api";
import type { CommercialRequest } from "../../features/stream/types";

const DURATIONS: CommercialRequest["length"][] = [30, 60, 90, 120, 150, 180];

export default function AdBreakControl() {
  const [length, setLength] = useState<CommercialRequest["length"]>(180);
  const [running, setRunning] = useState(false);
  const [result, setResult] = useState("");
  const [error, setError] = useState("");

  async function handleStart() {
    setRunning(true);
    setError("");
    setResult("");
    try {
      const res = await startCommercial({ length });
      setResult(`✅ Werbepause läuft (${res.length}s). Nächste möglich in ${res.retryAfter}s.`);
    } catch (err: any) {
      setError(err.message || "Konnte keine Werbepause starten");
    } finally {
      setRunning(false);
    }
  }

  return (
    <div className="bg-gray-800 rounded-lg p-4">
      <h2 className="text-lg font-bold mb-3 flex items-center gap-2">
        <span className="text-orange-400">📢</span> Werbepause
      </h2>

      {error && <p className="text-sm text-red-400 mb-2">{error}</p>}
      {result && <p className="text-sm text-green-400 mb-2">{result}</p>}

      <div className="flex flex-wrap gap-3 items-center">
        <select
          value={length}
          onChange={(e) => setLength(Number(e.target.value) as CommercialRequest["length"])}
          className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-sm text-white"
        >
          {DURATIONS.map((d) => (
            <option key={d} value={d}>
              {d} Sekunden
            </option>
          ))}
        </select>
        <button
          onClick={handleStart}
          disabled={running}
          className="px-4 py-2 bg-orange-600 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-orange-500 transition-colors text-sm font-semibold"
        >
          {running ? "Starte…" : "Werbepause starten"}
        </button>
      </div>
    </div>
  );
}
