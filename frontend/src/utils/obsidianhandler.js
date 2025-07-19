// src/utils/obsidianhandler.js
import { postJSON, submitAndRefresh } from "./helpers";

export async function handleAssignObsidian(playerId, gameId, refreshGameState, onClose) {
  try {
    await submitAndRefresh({
      requestFn: () =>
        postJSON("/relic/obsidian", {
          game_id: parseInt(gameId),
          player_id: parseInt(playerId),
        }),
      refreshGameState,
      closeModal: onClose,
    });
  } catch (err) {
    console.error("Failed to assign Obsidian:", err);
    alert("Failed to assign The Obsidian. See console.");
  }
}
