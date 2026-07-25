import { useEffect, useState } from "react";
import { getAutomodSettings, updateAutomodSettings, getAutomodEvents } from "../../features/automod/api";
import type { AutomodEvent } from "../../features/automod/types";

function linesToArray(text: string): string[] {
  return text
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
}

function formatRelative(dateString: string): string {
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(dateString).getTime()) / 1000));
  if (seconds < 60) return "gerade eben";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `vor ${minutes} Min.`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `vor ${hours} Std.`;
  return `vor ${Math.floor(hours / 24)} Tag(en)`;
}

function formatTimeout(seconds: number): string {
  if (seconds >= 86400) return `${Math.round(seconds / 86400)}d`;
  if (seconds >= 3600) return `${Math.round(seconds / 3600)}h`;
  if (seconds >= 60) return `${Math.round(seconds / 60)}m`;
  return `${seconds}s`;
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
  const [capsFilterEnabled, setCapsFilterEnabled] = useState(false);
  const [symbolFilterEnabled, setSymbolFilterEnabled] = useState(false);
  const [emoteFilterEnabled, setEmoteFilterEnabled] = useState(false);
  const [emoteThreshold, setEmoteThreshold] = useState(6);
  const [repetitionFilterEnabled, setRepetitionFilterEnabled] = useState(false);
  const [exemptVips, setExemptVips] = useState(true);
  const [exemptUsersText, setExemptUsersText] = useState("");

  const [events, setEvents] = useState<AutomodEvent[]>([]);
  const [eventsPage, setEventsPage] = useState(1);
  const [eventsTotalPages, setEventsTotalPages] = useState(0);
  const [eventsLoading, setEventsLoading] = useState(true);

  useEffect(() => {
    getAutomodSettings()
      .then((settings) => {
        setEnabled(settings.enabled);
        setBlockedWordsText((settings.blocked_words || []).join("\n"));
        setLinkFilterEnabled(settings.link_filter_enabled);
        setAllowedDomainsText((settings.allowed_domains || []).join("\n"));
        setCapsFilterEnabled(settings.caps_filter_enabled);
        setSymbolFilterEnabled(settings.symbol_filter_enabled);
        setEmoteFilterEnabled(settings.emote_filter_enabled);
        setEmoteThreshold(settings.emote_threshold || 6);
        setRepetitionFilterEnabled(settings.repetition_filter_enabled);
        setExemptVips(settings.exempt_vips);
        setExemptUsersText((settings.exempt_users || []).join("\n"));
      })
      .catch((err) => setError(err.message || "Fehler beim Laden"))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    setEventsLoading(true);
    getAutomodEvents(eventsPage, 10)
      .then((data) => {
        setEvents(data.data || []);
        setEventsTotalPages(data.total_pages || 0);
      })
      .catch((err) => console.error("Fehler beim Laden der Automod-Events:", err))
      .finally(() => setEventsLoading(false));
  }, [eventsPage]);

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
        caps_filter_enabled: capsFilterEnabled,
        symbol_filter_enabled: symbolFilterEnabled,
        emote_filter_enabled: emoteFilterEnabled,
        emote_threshold: emoteThreshold,
        repetition_filter_enabled: repetitionFilterEnabled,
        exempt_vips: exemptVips,
        exempt_users: linesToArray(exemptUsersText),
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
          Löscht automatisch Nachrichten mit verbotenen Wörtern, nicht erlaubten Links oder
          Spam-Mustern und verhängt einen eskalierenden Timeout. Mods und du selbst sind immer
          ausgenommen.
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
              Ohne diese Option wird nichts geprüft, auch wenn unten etwas konfiguriert ist.
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

        <div className="pt-4 border-t border-gray-700 space-y-3">
          <p className="font-semibold text-white">Spam-Filter</p>
          <p className="text-xs text-gray-400 -mt-2">
            Feste Schwellwerte (nicht einstellbar): CAPS ≥70% Großbuchstaben, Symbole ≥50%
            Sonderzeichen, dieselbe Nachricht 3× in 30 Sekunden. Die Emote-Anzahl kannst du unten
            selbst festlegen — Emotes im richtigen Moment (Esports-Highlights o.ä.) sind oft
            gewollt, kein Spam.
          </p>

          <label className="flex items-center justify-between cursor-pointer">
            <span className="text-sm text-gray-200">CAPS-Spam</span>
            <input
              type="checkbox"
              checked={capsFilterEnabled}
              onChange={(e) => setCapsFilterEnabled(e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
            />
          </label>
          <label className="flex items-center justify-between cursor-pointer">
            <span className="text-sm text-gray-200">Symbol-Spam</span>
            <input
              type="checkbox"
              checked={symbolFilterEnabled}
              onChange={(e) => setSymbolFilterEnabled(e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
            />
          </label>
          <div>
            <label className="flex items-center justify-between cursor-pointer">
              <span className="text-sm text-gray-200">Emote-Spam</span>
              <input
                type="checkbox"
                checked={emoteFilterEnabled}
                onChange={(e) => setEmoteFilterEnabled(e.target.checked)}
                className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
              />
            </label>
            {emoteFilterEnabled && (
              <div className="flex items-center gap-2 mt-2 pl-4">
                <label className="text-xs text-gray-400">Ab wie vielen Emotes pro Nachricht:</label>
                <input
                  type="number"
                  min={1}
                  value={emoteThreshold}
                  onChange={(e) => setEmoteThreshold(Math.max(1, Number(e.target.value)))}
                  className="w-20 p-1 rounded bg-gray-900 border border-gray-700 text-white text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
            )}
          </div>
          <label className="flex items-center justify-between cursor-pointer">
            <span className="text-sm text-gray-200">Nachrichten-Wiederholung</span>
            <input
              type="checkbox"
              checked={repetitionFilterEnabled}
              onChange={(e) => setRepetitionFilterEnabled(e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
            />
          </label>
        </div>

        <div className="pt-4 border-t border-gray-700 space-y-3">
          <p className="font-semibold text-white">Ausnahmen</p>

          <label className="flex items-center justify-between cursor-pointer">
            <div>
              <span className="text-sm text-gray-200">VIPs ausnehmen</span>
              <p className="text-xs text-gray-400">VIPs werden nie automoderiert.</p>
            </div>
            <input
              type="checkbox"
              checked={exemptVips}
              onChange={(e) => setExemptVips(e.target.checked)}
              className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-green-600 focus:ring-2 focus:ring-green-500"
            />
          </label>

          <div>
            <label className="block text-sm font-medium mb-1">Vertraute Nutzer</label>
            <p className="text-xs text-gray-400 mb-2">
              Ein Twitch-Nutzername pro Zeile - diese Zuschauer werden nie automoderiert (Ersatz
              für "Regulars", solange es keine Watchtime-Erkennung gibt).
            </p>
            <textarea
              value={exemptUsersText}
              onChange={(e) => setExemptUsersText(e.target.value)}
              rows={3}
              placeholder={"vertrauter_zuschauer"}
              className="w-full p-2 rounded bg-gray-900 border border-gray-700 text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
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

      <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
        <h3 className="font-semibold text-white mb-3">Letzte Verstöße</h3>

        {eventsLoading && <p className="text-sm text-gray-500">Lade…</p>}

        {!eventsLoading && events.length === 0 && (
          <p className="text-sm text-gray-500">Noch keine Verstöße erfasst.</p>
        )}

        <div className="space-y-2">
          {events.map((event) => (
            <div key={event.id} className="bg-gray-900 rounded-lg p-3 text-sm">
              <div className="flex items-center justify-between flex-wrap gap-2">
                <span className="font-semibold text-white">{event.offender_name}</span>
                <span className="text-xs text-gray-500">{formatRelative(event.created_at)}</span>
              </div>
              <p className="text-gray-300 mt-1">
                {event.reason} · Timeout {formatTimeout(event.timeout_seconds)}
              </p>
              {event.message_excerpt && (
                <p className="text-gray-500 text-xs mt-1 break-words">"{event.message_excerpt}"</p>
              )}
            </div>
          ))}
        </div>

        {!eventsLoading && eventsTotalPages > 1 && (
          <div className="flex justify-between items-center mt-4 pt-4 border-t border-gray-700">
            <button
              disabled={eventsPage <= 1}
              onClick={() => setEventsPage((p) => p - 1)}
              className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors text-sm"
            >
              ← Zurück
            </button>
            <span className="text-gray-400 text-sm">
              Seite {eventsPage} von {eventsTotalPages}
            </span>
            <button
              disabled={eventsPage >= eventsTotalPages}
              onClick={() => setEventsPage((p) => p + 1)}
              className="px-3 py-1.5 bg-gray-700 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors text-sm"
            >
              Weiter →
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
