import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";

interface OverlayLogEntry {
  time: string;
  text: string;
}

interface OverlayState {
  timeRemaining: number;
  isRunning: boolean;
  eventLog: OverlayLogEntry[];
}

function formatTime(totalSeconds: number): string {
  const clamped = Math.max(0, Math.floor(totalSeconds));
  const h = Math.floor(clamped / 3600).toString().padStart(2, "0");
  const m = Math.floor((clamped % 3600) / 60).toString().padStart(2, "0");
  const s = Math.floor(clamped % 60).toString().padStart(2, "0");
  return `${h}:${m}:${s}`;
}

// Bare, transparent page meant to be added as an OBS Browser Source - no
// site chrome, no auth (OBS can't carry a login cookie), identified purely
// by the ?user= Twitch ID in the URL. Lives at the top level of the router
// (see App.tsx's SiteLayout split) specifically so it never gets wrapped in
// the site's Header/Footer.
export default function SubathonOverlay() {
  const [searchParams] = useSearchParams();
  const userId = searchParams.get("user") ?? "";
  const [state, setState] = useState<OverlayState | null>(null);

  const textColor = `#${searchParams.get("color") || "ffffff"}`;
  const showBg = searchParams.get("bg") !== "false";
  const bgOpacity = searchParams.get("op") || "0.4";
  const showGlow = searchParams.get("glow") !== "false";

  useEffect(() => {
    const previousBg = document.body.style.background;
    document.body.style.setProperty("background", "transparent", "important");
    document.body.style.margin = "0";
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.background = previousBg;
      document.body.style.removeProperty("margin");
      document.body.style.removeProperty("overflow");
    };
  }, []);

  useEffect(() => {
    if (!userId) return;
    let cancelled = false;

    const fetchState = async () => {
      try {
        const res = await fetch(`/api/subathon/overlay/${userId}`);
        if (!res.ok) return;
        const data = await res.json();
        if (!cancelled) setState(data);
      } catch (err) {
        console.error("Overlay sync error:", err);
      }
    };

    fetchState();
    const interval = setInterval(fetchState, 1000);
    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, [userId]);

  const time = state ? formatTime(state.timeRemaining) : "00:00:00";
  const lastEvent = state?.eventLog?.[0] ?? null;

  if (!userId) {
    return (
      <div style={{ color: "#f87171", fontFamily: "monospace", padding: "1rem" }}>
        Missing ?user= parameter
      </div>
    );
  }

  return (
    <div
      style={{
        width: "100vw",
        height: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        fontFamily: "'Courier New', Courier, monospace",
        background: "transparent",
      }}
    >
      <div
        style={{
          textAlign: "center",
          padding: "20px",
          borderLeft: `4px solid ${textColor}`,
          background: showBg ? `rgba(0, 0, 0, ${bgOpacity})` : "transparent",
          backdropFilter: showBg ? "blur(5px)" : "none",
        }}
      >
        <div
          style={{
            fontSize: "0.8rem",
            color: textColor,
            opacity: 0.7,
            letterSpacing: "2px",
          }}
        >
          {state?.isRunning ? "SYSTEM ONLINE" : "SYSTEM PAUSED"}
        </div>

        <div
          style={{
            fontSize: "8rem",
            fontWeight: "bold",
            color: textColor,
            textShadow: showGlow ? `0 0 20px ${textColor}, 0 0 40px ${textColor}` : "none",
          }}
        >
          {time}
        </div>

        {lastEvent && (
          <div
            style={{
              marginTop: "10px",
              borderRight: `2px solid ${textColor}`,
              color: textColor,
              padding: "5px 15px",
              background: "rgba(0, 0, 0, 0.2)",
            }}
          >
            {lastEvent.text}
          </div>
        )}
      </div>
    </div>
  );
}
