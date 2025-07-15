// src/utils/relicHandlers.js
import { postJSON, submitAndRefresh } from "./helpers";

export async function handleShardSubmit(playerId, gameId, refreshGameState, closeModal) {
  try {
    await submitAndRefresh({
      requestFn: () =>
        postJSON("/relic/shard", {
          game_id: parseInt(gameId),
          new_holder_id: parseInt(playerId),
        }),
      refreshGameState,
      closeModal,
    });
  } catch (err) {
    console.error("Failed to submit Shard:", err);
    alert("Error assigning Shard of the Throne.");
  }
}

export async function handleScoreCrown(playerId, gameId, refresh) {
  try {
    await submitAndRefresh({
      requestFn: () =>
        postJSON("/relic/crown", {
          game_id: parseInt(gameId),
          player_id: parseInt(playerId),
        }),
      refreshGameState: refresh,
    });
  } catch (err) {
    console.error("Failed to submit Crown of Emphidia:", err);
    alert("Error scoring Crown of Emphidia.");
  }
}
