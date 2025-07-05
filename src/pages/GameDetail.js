import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

import PlayerSidebar from "../components/PlayerSidebar";
import GameControls from "../components/GameControls";
import ObjectivesGrid from "../components/ObjectivesGrid";
import GameNavbar from "../components/GameNavbar";
import MutinyModal from "../components/MutinyModal";

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
  } = useGameData(gameId);

  const [scoringMode, setScoringMode] = useState(false);
  const [expandedPlayers, setExpandedPlayers] = useState({});
  const [localScored, setLocalScored] = useState({});
  const [showAgendaModal, setShowAgendaModal] = useState(false);
  const [mutinyVotes, setMutinyVotes] = useState([]);
  const [mutinyResult, setMutinyResult] = useState("for");
  const [mutinyAbstained, setMutinyAbstained] = useState(false);
  const [custodiansScored, setCustodiansScored] = useState(false);

  useEffect(() => {
    setCustodiansScored(game?.AllScores?.some((s) => s.Type === "mecatol") || false);
  }, [game]);
  if (!game || !game.players) {
    return <div className="p-6">Loading game data...</div>;
  }

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
          [objectiveId]: [...new Set([...current, playerName])],
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

    const merged = (game?.players || []).map((p) => {
      const name = p.Player?.Name || "Unknown";
      const faction = p.Faction || "Unknown Faction";
      const factionKey =
        faction.replace(/^The\s+/i, "")
          .replace(/\s+/g, "")
          .replace(/[^a-zA-Z0-9]/g, "") +
        (faction.toLowerCase().includes("keleres") ? "FactionSymbol" : "");
      return {
        player_id: p.PlayerID,
        name,
        faction,
        factionKey,
        color: p.color || "#000",
        points: scoreMap.get(p.PlayerID)?.points || 0,
      };
    });

    return sort ? merged.sort((a, b) => b.points - a.points) : merged;
  };

  const playersUnsorted = getMergedPlayerData(false);
  const playersSorted = getMergedPlayerData(true);

  return (
    <>
      <GameNavbar mutinyUsed={mutinyUsed} setShowAgendaModal={setShowAgendaModal} />


      <div className="p-6 max-w-7xl mx-auto">
        <GameControls
          game={game}
          scoringMode={scoringMode}
          setScoringMode={setScoringMode}
          onAdvanceRound={advanceRound}
        />
        <div className="d-flex flex-row align-items-start gap-4 flex-wrap">
          <ObjectivesGrid
            objectives={objectives}
            playersUnsorted={playersUnsorted}
            objectiveScores={objectiveScores}
            localScored={localScored}
            scoringMode={scoringMode}
            scoreObjective={scoreObjective}
          />
          <div style={{ flex: "0 1 300px" }}>
            <PlayerSidebar
              playersSorted={playersSorted}
              expandedPlayers={expandedPlayers}
              setExpandedPlayers={setExpandedPlayers}
              objectiveScores={objectiveScores}
              game={game}
              setGame={setGame}
              secretObjectives={secretObjectives}
              setObjectiveScores={setObjectiveScores}
              secretCounts={secretCounts}
              setSecretCounts={setSecretCounts}
              scoreObjective={scoreObjective}
              unscoreObjective={unscoreObjective}
              custodiansScored={custodiansScored}
              gameId={gameId}
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

      </div>
    </>
  )
}