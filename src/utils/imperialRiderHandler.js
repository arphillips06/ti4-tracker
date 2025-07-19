// src/utils/imperialRiderHandler.js
import { postJSON, submitAndRefresh } from "./helpers";

export async function handleScoreImperialRider(playerId, gameId, refreshGameState, setShowImperialModal) {
  try {
    await submitAndRefresh({
      requestFn: () =>
        postJSON("/score/imperial", {
          game_id: parseInt(gameId),
          player_id: parseInt(playerId),
        }),
      refreshGameState,
      closeModal: () => setShowImperialModal(false),
    });
  } catch (err) {
    console.error(err);
    alert("Failed to score Imperial Rider. See console for details.");
  }
}
