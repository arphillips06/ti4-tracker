// src/hooks/useGroupedScoredSecrets.js
import { useMemo } from "react";

export default function useGroupedScoredSecrets(allScores, secretObjectives) {
  return useMemo(() => {
    const result = {};

    (allScores || []).forEach((s) => {
      if ((s.Type || "").toLowerCase() !== "secret") return;

      const scoreId = parseInt(s.ObjectiveID);
      const match = secretObjectives.find((o) => {
        const oid = parseInt(o.id ?? o.ID);
        return oid === scoreId;
      });

      if (match) {
        if (!result[s.PlayerID]) {
          result[s.PlayerID] = [];
        }
        result[s.PlayerID].push(match);
      }
    });

    return result;
  }, [allScores, secretObjectives]);
}
