import { useState, useEffect } from 'react';
import PlayerRow from '../components/PlayerRow';

export default function NewGamePage() {
  const [playerCount, setPlayerCount] = useState(3);
  const [players, setPlayers] = useState([]);
  const [factions, setFactions] = useState([]);
  const [winningPoints, setWinningPoints] = useState(10);
  const [useObjectiveDecks, setUseObjectiveDecks] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/factions')
      .then(res => res.json())
      .then(data => setFactions(data))
      .catch(err => console.error("Failed to load factions", err));
  }, []);

  useEffect(() => {
    const updated = [...players];
    while (updated.length < playerCount) {
      updated.push({ name: "", faction: "", color: "#000000" });
    }
    while (updated.length > playerCount) {
      updated.pop();
    }
    setPlayers(updated);
  }, [playerCount]);

  const handlePlayerChange = (index, field, value) => {
    const updated = [...players];
    updated[index][field] = value;
    setPlayers(updated);
  };

const startGame = async () => {
  try {
    const payload = {
      winning_points: winningPoints,
      use_objective_decks: useObjectiveDecks,
      players: players.map(p => ({
        name: p.name,
        faction: p.faction
      }))
    };
console.log("Submitting payload:", JSON.stringify(payload, null, 2));

    const res = await fetch('http://localhost:8080/games', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    });

    if (!res.ok) {
      const errorText = await res.text();  // <- helpful for debugging
      throw new Error(`Failed to create game: ${errorText}`);
    }

    const data = await res.json();
    const newGameId = data.game.id;

    window.location.href = `/games/${newGameId}`;
  } catch (error) {
    console.error("Error starting game:", error);
    alert("Failed to start game. See console for details.");
  }
};

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

      <div className="mb-6">
        <label className="block font-medium mb-1">Use Objective Decks:</label>
        <input
          type="checkbox"
          checked={useObjectiveDecks}
          onChange={(e) => setUseObjectiveDecks(e.target.checked)}
          className="mr-2"
        />
        <span>{useObjectiveDecks ? "Yes" : "No"}</span>
      </div>

      <label className="block mb-2 font-medium">Number of Players: {playerCount}</label>
      <input
        type="range"
        min={3}
        max={8}
        value={playerCount}
        onChange={(e) => setPlayerCount(Number(e.target.value))}
        className="w-full mb-6"
      />

      {players.map((player, idx) => (
        <PlayerRow
          key={idx}
          index={idx}
          player={player}
          factions={factions}
          onChange={handlePlayerChange}
        />
      ))}

      <div className="mt-6 text-right">
        <button
          onClick={startGame}
          className="bg-blue-600 text-white px-5 py-2 rounded hover:bg-blue-700"
        >
          Start Game
        </button>
      </div>
    </div>
  );
}
