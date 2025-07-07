// src/utils/relicHandlers.js
import API_BASE_URL from "../config";
export async function handleShardSubmit(playerId, gameId, refreshGameState, closeModal) {
  try {
    const res = await fetch(`${API_BASE_URL}/relic/shard`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        new_holder_id: parseInt(playerId),
      }),
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Failed to assign Shard: ${errorText}`);
    }

    await refreshGameState();
    closeModal();
  } catch (err) {
    console.error("Failed to submit Shard:", err);
    alert("Error assigning Shard of the Throne.");
  }
}
