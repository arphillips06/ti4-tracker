// src/pages/NewGamePage.js
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import PlayerInputRow from "../components/players/PlayerInputRow";
import factionColors from "../data/factionColors";
import API_BASE_URL from "../config";
import '../pages/NewGamePage.css';
import '../pages/stats.css';

const FACTIONS = Object.entries(factionColors).map(([key, data]) => ({
  key,
  label: data.label,
}));

export default function NewGamePage() {
  const navigate = useNavigate();

  const [numPlayers, setNumPlayers] = useState(3);
  const [players, setPlayers] = useState(
    Array.from({ length: 3 }, () => ({ name: "", faction: "", color: "" }))
  );
  const [useObjectives, setUseObjectives] = useState(true);
  const [winningPoints, setWinningPoints] = useState(10);

  const handlePlayerChange = (index, field, value) => {
    const updatedPlayers = [...players];

    if (field === "faction") {
      const normalize = (str) =>
        str.toLowerCase().replace(/^the\s+/, "").trim();

      const factionKey = Object.keys(factionColors).find(
        (k) => normalize(factionColors[k].label) === normalize(value) || k === value
      );

      if (factionKey) {
        updatedPlayers[index].faction = factionKey;

        const usedColors = updatedPlayers
          .filter((_, i) => i !== index)
          .map((p) => p.color);

        const colorOptions = factionColors[factionKey]?.colors || [];
        const firstAvailable = colorOptions.find((c) => !usedColors.includes(c));

        if (firstAvailable) {
          updatedPlayers[index].color = firstAvailable;
        }
      }
    } else {
      updatedPlayers[index][field] = value;
    }

    setPlayers(updatedPlayers);
  };


  const startGame = async () => {
    try {
      const payload = {
        winning_points: winningPoints,
        use_objective_decks: useObjectives,
        players: players.map((p) => ({
          name: p.name,
          faction: p.faction,
        })),
      };


      const res = await fetch(`${API_BASE_URL}/games`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`Failed to create game: ${errorText}`);
      }

      const data = await res.json();

      const newGameId = data.game.id; // <- likely fix here

      if (!data.game.id) throw new Error("No game ID returned from backend");
      navigate(`/game/${data.game.id}`);

      navigate(`/games/${newGameId}`);

    } catch (error) {
      console.error("Error starting game:", error);
      alert("Failed to start game. See console for details.");
    }
  };


  const selectedFactions = players.map((p) => p.faction).filter(Boolean);
  const selectedColors = players.map((p) => p.color).filter(Boolean);

  return (
    <div className="p-4 max-w-3xl mx-auto">
      <h1 className="ti-header">Start a New Game</h1>

      <div className="stat-card mb-6">
        <div className="mb-4">
          <label className="label-gold mb-1">
            Winning Points:
          </label>
          <select
            value={winningPoints}
            onChange={(e) => setWinningPoints(Number(e.target.value))}
            style={{
              width: "60px",
              padding: "4px 6px",
              fontSize: "0.9rem",
              borderRadius: "4px",
            }}
          >
            <option value={10}>10</option>
            <option value={14}>14</option>
          </select>
        </div>

        <div className="mb-4">
          <div className="checkbox-container">
            <label className="label-gold">Use Objective Decks:</label>
            <input
              type="checkbox"
              checked={useObjectives}
              onChange={(e) => setUseObjectives(e.target.checked)}
            />
          </div>
        </div>


        <div className="mb-4">
          <label className="label-gold mb-1">
            Number of Players: {numPlayers}
          </label>
          <input
            type="range"
            min={2}
            max={8}
            value={numPlayers}
            onChange={(e) => {
              const newCount = Number(e.target.value);
              setNumPlayers(newCount);
              setPlayers((prev) =>
                Array.from({ length: newCount }, (_, i) => prev[i] || { name: "", faction: "", color: "" })
              );
            }}
          />
        </div>
      </div>
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(350px, 1fr))",
          gap: "12px",
          marginBottom: "20px",
        }}
      >
        {players.map((player, i) => (
          <div
            key={i}
            style={{
              backgroundColor: "#0d0d1a",
              border: "1px solid #e0c87344",
              borderRadius: "8px",
              padding: "10px",
            }}
          >
            <PlayerInputRow
              index={i}
              value={player}
              onFactionChange={(index, faction) =>
                handlePlayerChange(index, "faction", faction)
              }
              onColorChange={(index, color) =>
                handlePlayerChange(index, "color", color)
              }
              onNameChange={(index, name) =>
                handlePlayerChange(index, "name", name)
              }
              factions={FACTIONS}
              selectedFactions={selectedFactions}
              selectedColors={selectedColors}
            />
          </div>
        ))}
      </div>


      <button onClick={startGame} className="btn btn-primary mt-4">
        Start Game
      </button>
    </div>
  );
}
