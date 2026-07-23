import { useEffect, useState } from "react";
import { getAutomodSettings, updateAutomodSettings } from "../../features/automod/api";

function linesToArray(text: string): string[] {
  return text
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
}

export default function AutomodDashboard() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);

  const [enabled, setEnabled] = useState(false);
  const [blockedWordsText, setBlockedWordsText] = useState("");
  const [linkFilterEnabled, setLinkFilterEnabled] = useState(false);
  const [allowedDomainsText, setAllowedDomainsText] = useState("");

  useEffect(() => {
    getAutomodSettings()
      .then((settings) => {
        setEnabled(settings.enabled);
        setBlockedWordsText((settings.blocked_words || []).join("\n"));
        setLinkFilterEnabled(settings.link_filter_enabled);
        setAllowedDomainsText((settings.allowed_domains || []).join("\n"));
      })
      .catch((err) => setError(err.message || "Fehler beim Laden"))
      .finally(() => setLoading(false));
  }, []);

  async function handleSave() {
    setSaving(true);
    setError("");
    setSaved(false);

    try {
      await updateAutomodSettings({
        enabled,
        blocked_words: linesToArray(blockedWordsText),
        link_filter_enabled: linkFilterEnabled,
        allowed_domains: linesToArray(allowedDomainsText),
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err: any) {
      setError(err.message || "Fehler beim Speichern");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="max-w-3xl mx-auto p-6 animate-pulse space-y-4">
        <div className="h-8 bg-gray-800 rounded w-1/3"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Automod</h1>
        <p className="text-sm text-gray-400 mt-1">
          Löscht automatisch Nachrichten mit verbotenen Wörtern oder nicht erlaubten Links und
          verhängt einen eskalierenden Timeout. Mods und du selbst sind immer ausgenommen.
        </p>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-6 space-y-6">
        <label className="flex items-center justify-between cursor-pointer">
          <div>
            <p className="font-semibold text-white">Automod aktiviert</p>
            <p className="text-xs text-gray-400 mt-0.5">
              Ohne diese Option wird nichts geprüft, auch wenn unten Wörter/Domains eingetragen sind.
            </p>
          </div>
          <input
            type="checkbox"
            checked={enabled}
            onChange={(e) => setEnabled(e.target.checked)}
            className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
          />
        </label>

        <div>
          <label className="block text-sm font-medium mb-1">Blockliste</label>
          <p className="text-xs text-gray-400 mb-2">Ein Wort/Phrase pro Zeile.</p>
          <textarea
            value={blockedWordsText}
            onChange={(e) => setBlockedWordsText(e.target.value)}
            rows={6}
            placeholder={"beispielwort\nverbotene phrase"}
            className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>

        <div className="pt-4 border-t border-gray-700">
          <label className="flex items-center justify-between cursor-pointer mb-3">
            <p className="font-semibold text-white">Link-Filter</p>
            <input
              type="checkbox"
              checked={linkFilterEnabled}
              onChange={(e) => setLinkFilterEnabled(e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
            />
          </label>

          <label className="block text-sm font-medium mb-1">Erlaubte Domains</label>
          <p className="text-xs text-gray-400 mb-2">
            Eine Domain pro Zeile (z.B. twitch.tv). Links zu anderen Domains werden gelöscht, wenn
            der Link-Filter aktiv ist.
          </p>
          <textarea
            value={allowedDomainsText}
            onChange={(e) => setAllowedDomainsText(e.target.value)}
            rows={4}
            placeholder={"twitch.tv\ndiscord.gg"}
            className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>

        <div className="flex items-center gap-3 pt-2">
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-2 bg-green-600 rounded-lg hover:bg-green-500 disabled:opacity-50 transition-colors"
          >
            {saving ? "Speichern…" : "Speichern"}
          </button>
          {saved && <span className="text-sm text-green-400">✓ Gespeichert</span>}
        </div>
      </div>
    </div>
  );
}
