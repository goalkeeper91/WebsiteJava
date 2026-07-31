import { useEffect, useState } from "react";
import { Zap, CheckCircle, XCircle, AlertTriangle, RefreshCw, Server } from "lucide-react";

interface N8NIntegration {
  id: number;
  enabled: boolean;
  workflowsUsed: number;
  votesThisMonth: number;
  lastResetAt: string;
  isReady: boolean;
}

export default function N8NIntegrationSetup() {
  const [integration, setIntegration] = useState<N8NIntegration | null>(null);
  const [loading, setLoading] = useState(true);
  const [hasProAccess, setHasProAccess] = useState(false);
  const [enabling, setEnabling] = useState(false);

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
      }
    } catch (err) {
      console.error("Failed to load n8n integration:", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleToggle() {
    const action = integration?.enabled ? "disable" : "enable";

    if (action === "disable" && !confirm("n8n Integration wirklich deaktivieren?")) {
      return;
    }

    setEnabling(true);
    try {
      const res = await fetch(`/api/n8n/integration/${action}`, {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        alert(action === "enable"
          ? "n8n Integration aktiviert! 🎉 Du kannst jetzt Advanced Commands erstellen."
          : "n8n Integration deaktiviert"
        );
        loadIntegration();
      } else {
        const error = await res.text();
        alert(`Fehler: ${error}`);
      }
    } catch (err) {
      console.error("Toggle error:", err);
      alert("Aktion fehlgeschlagen");
    } finally {
      setEnabling(false);
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 text-white">
        <div className="p-6 max-w-4xl mx-auto">
          <div className="animate-pulse space-y-4">
            <div className="h-8 bg-gray-800 rounded w-1/3"></div>
            <div className="h-64 bg-gray-800 rounded"></div>
          </div>
        </div>
      </div>
    );
  }

  if (!hasProAccess) {
    return (
      <div className="min-h-screen bg-gray-900 text-white">
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
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white">
    <div className="p-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">n8n Integration</h1>
          <p className="text-gray-400">
            Verwende Advanced Commands und Workflows powered by n8n
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
        <div className="flex flex-wrap items-center justify-between gap-2 mb-6">
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
                  Enabled, initializing...
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

        {/* Managed Infrastructure Info */}
        <div className="bg-gradient-to-r from-blue-900/20 to-purple-900/20 border border-blue-500/30 rounded-lg p-4 mb-6">
          <div className="flex items-start gap-3">
            <Server className="w-6 h-6 text-blue-400 flex-shrink-0 mt-0.5" />
            <div>
              <h3 className="font-bold text-blue-400 mb-1">Managed n8n Infrastructure</h3>
              <p className="text-sm text-gray-400">
                Deine Workflows laufen auf unserer managed n8n Instanz.
                Kein Setup nötig – aktiviere die Integration und los geht's!
              </p>
            </div>
          </div>
        </div>

        {/* Toggle Button */}
        <div className="flex justify-center">
          <button
            onClick={handleToggle}
            disabled={enabling}
            className={`px-8 py-3 rounded-lg font-semibold transition-all ${
              integration?.enabled
                ? "bg-red-600 hover:bg-red-500"
                : "bg-green-600 hover:bg-green-500"
            } disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            {enabling
              ? "Bitte warten..."
              : integration?.enabled
              ? "Integration deaktivieren"
              : "Integration aktivieren"}
          </button>
        </div>
      </div>

      {/* Usage Stats */}
      {integration?.enabled && (
        <div className="grid md:grid-cols-2 gap-6 mb-6">
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

      {/* Getting Started Guide */}
      <div className="bg-gray-800 rounded-xl p-6">
        <h3 className="text-lg font-bold mb-4">So funktioniert's</h3>
        <ol className="space-y-3 text-sm text-gray-400">
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              1
            </span>
            <div>
              <strong className="text-white">Integration aktivieren</strong>
              <p>Klicke oben auf "Integration aktivieren" - fertig!</p>
            </div>
          </li>
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              2
            </span>
            <div>
              <strong className="text-white">Workflow Template wählen</strong>
              <p>
                Wähle im n8n-Editor ein passendes Template aus (z.B. Hero Voting)
              </p>
            </div>
          </li>
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              3
            </span>
            <div>
              <strong className="text-white">Advanced Command erstellen</strong>
              <p>
                Erstelle einen Advanced Command unter <a href="/dashboard" className="text-blue-400 hover:underline">Chat Commands</a> und verknüpfe ihn mit einem Workflow
              </p>
            </div>
          </li>
          <li className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-600 text-white flex items-center justify-center text-xs font-bold">
              4
            </span>
            <div>
              <strong className="text-white">Testen und nutzen!</strong>
              <p>Teste deinen Command im Chat - die n8n Workflows laufen automatisch im Hintergrund</p>
            </div>
          </li>
        </ol>
      </div>

      {/* Feature Highlights */}
      {integration?.enabled && (
        <div className="mt-6 bg-gradient-to-r from-green-900/20 to-blue-900/20 border border-green-500/30 rounded-xl p-6">
          <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
            <Zap className="w-5 h-5 text-green-400" />
            Was du jetzt tun kannst:
          </h3>
          <div className="grid md:grid-cols-2 gap-4 text-sm">
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 flex-shrink-0 mt-0.5" />
              <div>
                <strong className="text-white">Advanced Commands</strong>
                <p className="text-gray-400">Commands mit komplexer Logik und API-Calls</p>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 flex-shrink-0 mt-0.5" />
              <div>
                <strong className="text-white">Vote Sessions</strong>
                <p className="text-gray-400">Interaktive Abstimmungen für deinen Chat</p>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 flex-shrink-0 mt-0.5" />
              <div>
                <strong className="text-white">Custom Workflows</strong>
                <p className="text-gray-400">Verknüpfe mit externen APIs (Discord, Spotify, etc.)</p>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-400 flex-shrink-0 mt-0.5" />
              <div>
                <strong className="text-white">Giveaway Management</strong>
                <p className="text-gray-400">Automatisierte Giveaways mit Teilnahme-Tracking</p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
    </div>
  );
}