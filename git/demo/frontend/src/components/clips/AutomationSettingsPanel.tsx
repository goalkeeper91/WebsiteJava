import { useEffect, useState } from "react";
import { Sparkles, Save, Lock } from "lucide-react";
import {
  getAutomationSettings,
  updateAutomationSettings,
} from "../../features/clips/api";
import type {
  AutomationSettings,
  AITone,
  AIStyle,
  PostingFrequency,
  Platform,
} from "../../features/clips/types";
import {
  LINK_SHARE_PLATFORMS,
  NATIVE_UPLOAD_PLATFORMS,
} from "../../features/clips/types";

const TONE_OPTIONS: AITone[] = ["neutral", "casual", "professional", "energetic", "funny"];
const STYLE_OPTIONS: AIStyle[] = ["standard", "short", "detailed", "viral", "storytelling"];
const FREQUENCY_OPTIONS: PostingFrequency[] = ["realtime", "hourly", "daily", "weekly", "manual"];

const PLATFORM_LABELS: Record<Platform, string> = {
  discord: "Discord",
  twitter: "X / Twitter",
  tiktok: "TikTok",
  youtube: "YouTube",
  instagram: "Instagram",
};

export default function AutomationSettingsPanel() {
  const [settings, setSettings] = useState<AutomationSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [notAvailable, setNotAvailable] = useState(false);

  useEffect(() => {
    load();
  }, []);

  async function load() {
    setLoading(true);
    setError(null);
    try {
      const data = await getAutomationSettings();
      setSettings(data);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Fehler beim Laden";
      if (message.toLowerCase().includes("pro oder premium")) {
        setNotAvailable(true);
      } else {
        setError(message);
      }
    } finally {
      setLoading(false);
    }
  }

  function togglePlatform(platform: Platform) {
    if (!settings || !LINK_SHARE_PLATFORMS.includes(platform)) return;
    const current = settings.target_platforms || [];
    const next = current.includes(platform)
      ? current.filter((p) => p !== platform)
      : [...current, platform];
    setSettings({ ...settings, target_platforms: next });
  }

  async function handleSave() {
    if (!settings) return;
    setSaving(true);
    setError(null);
    try {
      const updated = await updateAutomationSettings({
        is_enabled: settings.is_enabled,
        ai_tone: settings.ai_tone,
        ai_style: settings.ai_style,
        use_hashtags: settings.use_hashtags,
        max_hashtags: settings.max_hashtags,
        min_clip_duration: settings.min_clip_duration,
        max_clip_duration: settings.max_clip_duration,
        only_verified_clips: settings.only_verified_clips,
        target_platforms: settings.target_platforms,
        auto_post_enabled: settings.auto_post_enabled,
        posting_frequency: settings.posting_frequency,
        notify_on_post: settings.notify_on_post,
        notify_on_error: settings.notify_on_error,
      });
      setSettings(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Fehler beim Speichern");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="animate-pulse space-y-4">
        <div className="h-8 bg-gray-800 rounded w-1/3"></div>
        <div className="h-64 bg-gray-800 rounded"></div>
      </div>
    );
  }

  if (notAvailable) {
    return (
      <div className="bg-gray-800 rounded-xl p-8 text-center">
        <Lock className="w-16 h-16 text-gray-600 mx-auto mb-4" />
        <h2 className="text-2xl font-bold mb-2">Clip-Automatisierung</h2>
        <p className="text-gray-400 mb-6">
          Automatische Clip-Erkennung und Multi-Plattform-Posting sind nur in Pro und Premium verfügbar.
        </p>
        <a
          href="/dashboard/subscription"
          className="inline-block px-6 py-3 bg-green-600 hover:bg-green-500 rounded-lg transition-colors font-semibold"
        >
          Jetzt upgraden
        </a>
      </div>
    );
  }

  if (!settings) {
    return (
      <div className="bg-gray-800 rounded-xl p-6 text-red-400">
        {error || "Einstellungen konnten nicht geladen werden."}
      </div>
    );
  }

  return (
    <div className="bg-gray-800 rounded-xl p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Sparkles className="w-5 h-5 text-purple-400" />
          <h2 className="text-xl font-bold">Automation-Einstellungen</h2>
        </div>
        <label className="flex items-center gap-2 cursor-pointer">
          <span className="text-sm text-gray-400">Aktiviert</span>
          <input
            type="checkbox"
            checked={settings.is_enabled}
            onChange={(e) => setSettings({ ...settings, is_enabled: e.target.checked })}
            className="w-5 h-5 accent-purple-600"
          />
        </label>
      </div>

      {error && (
        <div className="bg-red-900/30 border border-red-500/50 text-red-300 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      <div className="grid md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm text-gray-400 mb-1">KI-Tonalität</label>
          <select
            value={settings.ai_tone}
            onChange={(e) => setSettings({ ...settings, ai_tone: e.target.value as AITone })}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2"
          >
            {TONE_OPTIONS.map((t) => (
              <option key={t} value={t}>{t}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm text-gray-400 mb-1">Caption-Stil</label>
          <select
            value={settings.ai_style}
            onChange={(e) => setSettings({ ...settings, ai_style: e.target.value as AIStyle })}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2"
          >
            {STYLE_OPTIONS.map((s) => (
              <option key={s} value={s}>{s}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm text-gray-400 mb-1">Min. Clip-Länge (Sek.)</label>
          <input
            type="number"
            min={5}
            max={300}
            value={settings.min_clip_duration}
            onChange={(e) => setSettings({ ...settings, min_clip_duration: Number(e.target.value) })}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2"
          />
        </div>

        <div>
          <label className="block text-sm text-gray-400 mb-1">Max. Clip-Länge (Sek.)</label>
          <input
            type="number"
            min={5}
            max={300}
            value={settings.max_clip_duration}
            onChange={(e) => setSettings({ ...settings, max_clip_duration: Number(e.target.value) })}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2"
          />
        </div>

        <div>
          <label className="block text-sm text-gray-400 mb-1">Posting-Frequenz</label>
          <select
            value={settings.posting_frequency}
            onChange={(e) => setSettings({ ...settings, posting_frequency: e.target.value as PostingFrequency })}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2"
          >
            {FREQUENCY_OPTIONS.map((f) => (
              <option key={f} value={f}>{f}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm text-gray-400 mb-1">Max. Hashtags</label>
          <input
            type="number"
            min={0}
            max={30}
            disabled={!settings.use_hashtags}
            value={settings.max_hashtags}
            onChange={(e) => setSettings({ ...settings, max_hashtags: Number(e.target.value) })}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 disabled:opacity-50"
          />
        </div>
      </div>

      <div className="flex flex-wrap gap-6">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={settings.use_hashtags}
            onChange={(e) => setSettings({ ...settings, use_hashtags: e.target.checked })}
            className="w-4 h-4 accent-purple-600"
          />
          <span className="text-sm">Hashtags generieren</span>
        </label>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={settings.only_verified_clips}
            onChange={(e) => setSettings({ ...settings, only_verified_clips: e.target.checked })}
            className="w-4 h-4 accent-purple-600"
          />
          <span className="text-sm">Nur verifizierte Clips posten</span>
        </label>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={settings.auto_post_enabled}
            onChange={(e) => setSettings({ ...settings, auto_post_enabled: e.target.checked })}
            className="w-4 h-4 accent-purple-600"
          />
          <span className="text-sm">Automatisch posten</span>
        </label>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={settings.notify_on_post}
            onChange={(e) => setSettings({ ...settings, notify_on_post: e.target.checked })}
            className="w-4 h-4 accent-purple-600"
          />
          <span className="text-sm">Benachrichtigen bei Post</span>
        </label>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={settings.notify_on_error}
            onChange={(e) => setSettings({ ...settings, notify_on_error: e.target.checked })}
            className="w-4 h-4 accent-purple-600"
          />
          <span className="text-sm">Benachrichtigen bei Fehler</span>
        </label>
      </div>

      <div>
        <label className="block text-sm text-gray-400 mb-2">Ziel-Plattformen</label>
        <div className="flex flex-wrap gap-3">
          {[...LINK_SHARE_PLATFORMS, ...NATIVE_UPLOAD_PLATFORMS].map((platform) => {
            const isLinkShare = LINK_SHARE_PLATFORMS.includes(platform);
            const active = (settings.target_platforms || []).includes(platform);
            return (
              <button
                key={platform}
                type="button"
                disabled={!isLinkShare}
                onClick={() => togglePlatform(platform)}
                title={isLinkShare ? undefined : "Natives Video-Posting folgt in einer späteren Phase"}
                className={`px-4 py-2 rounded-lg text-sm font-medium border transition-colors ${
                  active
                    ? "bg-purple-600 border-purple-500 text-white"
                    : "bg-gray-900 border-gray-700 text-gray-400"
                } ${!isLinkShare ? "opacity-40 cursor-not-allowed" : "hover:border-purple-500"}`}
              >
                {PLATFORM_LABELS[platform]}
                {!isLinkShare && <span className="ml-2 text-xs">(bald)</span>}
              </button>
            );
          })}
        </div>
        <p className="text-xs text-gray-500 mt-2">
          v1: Clips werden als Link geteilt (Discord/X). Natives Re-Upload zu TikTok/YouTube/Instagram folgt später.
        </p>
      </div>

      <div className="flex justify-end">
        <button
          onClick={handleSave}
          disabled={saving}
          className="flex items-center gap-2 px-6 py-3 bg-purple-600 hover:bg-purple-500 rounded-lg font-semibold transition-colors disabled:opacity-50"
        >
          <Save className="w-4 h-4" />
          {saving ? "Speichern..." : "Speichern"}
        </button>
      </div>
    </div>
  );
}
