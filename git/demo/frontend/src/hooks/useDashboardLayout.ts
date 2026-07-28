import { useCallback, useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";

export interface WidgetLayoutEntry {
  x: number;
  y: number;
  width: number;
  height?: number;
  visible: boolean;
  zIndex: number;
}

type Layout = Record<string, WidgetLayoutEntry>;

function storageKey(username: string | null): string {
  return `live-dashboard-layout:${username ?? "anon"}`;
}

function loadLayout(username: string | null): Layout {
  try {
    const raw = localStorage.getItem(storageKey(username));
    return raw ? JSON.parse(raw) : {};
  } catch {
    return {};
  }
}

// Manages per-widget position/visibility/stacking order for the Live-Dashboard's
// freely arrangeable widgets - persisted to localStorage per logged-in user
// (confirmed with the user: no cross-device sync needed for this, so no
// backend involved at all).
export function useDashboardLayout(defaults: Record<string, { x: number; y: number; width: number }>) {
  const { username } = useAuth();
  const [layout, setLayout] = useState<Layout>(() => loadLayout(username));

  useEffect(() => {
    setLayout(loadLayout(username));
  }, [username]);

  useEffect(() => {
    try {
      localStorage.setItem(storageKey(username), JSON.stringify(layout));
    } catch {
      // localStorage unavailable (private mode, quota) - layout just won't persist.
    }
  }, [layout, username]);

  const entryOrDefault = useCallback(
    (layoutSnapshot: Layout, id: string): WidgetLayoutEntry => {
      const fallback = defaults[id] || { x: 0, y: 0, width: 340 };
      return (
        layoutSnapshot[id] || {
          x: fallback.x,
          y: fallback.y,
          width: fallback.width,
          visible: true,
          zIndex: 1,
        }
      );
    },
    [defaults]
  );

  const getEntry = useCallback((id: string) => entryOrDefault(layout, id), [layout, entryOrDefault]);

  const updatePosition = useCallback(
    (id: string, x: number, y: number) => {
      setLayout((prev) => ({
        ...prev,
        [id]: { ...entryOrDefault(prev, id), x, y },
      }));
    },
    [entryOrDefault]
  );

  const updateSize = useCallback(
    (id: string, width: number, height: number) => {
      setLayout((prev) => ({
        ...prev,
        [id]: { ...entryOrDefault(prev, id), width, height },
      }));
    },
    [entryOrDefault]
  );

  const toggleVisible = useCallback(
    (id: string) => {
      setLayout((prev) => {
        const current = entryOrDefault(prev, id);
        return { ...prev, [id]: { ...current, visible: !current.visible } };
      });
    },
    [entryOrDefault]
  );

  const bringToFront = useCallback(
    (id: string, allIds: string[]) => {
      setLayout((prev) => {
        const maxZ = allIds.reduce((max, otherId) => {
          const z = prev[otherId]?.zIndex ?? 1;
          return z > max ? z : max;
        }, 1);
        const current = entryOrDefault(prev, id);
        return { ...prev, [id]: { ...current, zIndex: maxZ + 1 } };
      });
    },
    [entryOrDefault]
  );

  const resetLayout = useCallback(() => {
    setLayout({});
  }, []);

  return { getEntry, updatePosition, updateSize, toggleVisible, bringToFront, resetLayout };
}
