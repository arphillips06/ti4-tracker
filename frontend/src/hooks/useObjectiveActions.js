import { useCallback } from "react";
import API_BASE_URL from "../config"

export default function useObjectiveActions(gameId, refreshGameState, setLocalScored) {
  const scoreObjective = useCallback(
    async (playerId, objectiveId) => {
      const payload = {
        game_id: parseInt(gameId),
        player_id: playerId,
        objective_id: objectiveId,
      };

      try {
        const res = await fetch(`${API_BASE_URL}/score`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });

        if (!res.ok) {
          const errorText = await res.text();
          console.error("Scoring rejected:", errorText);
          alert("Scoring rejected by backend: " + errorText);
          return false;
        }

        setLocalScored((prev) => {
          const current = prev[objectiveId] || [];
          return {
            ...prev,
            [objectiveId]: [...new Set([...current, playerId])],
          };
        });

        await refreshGameState();
        return true;
      } catch (err) {
        console.error("Scoring failed:", err);
        alert("Scoring failed. See console.");
        return false;
      }
    },
    [gameId, refreshGameState, setLocalScored]
  );

  const unscoreObjective = useCallback(
    async (playerId, objectiveId) => {
      try {
        const res = await fetch(`${API_BASE_URL}/unscore`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            game_id: parseInt(gameId),
            player_id: playerId,
            objective_id: objectiveId,
          }),
        });

        if (!res.ok) {
          const errorText = await res.text();
          console.error("Unscoring failed:", errorText);
          alert("Unscoring failed: " + errorText);
          return false;
        }

        await refreshGameState();
        return true;
      } catch (err) {
        console.error("Unscoring error:", err);
        alert("Unscoring failed. See console.");
        return false;
      }
    },
    [gameId, refreshGameState]
  );

  const assignObjective = useCallback(
    async (roundId, objectiveId) => {
      try {
        const res = await fetch(`${API_BASE_URL}/assign_objective`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            game_id: parseInt(gameId),
            round_id: roundId,
            objective_id: objectiveId,
          }),
        });

        if (!res.ok) {
          const errText = await res.text();
          throw new Error(errText);
        }

        await refreshGameState();
      } catch (err) {
        console.error("Failed to assign objective:", err);
        alert("Failed to assign objective. See console for details.");
      }
    },
    [gameId, refreshGameState]
  );

  const advanceRound = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/games/${gameId}/advance-round`, {
        method: "POST",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText);
      }

      await refreshGameState();
    } catch (err) {
      console.error("Failed to advance round:", err);
      alert("Could not advance round. See console for details.");
    }
  }, [gameId, refreshGameState]);

  return {
    scoreObjective,
    unscoreObjective,
    assignObjective,
    advanceRound,
  };
}
