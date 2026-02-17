import { useEffect, useState } from "react";
import { Link, Zap, CheckCircle, XCircle, AlertTriangle, RefreshCw } from "lucide-react";

interface N8NIntegration {
  id: number;
  enabled: boolean;
  webhookBaseUrl: string;
  workflowsUsed: number;
  votesThisMonth: number;
  lastResetAt: string;
  isReady: boolean;
}

export default function N8NIntegrationSetup() {
  const [integration, setIntegration] = useState<N8NIntegration | null>(null);
  const [loading, setLoading] = useState(true);
  const [testing, setTesting] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [webhookUrl, setWebhookUrl] = useState("");
  const [hasProAccess, setHasProAccess] = useState(false);

  useEffect(() => {
    loadIntegration();
    checkSubscription();
  }, []);

  async function checkSubscription() {
    try {
      const res = await fetch("/api/subscription", { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        setHasProAccess(data.tierID === "pro" || data.tierID === "premium");
      }
    } catch (err) {
      console.error("Failed to check subscription:", err);
    }
  }

  async function loadIntegration() {
    setLoading(true);
    try {
      const res = await fetch("/api/n8n/integration", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setIntegration(data);
        setWebhookUrl(data.webhookBaseUrl || "");
      }
    } catch (err) {
      console.error("Failed to load n8n integration:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleEnable() {
    if (!webhookUrl) {
      alert("Bitte Webhook Base URL eingeben");
      return;
    }

    try {
      const res = await fetch("/api/n8n/integration/enable", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ webhook_base_url: webhookUrl }),
      });

      if (res.ok) {
        alert("n8n Integration aktiviert! 🎉");
        loadIntegration();
        setEditMode(false);
      } else {
        const error = await res.text();
        alert(`Fehler: ${error}`);
      }
    } catch (err) {
      console.error("Enable error:", err);
      alert("Aktivierung fehlgeschlagen");
    }
  }

  async function handleDisable() {
    if (!confirm("n8n Integration wirklich deaktivieren?")) {
      return;
    }

    try {
      const res = await fetch("/api/n8n/integration/disable", {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        alert("n8n Integration deaktiviert");
        loadIntegration();
      } else {
        alert("Deaktivierung fehlgeschlagen");
      }
    } catch (err) {
      console.error("Disable error:", err);
      alert("Deaktivierung fehlgeschlagen");
    }
  }

  async function handleTestConnection() {
    setTesting(true);
    try {
      const res = await fetch("/api/n8n/integration/test", {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        alert("✅ Verbindung erfolgreich getestet!");
      } else {
        alert("❌ Verbindungstest fehlgeschlagen");
      }
    } catch (err) {
      console.error("Test error:", err);
      alert("❌ Verbindungstest fehlgeschlagen");
    } finally {
      setTesting(false);
    }
  }

  if (loading) {
    return (
      <div className="p-6 max-w-4xl mx-auto">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-800 rounded w-1/3"></div>
          <div className="h-64 bg-gray-800 rounded"></div>
        </div>
      </div>
    );
  }

  if (!hasProAccess) {
    return (
      <div className="p-6 max-w-4xl mx-auto">
        <div className="bg-gray-800 rounded-xl p-8 text-center">
          <Zap className="w-16 h-16 text-gray-600 mx-auto mb-4" />
          <h2 className="text-2xl font-bold mb-2">n8n Integration</h2>
          <p className="text-gray-400 mb-6">
            Advanced Commands und Workflows mit n8n sind nur in Pro und Premium verfügbar.
          </p>
          <a
            href="/dashboard/subscription"
            className="inline-block px-6 py-3 bg-green-600 hover:bg-green-500 rounded-lg transition-colors font-semibold"
          >
            Jetzt upgraden
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">n8n Integration</h1>
          <p className="text-gray-400">
            Verbinde deine n8n Instanz für Advanced Commands und Workflows
          </p>
        </div>
        <button
          onClick={loadIntegration}
          className="p-2 hover:bg-gray-800 rounded-lg transition-colors"
          title="Aktualisieren"
        >
          <RefreshCw className="w-5 h-5" />
        </button>
      </div>

      {/* Status Card */}
      <div className="bg-gray-800 rounded-xl p-6 mb-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold">Connection Status</h2>
          <div className="flex items-center gap-2">
            {integration?.isReady ? (
              <>
                <CheckCircle className="w-5 h-5 text-green-500" />
                <span className="text-green-500 font-semibold">Connected</span>
              </>
            ) : integration?.enabled ? (
              <>
                <AlertTriangle className="w-5 h-5 text-yellow-500" />
                <span className="text-yellow-500 font-semibold">
                  Enabled, but not ready
                </span>
              </>
            ) : (
              <>
                <XCircle className="w-5 h-5 text-gray-500" />
                <span className="text-gray-500 font-semibold">Disabled</span>
              </>
            )}
          </div>
        </div>

        {/* Webhook URL Config */}
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">
              Webhook Base URL
            </label>
            {editMode || !integration?.webhookBaseUrl ? (
              <div className="flex gap-2">
                <input
                  type="text"
                  value={webhookUrl}
                  onChange={(e) => setWebhookUrl(e.target.value)}
                  placeholder="https://n8n.example.com/webhook"
                  className="flex-1 bg-gray-900 rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <button
                  onClick={handleEnable}
                  className="px-4 py-2 bg-green-600 hover:bg-green-500 rounded transition-colors"
                >
                  Speichern
                </button>
                {integration?.webhookBaseUrl && (
                  <button
                    onClick={() => {
                      setEditMode(false);
                      setWebhookUrl(integration.webhookBaseUrl);
                    }}
                    className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded transition-colors"
                  >
                    Abbrechen
                  </button>
                )}
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <div className="flex-1 bg-gray-900 rounded px-4 py-2 font-mono text-sm">
                  {integration.webhookBaseUrl}
                </div>
                <button
                  onClick={() => setEditMode(true)}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded transition-colors"
                >
                  Bearbeiten
                </button>
              </div>
            )}
          </div>

          {/* Example Webhook URLs */}
          {integration?.webhookBaseUrl && (
            <div className="bg-gray-900 rounded p-4">
              <p className="text-sm text-gray-400 mb-2">Beispiel Webhook URLs:</p>
              <div className="space-y-1 text-sm font-mono">
                <div className="text-gray-500">
                  Commands: {integration.webhookBaseUrl}/command-handler
                </div>
                <div className="text-gray-500">
                  Votes: {integration.webhookBaseUrl}/vote-handler
                </div>
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex gap-2">
            {integration?.enabled && integration?.webhookBaseUrl && (
              <>
                <button
                  onClick={handleTestConnection}
                  disabled={testing}
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 rounded transition-colors flex items-center gap-2"
                >
                  <Link className="w-4 h-4" />
                  {testing ? "Teste..." : "Verbindung testen"}
                </button>
                <button
                  onClick={handleDisable}
                  className="px-4 py-2 bg-red-600 hover:bg-red-500 rounded transition-colors"
                >
                  Deaktivieren
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Usage Stats */}
      {integration?.enabled && (
        <div className="grid md:grid-cols-2 gap-6">
          <div className="bg-gray-800 rounded-xl p-6">
            <h3 className="text-lg font-bold mb-2">Workflows</h3>
            <div className="text-3xl font-bold text-green-500">
              {integration.workflowsUsed}
            </div>
            <p className="text-sm text-gray-400 mt-1">In Verwendung</p>
          </div>

          <div className="bg-gray-800 rounded-xl p-6">
            <h3 className="text-lg font-bold mb-2">Votes</h3>
            <div className="text-3xl font-bold text-blue-500">
              {integration.votesThisMonth}
            </div>
            <p className="text-sm text-gray-400 mt-1">
              Diesen Monat (Reset:{" "}
              {new Date(integration.lastResetAt).toLocaleDateString()})
            </p>
          </div>
        </div>
      )}

      {/* Setup Guide */}
      <div className="mt-6 bg-gray-800 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-4">Setup Guide</h3>
        <ol className="space-y-3 text-sm text-gray-400">
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              1
            </span>
            <div>
              <strong className="text-white">n8n Instanz aufsetzen</strong>
              <p>Installiere n8n auf einem Server (Cloud oder Self-Hosted)</p>
            </div>
          </li>
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              2
            </span>
            <div>
              <strong className="text-white">Workflow erstellen</strong>
              <p>
                Erstelle Workflows mit Webhook-Triggern für Commands oder Votes
              </p>
            </div>
          </li>
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              3
            </span>
            <div>
              <strong className="text-white">Base URL eintragen</strong>
              <p>Trage deine n8n Webhook Base URL oben ein (ohne /webhook am Ende)</p>
            </div>
          </li>
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              4
            </span>
            <div>
              <strong className="text-white">Verbindung testen</strong>
              <p>Klicke auf "Verbindung testen" um sicherzustellen dass alles läuft</p>
            </div>
          </li>
        </ol>
      </div>
    </div>
  );
}