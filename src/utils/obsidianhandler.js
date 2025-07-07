import API_BASE_URL from "../config";

export const handleScoreObsidian = async (playerId, gameId, refresh, onClose) => {
  try {
    await fetch(`${API_BASE_URL}/relic/obsidian`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ player_id: playerId, game_id: gameId }),
    });
    refresh();
    onClose();
  } catch (err) {
    console.error("Failed to score Obsidian:", err);
  }
};
