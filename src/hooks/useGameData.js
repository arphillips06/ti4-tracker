import { useEffect, useState } from "react";

export default function useGameData(gameId) {
  const [game, setGame] = useState(null);
  const [objectives, setObjectives] = useState([]);
  const [objectiveScores, setObjectiveScores] = useState({});
  const [secretObjectives, setSecretObjectives] = useState([]);
  const [secretCounts, setSecretCounts] = useState({});
  const [mutinyUsed, setMutinyUsed] = useState(false);
  const [censureHolder, setCensureHolder] = useState(null);


  const fetchGame = async () => {
    const res = await fetch(`http://localhost:8080/games/${gameId}`);
    return res.json();
  };

  const fetchObjectives = async () => {
    const res = await fetch(`http://localhost:8080/games/${gameId}/objectives`);
    return res.json();
  };

  const fetchScores = async () => {
    const res = await fetch(`http://localhost:8080/games/${gameId}/objectives/scores`);
    return res.json();
  };

  const fetchSecrets = async () => {
    const res = await fetch("http://localhost:8080/objectives/secrets/all");
    return res.json();
  };

  const refreshGameState = async () => {
    const [gameData, objectivesData, scoresData] = await Promise.all([
      fetchGame(),
      fetchObjectives(),
      fetchScores(),
    ]);

    setGame(gameData);

    const map = {};
    (Array.isArray(scoresData) ? scoresData : scoresData?.value || []).forEach((entry) => {
      map[entry.objective_id] = entry.scored_by || [];
    });
    setObjectiveScores(map);

    setObjectives(Array.isArray(objectivesData) ? objectivesData : objectivesData?.value || []);
  };

  useEffect(() => {
    (async () => {
      const [gameData, secretData, scoresData, objectiveData] = await Promise.all([
        fetchGame(),
        fetchSecrets(),
        fetchScores(),
        fetchObjectives(),
      ]);

      setGame(gameData);
      setMutinyUsed(gameData.AllScores?.some((s) => s.AgendaTitle === "Mutiny"));

      const initialSecrets = {};
      gameData.players.forEach((p) => {
        initialSecrets[p.PlayerID] = 0;
      });
      setSecretCounts(initialSecrets);

      const map = {};
      (Array.isArray(scoresData) ? scoresData : scoresData?.value || []).forEach((entry) => {
        map[entry.objective_id] = entry.scored_by || [];
      });

      setObjectiveScores(map);
      setObjectives(Array.isArray(objectiveData) ? objectiveData : objectiveData?.value || []);

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
  };
}
