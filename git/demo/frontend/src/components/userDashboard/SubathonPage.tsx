import { useCallback, useEffect, useState } from "react";
import { Play, Pause, RotateCcw, Copy, Radio } from "lucide-react";
import {
  getSubathonState,
  startSubathon,
  pauseSubathon,
  resetSubathon,
  updateSubathonSettings,
} from "../../features/subathon/api";
import type { SubathonState } from "../../features/subathon/types";

function formatTime(totalSeconds: number): string {
  const clamped = Math.max(0, Math.floor(totalSeconds));
  const h = Math.floor(clamped / 3600).toString().padStart(2, "0");
  const m = Math.floor((clamped % 3600) / 60).toString().padStart(2, "0");
  const s = Math.floor(clamped % 60).toString().padStart(2, "0");
  return `${h}:${m}:${s}`;
}

export default function SubathonPage() {
  const [state, setState] = useState<SubathonState | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [settingsForm, setSettingsForm] = useState({ initialTime: 30, subTime: 120, bitsTime: 30 });
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    try {
      const data = await getSubathonState();
      setState(data);
      setSettingsForm({
        initialTime: data.initialTime,
        subTime: data.subTime,
        bitsTime: data.bitsTime,
      });
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Laden");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
    const interval = setInterval(load, 1000);
    return () => clearInterval(interval);
  }, [load]);

  async function handleStart() {
    try {
      setState(await startSubathon());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Starten");
    }
  }

  async function handlePause() {
    try {
      setState(await pauseSubathon());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Pausieren");
    }
  }

  async function handleReset() {
    if (!confirm("Timer wirklich zurücksetzen? Stats und Log werden geleert.")) return;
    try {
      setState(await resetSubathon());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Zurücksetzen");
    }
  }

  async function handleSaveSettings() {
    setSaving(true);
    try {
      setState(await updateSubathonSettings(settingsForm));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Speichern");
    } finally {
      setSaving(false);
    }
  }

  const overlayUrl = state
    ? `https://timer.goalkeeper91.de/overlay?user=${state.userId ?? ""}`
    : "";

  function copyOverlayUrl() {
    navigator.clipboard.writeText(overlayUrl);
    alert("OBS-Overlay-URL kopiert!");
  }

  if (loading) {
    return (
      <div className="p-6 animate-pulse space-y-4">
        <div className="h-8 bg-gray-800 rounded w-1/3"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  if (!state) {
    return (
      <div className="p-6">
        <div className="bg-gray-800 rounded-xl p-6 text-red-400">
          {error || "Timer konnte nicht geladen werden."}
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <Radio className="w-5 h-5 text-purple-400" />
            Subathon Timer
          </h2>
          <p className="text-sm text-gray-400 mt-1">
            {state.isRunning ? "Läuft" : "Pausiert"} · Subs verlängern um {state.subTime}s, 100 Bits um {state.bitsTime}s
          </p>
        </div>
        <span
          className={`text-xs px-3 py-1 rounded-full font-semibold ${
            state.isRunning ? "bg-green-600 text-white" : "bg-gray-700 text-gray-300"
          }`}
        >
          {state.isRunning ? "● LIVE" : "○ STANDBY"}
        </span>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      <div className="grid md:grid-cols-3 gap-6">
        {/* Timer + Controls */}
        <div className="md:col-span-2 bg-gray-800 border border-gray-700 rounded-xl p-6 space-y-6">
          <div className="text-center">
            <div className="text-6xl md:text-7xl font-bold text-purple-400 tracking-wider tabular-nums">
              {formatTime(state.timeRemaining)}
            </div>
          </div>

          <div className="grid grid-cols-3 gap-3">
            <StatBox label="Subs" value={state.totalSubs} />
            <StatBox label="Bits" value={state.totalBits} />
            <StatBox label="Events" value={state.totalEvents} />
          </div>

          <div className="flex gap-3">
            <button
              onClick={handleStart}
              disabled={state.isRunning}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-green-600 hover:bg-green-500 disabled:bg-green-900/50 disabled:text-gray-500 text-white rounded-lg font-semibold transition-colors"
            >
              <Play className="w-4 h-4" />
              Start
            </button>
            <button
              onClick={handlePause}
              disabled={!state.isRunning}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-yellow-600 hover:bg-yellow-500 disabled:bg-yellow-900/50 disabled:text-gray-500 text-white rounded-lg font-semibold transition-colors"
            >
              <Pause className="w-4 h-4" />
              Pause
            </button>
            <button
              onClick={handleReset}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-red-600 hover:bg-red-500 text-white rounded-lg font-semibold transition-colors"
            >
              <RotateCcw className="w-4 h-4" />
              Reset
            </button>
          </div>

          <div className="grid sm:grid-cols-3 gap-3 pt-2 border-t border-gray-700">
            <div>
              <label className="block text-xs text-gray-400 mb-1">Start-Zeit (Min.)</label>
              <input
                type="number"
                min={1}
                value={settingsForm.initialTime}
                onChange={(e) =>
                  setSettingsForm({ ...settingsForm, initialTime: Number(e.target.value) })
                }
                className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-400 mb-1">Sub (+Sek.)</label>
              <input
                type="number"
                min={0}
                value={settingsForm.subTime}
                onChange={(e) =>
                  setSettingsForm({ ...settingsForm, subTime: Number(e.target.value) })
                }
                className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-400 mb-1">100 Bits (+Sek.)</label>
              <input
                type="number"
                min={0}
                value={settingsForm.bitsTime}
                onChange={(e) =>
                  setSettingsForm({ ...settingsForm, bitsTime: Number(e.target.value) })
                }
                className="w-full bg-gray-900 border border-gray-700 text-white p-2 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              />
            </div>
          </div>
          <button
            onClick={handleSaveSettings}
            disabled={saving}
            className="w-full px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-white rounded-lg transition-colors text-sm"
          >
            {saving ? "Speichern..." : "Einstellungen speichern"}
          </button>
        </div>

        {/* Side panel: OBS + Log */}
        <div className="space-y-6">
          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
            <h3 className="font-semibold text-white mb-3">OBS Overlay</h3>
            <input
              type="text"
              readOnly
              value={overlayUrl}
              className="w-full bg-gray-900 border border-gray-700 text-gray-400 text-xs p-2 rounded-lg mb-3"
            />
            <button
              onClick={copyOverlayUrl}
              className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-lg transition-colors text-sm"
            >
              <Copy className="w-4 h-4" />
              URL kopieren
            </button>
            <p className="text-xs text-gray-500 mt-2">
              Als Browser-Quelle in OBS einfügen.
            </p>
          </div>

          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
            <h3 className="font-semibold text-white mb-3">Event-Log</h3>
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {state.eventLog.length === 0 && (
                <p className="text-sm text-gray-500">Noch keine Events.</p>
              )}
              {state.eventLog.map((entry, i) => (
                <div
                  key={`${entry.time}-${i}`}
                  className="text-xs bg-gray-900 border-l-2 border-purple-500 rounded px-3 py-2"
                >
                  <span className="text-purple-400">[{entry.time}]</span>{" "}
                  <span className="text-gray-300">{entry.text}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function StatBox({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-3 text-center">
      <div className="text-xs text-gray-400 uppercase tracking-wide">{label}</div>
      <div className="text-2xl font-bold text-white mt-1">{value}</div>
    </div>
  );
}
