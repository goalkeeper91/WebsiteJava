import type { ChatCommand, PageResponse } from "./types";

export async function fetchCommands(
  page = 0,
  size = 10
): Promise<PageResponse<ChatCommand>> {
  const res = await fetch(
    `/api/dashboard/commands?page=${page}&size=${size}`,
    { credentials: "include" }
  );

  if (!res.ok) {
    throw new Error("Fehler beim Laden der Commands");
  }

  return res.json();
}

export async function createCommand(payload: {
  trigger: string;
  response: string;
  cooldown: number;
}) {
  const res = await fetch(`/api/dashboard/commands`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) throw new Error("Command konnte nicht erstellt werden");
  return res.json();
}

export async function updateCommand(
  id: number | string,
  payload: {
    trigger?: string;
    response?: string;
    cooldown?: number;
    enabled?: boolean;
  }
) {
  const res = await fetch(`/api/dashboard/commands/${id}`, {
    method: "PUT",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) throw new Error("Command konnte nicht aktualisiert werden");
  return res.json();
}

export async function deleteCommand(id: number | string) {
  const res = await fetch(`/api/dashboard/commands/${id}`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!res.ok) throw new Error("Command konnte nicht gelöscht werden");
}
