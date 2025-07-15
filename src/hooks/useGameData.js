import { useEffect, useState } from "react";
import API_BASE_URL from "../config"
export default function useGameData(gameId) {
  const [game, setGame] = useState(null);
  const [objectives, setObjectives] = useState([]);
  const [objectiveScores, setObjectiveScores] = useState({});
  const [secretObjectives, setSecretObjectives] = useState([]);
  const [secretCounts, setSecretCounts] = useState({});
  const [mutinyUsed, setMutinyUsed] = useState(false);
  const [censureHolder, setCensureHolder] = useState(null);
  const [cdlUsed, setCdlUsed] = useState(false);
  const [crownUsed, setCrownUsed] = useState(false);
const [obsidianHolderId, setObsidianHolderId] = useState(null);



  const fetchGame = async () => {
    const res = await fetch(`${API_BASE_URL}/games/${gameId}`);
    const data = await res.json();
    return data; // Support both wrapped and direct response
  };

  const fetchObjectives = async () => {
    const res = await fetch(`${API_BASE_URL}/games/${gameId}/objectives`);
    return res.json();
  };

  const fetchScores = async () => {
    const res = await fetch(`${API_BASE_URL}/games/${gameId}/objectives/scores`);
    return res.json();
  };

  const fetchSecrets = async () => {
    const res = await fetch(`${API_BASE_URL}/objectives/secrets/all`);
    return res.json();
  };

const refreshGameState = async () => {
  const [gameData, objectivesData, scoresData] = await Promise.all([
    fetchGame(),
    fetchObjectives(),
    fetchScores(),
  ]);

  gameData.AllScores = gameData.all_scores || [];
  gameData.game_players = gameData.players || [];
  gameData.winner_id = gameData.winner_id ?? gameData.WinnerID;

  setGame(gameData);

  const obsidianScore = gameData.AllScores?.find(
    (s) => s.Type === "relic" && s.RelicTitle === "The Obsidian"
  );
  setObsidianHolderId(obsidianScore?.PlayerID || null);

  const map = {};
  (Array.isArray(scoresData) ? scoresData : scoresData?.value || []).forEach((entry) => {
    map[entry.objective_id] = entry.scored_by || [];
  });
  setObjectiveScores(map);

  setObjectives(Array.isArray(objectivesData) ? objectivesData : objectivesData?.value || []);
};
  useEffect(() => {
    (async () => {
      const [gameData, objectiveData, secretData, scoresData] = await Promise.all([
        fetchGame(),
        fetchObjectives(),
        fetchSecrets(),
        fetchScores(),
      ]);


      // 🛠 Fix: Normalize AllScores for compatibility
      gameData.AllScores = gameData.all_scores || [];
      gameData.game_players = gameData.players || [];
      setGame(gameData);

      setMutinyUsed(gameData.AllScores?.some((s) => s.AgendaTitle === "Mutiny"));
      setCdlUsed(gameData.AllScores?.some((s) => s.AgendaTitle === "Classified Document Leaks"));
      setCrownUsed(gameData.AllScores?.some((s) => s.Type?.toLowerCase() === "relic" && s.RelicTitle === "The Crown of Emphidia"));
      const obsidianScore = gameData.AllScores?.find(
        (s) => s.Type === "relic" && s.RelicTitle === "The Obsidian"
      );
      setObsidianHolderId(obsidianScore?.PlayerID || null);
      const initialSecrets = {};
      (gameData.players || []).forEach((p) => {
        initialSecrets[p.PlayerID || p.id] = 0;
      });
      setSecretCounts(initialSecrets);

      const scoreMap = {};
      (Array.isArray(scoresData) ? scoresData : scoresData?.value || []).forEach((entry) => {
        scoreMap[entry.objective_id] = entry.scored_by || [];
      });
      setObjectiveScores(scoreMap);

      const normalizedObjectives = Array.isArray(objectiveData)
        ? objectiveData
        : objectiveData?.value || [];


      setObjectives(normalizedObjectives);

      const normalizedSecrets = secretData.map((obj) => ({
        id: obj.ID,
        name: obj.name,
        phase: obj.phase?.toLowerCase() || obj.Phase?.toLowerCase() || "",
        ...obj,
      }));
      setSecretObjectives(normalizedSecrets);

    })();

    const match = game?.AllScores?.find(
      (s) => s.Type === "Agenda" && s.AgendaTitle === "Political Censure"
    );
    setCensureHolder(match?.PlayerID || null);
  }, [gameId]);


  return {
    game,
    objectives,
    secretObjectives,
    secretCounts,
    setSecretCounts,
    objectiveScores,
    mutinyUsed,
    setGame,
    setObjectiveScores,
    refreshGameState,
    censureHolder,
    setMutinyUsed,
    setCdlUsed,
    crownUsed,
    obsidianHolderId,
    setObsidianHolderId,
  };
}
