import API_BASE_URL from "../config";

export async function handleScoreCrown(playerId, gameId, refresh) {
  console.log("Sending Crown payload:", {
    player_id: parseInt(playerId),
    game_id: parseInt(gameId),
  });
  const res = await fetch(`${API_BASE_URL}/relic/crown`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      player_id: parseInt(playerId), // ✅ key name must be player_id
      game_id: parseInt(gameId),     // ✅ key name must be game_id
    }),
  });

  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Failed to score Crown of Emphidia");
  }

  refresh();
}
