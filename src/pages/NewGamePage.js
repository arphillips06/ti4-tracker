import { useState, useEffect } from 'react';
import PlayerRow from '../components/PlayerRow';

export default function NewGamePage() {
  const [playerCount, setPlayerCount] = useState(3);
  const [players, setPlayers] = useState([]);
  const [factions, setFactions] = useState([]);

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

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">Start a New Game</h2>

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
    </div>
  );
}
