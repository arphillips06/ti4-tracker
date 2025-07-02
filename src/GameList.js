import React, { useEffect, useState } from 'react';
import axios from 'axios';

function GameList() {
  const [games, setGames] = useState([]);

useEffect(() => {
  axios.get('http://localhost:8080/games')
    .then(response => {
      console.log("API response:", response.data); // 👈 See the shape
      setGames(response.data);
    })
    .catch(error => console.error('Error fetching games:', error));
}, []);

  return (
    <div className="container mt-4">
      <h1 className="mb-4">TI4 Games</h1>
      {games.length === 0 ? (
        <p>No games found.</p>
      ) : (
<ul className="list-group">
  {games.map(game => (
    <li key={game.id} className="list-group-item">
      <strong>Game #{game.id}</strong> – {game.players?.length || 0} players – Round {game.current_round}
    </li>
  ))}
</ul>
      )}
    </div>
  );
}

export default GameList;
