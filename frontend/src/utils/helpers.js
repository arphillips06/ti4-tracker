import API_BASE_URL from "../config";
import { isAgendaScore } from "./selectors";

export async function postJSON(url, payload) {
  const finalUrl = url.startsWith("http")
    ? url
    : `${API_BASE_URL}${url}`;

  const res = await fetch(finalUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) throw new Error(`POST ${finalUrl} failed: ${res.status}`);
  return res.json();
}

export async function getJSON(url) {
  const finalUrl = url.startsWith("http")
    ? url
    : `${API_BASE_URL}${url}`;

  const res = await fetch(finalUrl, { headers: { "Content-Type": "application/json" } });
  if (!res.ok) throw new Error(`GET ${finalUrl} failed: ${res.status}`);
  return await res.json();
}

export async function submitAndRefresh({ requestFn, refreshGameState, setGame, closeModal }) {
  try {
    const result = await requestFn();
    if (setGame && result) setGame(result);
    if (refreshGameState) await refreshGameState();
    if (closeModal) closeModal();
  } catch (err) {
    console.error("❌ Request failed:", err);
    alert(err.message || "Something went wrong.");
  }
}

export const isAgendaUsed = (game, title) =>
  game?.AllScores?.some((s) => isAgendaScore(s, title));
