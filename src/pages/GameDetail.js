import React, { useState, useEffect, useMemo } from "react";
import { useParams } from "react-router-dom";

import PlayerSidebar from "../components/PlayerSidebar";
import GameControls from "../components/GameControls";
import ObjectivesGrid from "../components/ObjectivesGrid";
import GameNavbar from "../components/GameNavbar";
import MutinyModal from "../components/MutinyModal";
import PoliticalCensureModal from "../components/PoliticalCensureModal";
import SeedOfEmpireModal from "../components/SeedOfEmpireModal";
import ClassifiedDocumentLeaksModal from "../components/ClassifiedDocumentLeaksModal";


import useGameData from "../hooks/useGameData";
const API_URL = "http://localhost:8080";


export default function GameDetail() {
  const { gameId } = useParams();

  const {
    game,
    objectives,
    secretObjectives,
    secretCounts,
    setSecretCounts,
    objectiveScores,
    mutinyUsed,
    setMutinyUsed,
    setGame,
    setObjectiveScores,
    refreshGameState,
    censureHolder,
    scores,
  } = useGameData(gameId);

  const [scoringMode, setScoringMode] = useState(false);
  const [expandedPlayers, setExpandedPlayers] = useState({});
  const [localScored, setLocalScored] = useState({});
  const [showAgendaModal, setShowAgendaModal] = useState(false);
  const [mutinyVotes, setMutinyVotes] = useState([]);
  const [mutinyResult, setMutinyResult] = useState("for");
  const [mutinyAbstained, setMutinyAbstained] = useState(false);
  const [custodiansScorerId, setCustodiansScorerId] = useState(null);
  const [showCensureModal, setShowCensureModal] = useState(false);
  const [showSeedModal, setShowSeedModal] = useState(false);
  const [assigningObjective, setAssigningObjective] = useState(null); // { roundNumber, stage }
  const [selectedObjectiveId, setSelectedObjectiveId] = useState(null);
  const [agendaModal, setAgendaModal] = useState(null);



const isLoading = !game || !game.ID;


  useEffect(() => {
    const scorer = game?.all_scores?.find((s) => s.Type === "mecatol");
    setCustodiansScorerId(scorer?.PlayerID || null);
  }, [game]);

  const unscoreObjective = async (playerId, objectiveId) => {
    try {
      const res = await fetch(`${API_URL}/unscore`, {
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
  };

const assignObjective = async (gameId, roundId, objectiveId) => {
    console.log("🚀 assignObjective payload", { gameId, roundId, objectiveId });

  try {
    const res = await fetch("http://localhost:8080/assign_objective", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: gameId,
        round_id: roundId,
        objective_id: objectiveId,
      }),
    });

    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText);
    }

    console.log(`Assigned objective ${objectiveId} to round ${roundId} in game ${gameId}`);
    await refreshGameState();
  } catch (err) {
    console.error("Failed to assign objective:", err);
    alert("Failed to assign objective. See console for details.");
  }
};



const handleSeedSubmit = async (result) => {
  try {
    await fetch(`${API_URL}/agenda/seed`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        round_id: game.current_round_id || 0,
        result,
      }),
    });
    await refreshGameState();
  } catch (err) {
    console.error("Failed to apply Seed of an Empire:", err);
    alert("Failed to apply agenda. See console.");
  }
};

const handleMutinySubmit = async () => {
  try {
    await fetch(`${API_URL}/agenda/mutiny`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        round_id: game.current_round_id || 0,
        result: mutinyResult,
        for_votes: mutinyAbstained ? [] : mutinyVotes,
      }),
    });

    const updatedGame = await fetch(`${API_URL}/games/${gameId}`).then((r) => r.json());
    setGame(updatedGame);

    setShowAgendaModal(false);
    setMutinyUsed(true);
  } catch (err) {
    console.error("Failed to apply mutiny agenda:", err);
    alert("Failed to apply agenda. See console.");
  }
};
const handlePoliticalCensureSubmit = async ({ playerId, gained }) => {
  try {
    await fetch(`${API_URL}/agenda/political-censure`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        round_id: game.current_round_id || 0,
        player_id: playerId,
        gained: gained,
      }),
    });

    await refreshGameState();
  } catch (err) {
    console.error("Failed to apply political censure:", err);
    alert("Failed to apply agenda. See console.");
  }
};

const handleClassifiedSubmit = async (playerId, objectiveId) => {
  try {
    await fetch(`${API_URL}/agenda/classified-document-leaks`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        game_id: parseInt(gameId),
        player_id: playerId,
        objective_id: objectiveId,
      }),
    });

    await refreshGameState();
    setAgendaModal(null);
  } catch (err) {
    console.error("Failed to apply Classified Document Leaks:", err);
    alert("Failed to apply agenda. See console.");
  }
};


const scoreObjective = async (playerId, objectiveId, playerName) => {
  const payload = {
    game_id: parseInt(gameId),
    player_id: playerId,
    objective_id: objectiveId,
  };

  try {
    const res = await fetch(`${API_URL}/score`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const errorText = await res.text();
      console.error("Scoring rejected:", errorText);
      alert("Scoring rejected by backend: " + errorText);
      return;
    }
    setLocalScored((prev) => {
      const current = prev[objectiveId] || [];
      return {
        ...prev,
        [objectiveId]: [...new Set([...current, playerId])],
      };
    });

    await refreshGameState();
    return true
  } catch (err) {
    console.error("Scoring failed:", err);
    alert("Scoring failed. See console.");
    return false
  }
};


const advanceRound = async () => {
  try {
    const res = await fetch(`${API_URL}/games/${gameId}/advance-round`, {
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
};
const getMergedPlayerData = (sort = true) => {
  const scoreMap = new Map(game?.scores?.map((s) => [s.player_id, s]) || []);

  const merged = (game?.game_players || []).map((gp) => {
    const p = gp.Player || {};
    const name = p.Name || "Unknown";
    const faction = gp.Faction || "Unknown Faction";  // ← fixed this line
    const factionKey =
      faction.replace(/^The\s+/i, "")
        .replace(/\s+/g, "")
        .replace(/[^a-zA-Z0-9]/g, "") +
      (faction.toLowerCase().includes("keleres") ? "FactionSymbol" : "");

    return {
      player_id: p.ID,
      name,
      faction,
      factionKey,
      color: gp.color || "#000",
      points: scoreMap.get(p.ID)?.points || 0,
    };
  });

  return sort ? merged.sort((a, b) => b.points - a.points) : merged;
};



const groupedScoredSecrets = useMemo(() => {
  const result = {};

  (game?.all_scores || []).forEach((s) => {
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
}, [game?.all_scores, secretObjectives]);

const playersUnsorted = getMergedPlayerData(false);
const playersSorted = getMergedPlayerData(true);
console.log("🔎 playersSorted:", playersSorted);
if (!game) return <div className="p-4 text-warning">Loading game data...</div>;

return (
  <>
    <GameNavbar
      mutinyUsed={mutinyUsed}
      setShowAgendaModal={setShowAgendaModal}
      setShowCensureModal={setShowCensureModal}
      setShowSeedModal={setShowSeedModal}
      setAgendaModal={setAgendaModal}
    />


    <div className="p-6 max-w-7xl mx-auto">
      <GameControls
        game={game}
        scoringMode={scoringMode}
        setScoringMode={setScoringMode}
        onAdvanceRound={advanceRound}
      />
      <div className="d-flex flex-row align-items-start gap-4 flex-wrap">
        <ObjectivesGrid
          game={game}
          refreshGameState={refreshGameState}
          objectives={objectives}
          playersUnsorted={playersUnsorted}
          gameId={game.id || game.ID}
          objectiveScores={objectiveScores}
          localScored={localScored}
          scoringMode={scoringMode}
          scoreObjective={scoreObjective}
          useObjectiveDecks={game?.use_objective_decks}
          assigningObjective={assigningObjective}
          setAssigningObjective={setAssigningObjective}
          assignObjective={assignObjective}

        />



        <div style={{ flex: "0 1 300px" }}>
          <PlayerSidebar
            playersSorted={playersSorted}
            expandedPlayers={expandedPlayers}
            setExpandedPlayers={setExpandedPlayers}
            game={game}
            scoreObjective={scoreObjective}
            unscoreObjective={unscoreObjective}
            secretCounts={secretCounts}
            setSecretCounts={setSecretCounts}
            secretObjectives={secretObjectives}
            objectiveScores={objectiveScores}
            gameId={gameId}
            setGame={setGame}
            setObjectiveScores={setObjectiveScores}
            refreshGameState={refreshGameState}
            custodiansScored={!!game?.all_scores?.some((s) => s.Type === "mecatol")}
          />
        </div>
      </div>
      <MutinyModal
        show={showAgendaModal}
        onClose={() => setShowAgendaModal(false)}
        onSubmit={handleMutinySubmit}
        mutinyResult={mutinyResult}
        setMutinyResult={setMutinyResult}
        mutinyAbstained={mutinyAbstained}
        setMutinyAbstained={setMutinyAbstained}
        mutinyVotes={mutinyVotes}
        setMutinyVotes={setMutinyVotes}
        players={playersUnsorted}
      />
      <SeedOfEmpireModal
        show={showSeedModal}
        onClose={() => setShowSeedModal(false)}
        onSubmit={handleSeedSubmit}
      />

      <PoliticalCensureModal
        show={showCensureModal}
        onClose={() => setShowCensureModal(false)}
        onSubmit={handlePoliticalCensureSubmit}
        players={playersUnsorted.map((p) => ({
          ...p,
          agendaScores: game.all_scores?.filter(
            (s) => s.PlayerID === p.player_id && s.Type?.toLowerCase() === "agenda"
          ) || [],
        }))}
      />
      <ClassifiedDocumentLeaksModal
        show={agendaModal === "Classified Document Leaks"}
        players={playersUnsorted}
        secretObjectives={secretObjectives}
        scoredSecrets={groupedScoredSecrets}
        onClose={() => setAgendaModal(null)}
        onSubmit={handleClassifiedSubmit}
      />


    </div>
  </>
)
}