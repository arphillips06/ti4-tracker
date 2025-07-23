import API_BASE_URL from "../config";

/**
 * Perform a POST request with JSON payload and basic error handling.
 * @param {string} path - Relative API path, e.g., "/agenda/mutiny"
 * @param {Object} body - Payload to send
 * @returns {Promise<any>} - Parsed JSON or throws error
 */
export async function postJSON(path, body) {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    const errorText = await res.text().catch(() => "Unknown error");
    throw new Error(`POST ${path} failed: ${errorText}`);
  }

  try {
    return await res.json();
  } catch {
    return {};
  }
}

/**
 * Submit a request and optionally refresh game state or set modal/game state.
 * @param {Function} requestFn - async function that sends the request
 * @param {Function} [refreshGameState]
 * @param {Function} [setGame]
 * @param {Function} [closeModal]
 */
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

/**
 * Utility to check if an agenda has already been used in a game.
 * @param {Object} game - The full game object
 * @param {string} title - The agenda title to check
 * @returns {boolean} - True if used, false otherwise
 */
export const isAgendaUsed = (game, title) =>
  game?.AllScores?.some(
    (s) => s.Type?.toLowerCase() === "agenda" && s.AgendaTitle === title
  );
