import { useState, useEffect } from "react";
import { Bot, CheckCircle, AlertTriangle, RefreshCcw } from "lucide-react";

interface BotAuthStatus {
  authenticated: boolean;
  botUsername?: string;
  botTwitchId?: string;
  tokenPresent?: boolean;
  expiresAt?: string;
}

export default function BotAuthCard() {
  const [status, setStatus] = useState<BotAuthStatus | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    checkBotStatus();
  }, []);

  async function checkBotStatus() {
    setLoading(true);
    try {
      // TODO: Richtigen API Endpoint implementieren
      const res = await fetch("/api/bot/auth/status", {
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();
        setStatus(data);
      }
    } catch (err) {
      console.error("Fehler beim Laden des Bot Auth Status:", err);
      // Fallback für Development
      setStatus({
        authenticated: false,
        tokenPresent: false,
      });
    } finally {
      setLoading(false);
    }
  }

  function handleBotAuth() {
    // Redirect zum Bot Auth Flow
    window.location.href = "/auth/login/bot";
  }

  const renderTokenProgress = () => {
    if (!status?.expiresAt) return <p className="text-sm text-gray-500">Kein Token vorhanden</p>;

    const expiryDate = new Date(status.expiresAt);
    const now = new Date();
    const totalMs = expiryDate.getTime() - now.getTime();
    const daysLeft = totalMs / (1000 * 60 * 60 * 24);
    const percent = Math.max(0, Math.min(100, (daysLeft / 60) * 100)); // 60 Tage max

    return (
      <div>
        <div className="w-full rounded-full h-3 bg-gray-700">
          <div
            className={`h-3 rounded-full transition-all ${
              percent > 30 ? "bg-green-500" : "bg-red-500"
            }`}
            style={{ width: `${percent}%` }}
          />
        </div>
        <p className="text-sm text-gray-400 mt-1">
          Token läuft ab am {expiryDate.toLocaleString()} ({daysLeft.toFixed(1)} Tage)
        </p>
      </div>
    );
  };

  if (loading) {
    return (
      <div className="p-4 rounded-xl shadow animate-pulse">
        <div className="h-6 bg-gray-700 rounded w-1/3 mb-3"></div>
        <div className="h-4 bg-gray-700 rounded w-2/3"></div>
      </div>
    );
  }

  return (
    <div className="p-4 rounded-xl shadow">
      {/* Header mit Status */}
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-bold flex items-center gap-2">
          <Bot className="w-5 h-5" />
          Bot Account
        </h3>
        <button
          onClick={checkBotStatus}
          className="p-1 hover:bg-gray-700 rounded transition-colors"
          title="Status aktualisieren"
        >
          <RefreshCcw className="w-4 h-4" />
        </button>
      </div>

      {/* Status Anzeige */}
      <div className="text-sm mb-3">
        {status?.authenticated ? (
          <span className="flex items-center text-green-600">
            <CheckCircle className="w-5 h-5 mr-1" />
            Authentifiziert als {status.botUsername}
          </span>
        ) : (
          <span className="flex items-center text-red-600">
            <AlertTriangle className="w-5 h-5 mr-1" />
            Nicht authentifiziert
          </span>
        )}
      </div>

      {/* Token Progress (wenn authentifiziert) */}
      {status?.authenticated && status?.tokenPresent ? (
        <div className="mb-3">
          {renderTokenProgress()}
        </div>
      ) : null}

      {/* Action Button */}
      {status?.authenticated ? (
        <div className="space-y-2">
          {/* Bot Info */}
          {status.botTwitchId && (
            <p className="text-xs text-gray-400">
              Bot ID: {status.botTwitchId}
            </p>
          )}

          {/* Re-Auth Button */}
          <button
            onClick={handleBotAuth}
            className="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded transition-colors"
          >
            Token erneuern
          </button>
        </div>
      ) : (
        <button
          onClick={handleBotAuth}
          className="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded transition-colors"
        >
          Bot authentifizieren
        </button>
      )}
    </div>
  );
}