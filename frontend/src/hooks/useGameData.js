import { useEffect, useState } from "react";
import { getJSON } from "../utils/helpers";
import {
  isAgendaScore,
  isRelicScore,
} from "../utils/selectors";

const normalizeScore = (s) => ({
  ...s,
  PlayerID: s.PlayerID ?? s.player_id,
  ObjectiveID: s.ObjectiveID ?? s.objective_id,
  RoundID: s.RoundID ?? s.round_id,
  Type: s.Type ?? s.type,
  Points: s.Points ?? s.points,
  AgendaTitle: s.AgendaTitle ?? s.agenda_title,
  RelicTitle: s.RelicTitle ?? s.relic_title,
});

const normalizeObjective = (obj) => {
  const rawPhase = obj.phase ?? obj.Phase ?? "";
  return {
    id: obj.ID ?? obj.id,
    name: obj.name ?? obj.Name,
    phase: typeof rawPhase === "string" ? rawPhase.toLowerCase() : "",
    ...obj,
  };
};

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



  const fetchGame = () => getJSON(`/games/${gameId}`);
  const fetchObjectives = () => getJSON(`/games/${gameId}/objectives`);
  const fetchScores = () => getJSON(`/games/${gameId}/objectives/scores`);
  const fetchSecrets = () => getJSON(`/objectives/secrets/all`);

  const refreshGameState = async () => {
    const [gameData, objectivesData, scoresData] = await Promise.all([
      fetchGame(),
      fetchObjectives(),
      fetchScores(),
    ]);

    const normalizedScores = (gameData.all_scores || []).map(normalizeScore);
    const normalizedObjectives = Array.isArray(objectivesData)
      ? objectivesData.map(normalizeObjective)
      : objectivesData?.value || [];

    const scoreMap = {};
    (Array.isArray(scoresData) ? scoresData : scoresData?.value || []).forEach((entry) => {
      scoreMap[entry.objective_id] = entry.scored_by || [];
    });

    setGame({
      ...gameData,
      AllScores: normalizedScores,
      game_players: gameData.players || [],
      winner_id: gameData.winner_id ?? gameData.WinnerID,
    });
    setObjectives(normalizedObjectives);
    setObjectiveScores(scoreMap);

    const obsidianScore = normalizedScores.find((s) => isRelicScore(s, "The Obsidian"));
    setObsidianHolderId(obsidianScore?.PlayerID || null);
  };
  useEffect(() => {
    (async () => {
      const [gameData, objectiveData, secretData, scoresData] = await Promise.all([
        fetchGame(),
        fetchObjectives(),
        fetchSecrets(),
        fetchScores(),
      ]);

      const normalizedScores = (gameData.all_scores || []).map(normalizeScore);
      const normalizedObjectives = Array.isArray(objectiveData)
        ? objectiveData.map(normalizeObjective)
        : objectiveData?.value || [];
      const normalizedSecrets = (Array.isArray(secretData) ? secretData : []).map(normalizeObjective);

      const scoreMap = {};
      (Array.isArray(scoresData) ? scoresData : scoresData?.value || []).forEach((entry) => {
        scoreMap[entry.objective_id] = entry.scored_by || [];
      });

      setGame({
        ...gameData,
        AllScores: normalizedScores,
        game_players: gameData.players || [],
      });
      setObjectives(normalizedObjectives);
      setSecretObjectives(normalizedSecrets);
      setObjectiveScores(scoreMap);

      setMutinyUsed(normalizedScores.some((s) => s.AgendaTitle === "Mutiny"));
      setCdlUsed(normalizedScores.some((s) => s.AgendaTitle === "Classified Document Leaks"));
      setCrownUsed(normalizedScores.some((s) => isRelicScore(s, "The Crown of Emphidia")));
      setObsidianHolderId(normalizedScores.find((s) => isRelicScore(s, "The Obsidian"))?.PlayerID || null);

      const initialSecrets = {};
      (gameData.players || []).forEach((p) => {
        initialSecrets[p.PlayerID || p.id] = 0;
      });
      setSecretCounts(initialSecrets);

      const censureMatch = normalizedScores.find((s) => isAgendaScore(s, "Political Censure"));
      setCensureHolder(censureMatch?.PlayerID || null);
    })();
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
