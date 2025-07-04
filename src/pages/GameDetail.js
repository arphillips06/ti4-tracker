// GameDetail.js

import { useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';
import factionImageMap from '../data/factionIcons';
import PlayerRow from "../components/PlayerRow";



export default function GameDetail() {
  const { gameId } = useParams();
  const [game, setGame] = useState(null);
  const [objectives, setObjectives] = useState([]);
  const [scoringMode, setScoringMode] = useState(false);
  

  useEffect(() => {
    fetch(`http://localhost:8080/games/${gameId}`)
      .then(async (res) => {
        const text = await res.text();
        console.log("Raw response from backend:", text);
        if (!res.ok) throw new Error(`Game not found. Server said: ${text}`);

        try {
          const data = JSON.parse(text);
          console.log("Loaded game:", data);
          setGame(data);
        } catch (err) {
          console.error("Failed to parse JSON:", err);
          throw new Error("Invalid JSON from server");
        }
      })
      .catch((err) => console.error("Error loading game:", err));

    fetch(`http://localhost:8080/games/${gameId}/objectives`)
      .then((res) => res.json())
      .then(setObjectives)
      .catch((err) => console.error("Error loading objectives:", err));
  }, [gameId]);

const scoreObjective = async (playerId, objectiveId) => {
  const payload = {
    game_id: parseInt(gameId),
    player_id: playerId,
    objective_id: objectiveId,
  };

  console.log("Scoring payload:", payload);

  try {
    const res = await fetch("http://localhost:8080/score", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(errorText);
    }

    const result = await res.json();
    console.log("Score response:", result);

    // ✅ Refresh game and objectives after scoring
    const gameRes = await fetch(`http://localhost:8080/games/${gameId}`);
    const gameText = await gameRes.text();
    if (!gameRes.ok) throw new Error(gameText);
    setGame(JSON.parse(gameText));

    const objRes = await fetch(`http://localhost:8080/games/${gameId}/objectives`);
    if (!objRes.ok) throw new Error("Failed to fetch objectives");
    const updatedObjectives = await objRes.json();
    setObjectives(updatedObjectives);
  } catch (err) {
    console.error("Scoring failed:", err);
    alert("Scoring failed. See console for details.");
  }
};

  const advanceRound = async () => {
    try {
      const res = await fetch(`http://localhost:8080/games/${gameId}/advance-round`, {
        method: "POST",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText);
      }

      const result = await res.json();
      console.log("Round advanced:", result);

      // Re-fetch game data
      const gameRes = await fetch(`http://localhost:8080/games/${gameId}`);
      const gameText = await gameRes.text();
      if (!gameRes.ok) throw new Error(gameText);
      setGame(JSON.parse(gameText));
      console.log("Refetched game after advance:", gameText);

      // Re-fetch objectives
      const objRes = await fetch(`http://localhost:8080/games/${gameId}/objectives`);
      if (!objRes.ok) throw new Error("Failed to fetch objectives");
      const updatedObjectives = await objRes.json();
      setObjectives(updatedObjectives);
      console.log("Refetched objectives after advance:", updatedObjectives);

    } catch (err) {
      console.error("Failed to advance round:", err);
      alert("Could not advance round. See console for details.");
    }
  };

const getMergedPlayerData = () => {
  
  const scoreMap = new Map(game.scores?.map(s => [s.player_id, s]) || []);

  return game.players
    .map((p) => {
      const playerId = p.PlayerID;
      const name = p.Player?.Name || "Unknown";
      const faction = p.Faction || "Unknown Faction";
      const color = p.color || "#000000";
      const scoreEntry = scoreMap.get(playerId);
      const points = scoreEntry?.points || 0;

      // Generate a faction icon key from the faction name
      const factionKey = faction
        .replace(/^The\s+/i, "")      // remove "The "
        .replace(/\s+/g, "")          // remove all spaces
        .replace(/[^a-zA-Z0-9]/g, "") // remove symbols
        + (faction.toLowerCase().includes("keleres") ? "FactionSymbol" : ""); // handle Keleres oddity

      return {
        player_id: playerId,
        name,
        faction,
        factionKey,
        color,
        points,
      };
    })
    .sort((a, b) => b.points - a.points);
};



  if (!game || !game.players) {
    return <div className="p-6">Loading game data...</div>;
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      {/* Game Header */}
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2 className="h4">
          Round {game.current_round} | {game.winning_points} Point Game
        </h2>
        <div className="d-flex gap-3 align-items-center">
          <div className="form-check form-switch">
            <input
              className="form-check-input"
              type="checkbox"
              checked={scoringMode}
              onChange={(e) => setScoringMode(e.target.checked)}
            />
            <label className="form-check-label">Score Objectives</label>
          </div>
          <button className="btn btn-outline-primary btn-sm" onClick={advanceRound}>
            Advance Round
          </button>
        </div>
      </div>

      <div className="row">
        {/* Objectives Section */}
        <div className="col-md-6">
          <h4>Objectives</h4>
          {objectives.map((obj) => (
            <div
              key={obj.ID}
              className="card mb-3 border-warning"
              style={{ backgroundColor: "#fff8dc" }}
            >
              <div className="card-body">
                <h5 className="card-title">
                  {obj.Objective?.Name || "Unnamed Objective"}
                </h5>

                {scoringMode && (
                  <div className="d-flex flex-wrap gap-2 mt-2">
                    {(game.players || []).map((p) => {
                      const name = p.Player?.Name || "Unknown";
                      const color = p.color || "#000";
                      const playerId = p.PlayerID;
                      const objectiveId = obj.Objective?.ID;

                      return (
                        <button
                          key={playerId}
                          className="btn btn-sm text-white"
                          style={{ backgroundColor: color }}
                          onClick={() => scoreObjective(playerId, objectiveId)}
                        >
                          {name}
                        </button>
                      );
                    })}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>

        {/* Players Section */}
        <div className="col-md-6">
{getMergedPlayerData().map((entry) => {
  return (
    <div
      key={entry.player_id}
      className="card mb-3 border-start border-5"
      style={{ borderColor: entry.color }}
    >
      <div className="card-body">
        <div className="fw-semibold small mb-1">{entry.name}</div>
        <div className="d-flex align-items-center gap-2">
          <img
            src={`/faction-icons/${entry.factionKey}.webp`}
            alt={entry.faction}
            style={{
              width: "24px",
              height: "24px",
              borderRadius: "50%",
              objectFit: "contain",
              backgroundColor: "transparent",
            }}
            onError={(e) => (e.target.style.display = "none")}
          />
          <div className="text-muted small fst-italic">{entry.faction}</div>
        </div>
        <div className="mt-1 small">Points: {entry.points}</div>
      </div>
    </div>
  );
})}
        </div>
      </div>
    </div>
  );
}
