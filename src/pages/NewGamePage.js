// src/pages/NewGamePage.js
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import PlayerInputRow from "../components/PlayerInputRow";
import factionColors from "../data/factionColors";
import API_BASE_URL from "../config";

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

      console.log("Submitting payload:", JSON.stringify(payload, null, 2));

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
      console.log("Response data from backend:", data);
      console.log("data.game:", data.game);

      const newGameId = data.game.id; // <- likely fix here

      if (!data.game.id) throw new Error("No game ID returned from backend");
      navigate(`/game/${data.game.id}`);

      console.log("Navigating to new game:", newGameId);
      navigate(`/games/${newGameId}`);

    } catch (error) {
      console.error("Error starting game:", error);
      alert("Failed to start game. See console for details.");
    }
  };


  const selectedFactions = players.map((p) => p.faction).filter(Boolean);
  const selectedColors = players.map((p) => p.color).filter(Boolean);

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">Start a New Game</h2>

      <div className="mb-4">
        <label className="block font-medium mb-1">Winning Points:</label>
        <select
          value={winningPoints}
          onChange={(e) => setWinningPoints(Number(e.target.value))}
          className="border rounded px-3 py-2"
        >
          <option value={10}>10</option>
          <option value={14}>14</option>
        </select>
      </div>

      <div className="mb-4">
        <label className="block font-medium mb-1">Use Objective Decks:</label>
        <input
          type="checkbox"
          checked={useObjectives}
          onChange={(e) => setUseObjectives(e.target.checked)}
        />
      </div>

      <div className="mb-4">
        <label className="block font-medium mb-1">
          Number of Players: {numPlayers}
        </label>
        <input
          type="range"
          min={1}
          max={10}
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

      {players.map((player, i) => (
        <PlayerInputRow
          key={i}
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
      ))}

      <button onClick={startGame} className="btn btn-primary mt-4">
        Start Game
      </button>
    </div>
  );
}
