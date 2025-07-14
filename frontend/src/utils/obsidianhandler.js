import API_BASE_URL from "../config";

export async function handleAssignObsidian(playerId, gameId, refreshGameState, onClose) {
  try {
    await fetch(`${API_BASE_URL}/relic/obsidian`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        player_id: parseInt(playerId),
      }),
    });
    await refreshGameState();

    if (onClose) onClose();
  } catch (err) {
    console.error("Failed to assign Obsidian:", err);
    alert("Failed to assign The Obsidian. See console.");
  }
}