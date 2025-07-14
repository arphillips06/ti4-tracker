// imperialRiderHandler.js
import API_BASE_URL from "../config";
export async function handleScoreImperialRider(playerId, gameId, refreshGameState, setShowImperialModal) {
  try {
    const res = await fetch(`${API_BASE_URL}/score/imperial`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        player_id: parseInt(playerId),
      }),
    });
    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Failed to score Imperial Rider: ${errorText}`);
    }
    await refreshGameState();
    setShowImperialModal(false);
  } catch (err) {
    console.error(err);
    alert("Failed to score Imperial Rider. See console for details.");
  }
}
