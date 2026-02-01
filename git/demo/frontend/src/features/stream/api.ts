import type {
  StreamInfo,
  LiveStream,
  DashboardStats,
  Category,
  UpdateStreamInfoRequest,
} from "./types";

const API_BASE = "/api/dashboard/stream";

export async function fetchStreamInfo(): Promise<StreamInfo> {
  const res = await fetch(`${API_BASE}/info`, {
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("Fehler beim Laden der Stream-Info");
  }

  return res.json();
}

export async function updateStreamInfo(
  data: UpdateStreamInfoRequest
): Promise<StreamInfo> {
  const res = await fetch(`${API_BASE}/info`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({}));
    throw new Error(error.message || "Fehler beim Aktualisieren der Stream-Info");
  }

  return res.json();
}

export async function fetchLiveStatus(): Promise<LiveStream> {
  const res = await fetch(`${API_BASE}/live`, {
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("Fehler beim Laden des Live-Status");
  }

  return res.json();
}

export async function fetchDashboardStats(): Promise<DashboardStats> {
  const res = await fetch(`${API_BASE}/stats`, {
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("Fehler beim Laden der Statistiken");
  }

  return res.json();
}

export async function searchCategories(query: string): Promise<Category[]> {
  const res = await fetch(
    `${API_BASE}/categories/search?query=${encodeURIComponent(query)}`,
    {
      credentials: "include",
    }
  );

  if (!res.ok) {
    throw new Error("Fehler beim Suchen der Kategorien");
  }

  return res.json();
}