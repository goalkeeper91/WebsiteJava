import { apiFetch } from "../../lib/apiFetch";
import type {
  StreamInfo,
  LiveStream,
  DashboardStats,
  Category,
  UpdateStreamInfoRequest,
  CommercialRequest,
  CommercialResult,
} from "./types";

const API_BASE = "/api/dashboard/stream";

export async function fetchStreamInfo(): Promise<StreamInfo> {
  const res = await apiFetch(`${API_BASE}/info`);

  if (!res.ok) {
    throw new Error("Fehler beim Laden der Stream-Info");
  }

  return res.json();
}

export async function updateStreamInfo(
  data: UpdateStreamInfoRequest
): Promise<StreamInfo> {
  const res = await apiFetch(`${API_BASE}/info`, {
    method: "PATCH",
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
  const res = await apiFetch(`${API_BASE}/live`);

  if (!res.ok) {
    throw new Error("Fehler beim Laden des Live-Status");
  }

  return res.json();
}

export async function fetchDashboardStats(): Promise<DashboardStats> {
  const res = await apiFetch(`${API_BASE}/stats`);

  if (!res.ok) {
    throw new Error("Fehler beim Laden der Statistiken");
  }

  return res.json();
}

export async function searchCategories(query: string): Promise<Category[]> {
  const res = await apiFetch(
    `${API_BASE}/categories/search?query=${encodeURIComponent(query)}`
  );

  if (!res.ok) {
    throw new Error("Fehler beim Suchen der Kategorien");
  }

  return res.json();
}

export async function startCommercial(data: CommercialRequest): Promise<CommercialResult> {
  const res = await apiFetch(`${API_BASE}/commercial`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    throw new Error("Konnte keine Werbepause starten");
  }

  return res.json();
}
