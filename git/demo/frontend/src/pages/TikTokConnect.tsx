import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { CheckCircle2, Loader2, AlertTriangle } from "lucide-react";
import Seo from "../components/Seo";

// Unlisted page (noindex, not linked from the navigation) where a visitor can
// link a TikTok account through Composio's OAuth flow: TikTok Login Kit
// consent -> back here -> the connected account's display name + avatar
// (user.info.basic). The same connection is what lets a video be sent to the
// account's TikTok inbox as a draft (video.upload) - that step happens in
// Claude/Composio, not on this page.
//
// All calls go to our own /api/tiktok/* endpoints (backend-go): the Composio
// API key never reaches the browser, and nothing about the visitor is stored
// in our database.

type Profile = { displayName: string; avatarUrl: string };
type Status =
  | { kind: "loading" }
  | { kind: "disconnected" }
  | { kind: "connected"; profile: Profile }
  | { kind: "error"; message: string };

const GENERIC_ERROR = "Something went wrong. Please try again later.";

async function readError(res: Response): Promise<string> {
  try {
    const body = await res.json();
    if (body && typeof body.error === "string") return body.error;
  } catch {
    /* fall through */
  }
  return GENERIC_ERROR;
}

export default function TikTokConnect() {
  const [status, setStatus] = useState<Status>({ kind: "loading" });
  const [busy, setBusy] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  // Composio appends ?status=success|failed&connected_account_id=... when it
  // sends the visitor back. Only used as a hint to word a failure - the real
  // state always comes from /api/tiktok/me, never from the query string.
  const returnedFailed = searchParams.get("status") === "failed";

  const loadStatus = useCallback(async () => {
    try {
      const res = await fetch("/api/tiktok/me", { credentials: "include" });
      if (!res.ok) {
        setStatus({ kind: "error", message: await readError(res) });
        return;
      }
      const body = await res.json();
      if (body.connected) {
        setStatus({
          kind: "connected",
          profile: { displayName: body.displayName ?? "", avatarUrl: body.avatarUrl ?? "" },
        });
      } else {
        setStatus({ kind: "disconnected" });
      }
    } catch {
      setStatus({ kind: "error", message: GENERIC_ERROR });
    }
  }, []);

  useEffect(() => {
    void loadStatus();
  }, [loadStatus]);

  // Drop the callback query params from the address bar once handled.
  useEffect(() => {
    if (searchParams.has("status") || searchParams.has("connected_account_id")) {
      navigate("/tiktok", { replace: true });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const connect = async () => {
    setBusy(true);
    setActionError(null);
    try {
      const res = await fetch("/api/tiktok/connect", { method: "POST", credentials: "include" });
      if (!res.ok) {
        setActionError(await readError(res));
        setBusy(false);
        return;
      }
      const body = await res.json();
      if (typeof body.redirectUrl === "string" && body.redirectUrl.startsWith("https://")) {
        window.location.href = body.redirectUrl;
        return; // leaving the page - keep the button disabled
      }
      setActionError(GENERIC_ERROR);
    } catch {
      setActionError(GENERIC_ERROR);
    }
    setBusy(false);
  };

  const disconnect = async () => {
    setBusy(true);
    setActionError(null);
    try {
      const res = await fetch("/api/tiktok/disconnect", { method: "POST", credentials: "include" });
      if (!res.ok) {
        setActionError(await readError(res));
      } else {
        setStatus({ kind: "disconnected" });
      }
    } catch {
      setActionError(GENERIC_ERROR);
    }
    setBusy(false);
  };

  return (
    <section className="w-full min-h-screen bg-slate-950 flex items-center justify-center py-10 px-4">
      <Seo
        title="Connect TikTok"
        description="Connect your TikTok account."
        path="/tiktok"
        noindex
      />
      <div className="w-full max-w-xl bg-slate-900 text-white p-8 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold mb-3">Connect your TikTok account</h1>
        <p className="text-gray-300 mb-6">
          Link your TikTok account through TikTok Login Kit. We only see what you approve on TikTok's
          consent screen, and you can disconnect again at any time.
        </p>

        <h2 className="text-lg font-semibold mb-2">Permissions requested</h2>
        <ul className="list-disc list-inside text-gray-300 mb-6 space-y-1">
          <li>
            <code className="text-goalyBlue">user.info.basic</code> &ndash; read your display name and profile
            picture, shown below once connected.
          </li>
          <li>
            <code className="text-goalyBlue">video.upload</code> &ndash; send a video to your TikTok inbox as a
            draft. Nothing is published: you review and post it yourself in the TikTok app.
          </li>
        </ul>

        {status.kind === "loading" && (
          <div className="flex items-center gap-3 text-gray-300" role="status">
            <Loader2 className="w-5 h-5 animate-spin" aria-hidden /> Checking connection&hellip;
          </div>
        )}

        {status.kind === "error" && (
          <div className="flex items-start gap-3 rounded-md bg-red-950/60 border border-red-800 p-4 mb-4" role="alert">
            <AlertTriangle className="w-5 h-5 mt-0.5 text-red-400 shrink-0" aria-hidden />
            <div>
              <p className="text-red-200">{status.message}</p>
              <button
                type="button"
                onClick={() => {
                  setStatus({ kind: "loading" });
                  void loadStatus();
                }}
                className="mt-2 underline text-sm text-red-200"
              >
                Retry
              </button>
            </div>
          </div>
        )}

        {status.kind === "disconnected" && (
          <>
            {returnedFailed && (
              <div className="flex items-start gap-3 rounded-md bg-red-950/60 border border-red-800 p-4 mb-4" role="alert">
                <AlertTriangle className="w-5 h-5 mt-0.5 text-red-400 shrink-0" aria-hidden />
                <p className="text-red-200">The connection was not completed. You can try again below.</p>
              </div>
            )}
            <button
              type="button"
              onClick={connect}
              disabled={busy}
              className="w-full py-3 px-6 bg-white text-black hover:bg-gray-200 disabled:opacity-60 disabled:cursor-not-allowed rounded-lg transition-colors font-semibold"
            >
              {busy ? "Redirecting to TikTok…" : "Connect with TikTok"}
            </button>
          </>
        )}

        {status.kind === "connected" && (
          <div>
            <div className="flex items-center gap-4 rounded-lg bg-slate-800 p-4 mb-4">
              {status.profile.avatarUrl ? (
                <img
                  src={status.profile.avatarUrl}
                  alt=""
                  referrerPolicy="no-referrer"
                  className="w-16 h-16 rounded-full object-cover bg-slate-700"
                />
              ) : (
                <div className="w-16 h-16 rounded-full bg-slate-700" aria-hidden />
              )}
              <div className="min-w-0">
                <p className="flex items-center gap-2 text-green-400 text-sm font-semibold">
                  <CheckCircle2 className="w-4 h-4" aria-hidden /> Connected
                </p>
                <p className="text-xl font-bold break-words">{status.profile.displayName || "TikTok user"}</p>
              </div>
            </div>
            <button
              type="button"
              onClick={disconnect}
              disabled={busy}
              className="py-2 px-4 border border-gray-600 hover:bg-slate-800 disabled:opacity-60 disabled:cursor-not-allowed rounded-lg transition-colors text-sm"
            >
              {busy ? "Disconnecting…" : "Disconnect"}
            </button>
          </div>
        )}

        {actionError && (
          <p className="mt-4 text-red-300 text-sm" role="alert">
            {actionError}
          </p>
        )}

        <p className="text-sm text-gray-400 mt-8">
          We don't store your TikTok data in our own database. Your authorization is held by our integration
          provider Composio until you disconnect, and a temporary cookie (24&nbsp;hours) links this browser to
          it. Details in our{" "}
          <Link to="/legal/datenschutz" className="underline text-goalyBlue">
            privacy policy
          </Link>{" "}
          (German) and{" "}
          <Link to="/legal/agb" className="underline text-goalyBlue">
            terms
          </Link>
          .
        </p>
      </div>
    </section>
  );
}
