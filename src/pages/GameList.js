import React, { useEffect, useState } from 'react';
import axios from 'axios';
import API_BASE_URL from "../config";

function GameList() {
  const [games, setGames] = useState([]);

  function formatDuration(start, end) {
    if (!end) return "Ongoing";

    const ms = end - start;
    const totalMinutes = Math.floor(ms / (1000 * 60));
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  }

  async function fetchScoresForGame(id) {
    try {
      const res = await axios.get(`${API_BASE_URL}/games/${id}/score-summary`);
      return res.data;
    } catch (e) {
      console.error("Error fetching scores for game", id, e);
      return [];
    }
  }

  useEffect(() => {
    async function loadGames() {
      try {
        const res = await axios.get(`${API_BASE_URL}/games`);
        const gamesData = res.data;

        const withScores = await Promise.all(
          gamesData.map(async (game) => {
            const scores = await fetchScoresForGame(game.ID);
            return { ...game, Scores: scores };
          })
        );

        setGames(withScores);
      } catch (error) {
        console.error('Error fetching games:', error);
      }
    }

    loadGames();
  }, []);

  return (
    <div className="container mt-4">
      <h1 className="mb-4">TI4 Games</h1>
      {games.length === 0 ? (
        <p>No games found.</p>
      ) : (
        <ul className="list-group">
          {games.map((game) => {
            const start = new Date(game.CreatedAt);
            const end = game.FinishedAt ? new Date(game.FinishedAt) : null;
            const durationText = formatDuration(start, end);
            const winner = game.Winner?.Name;

            return (
              <li className="list-group-item" key={game.ID}>
                <div className="mb-2">
                  <strong>Game #{game.ID}</strong> –{" "}
                  <span>{game.GamePlayers?.length || 0} players</span> –{" "}
                  <span>Round {game.CurrentRound || "?"}</span>
                </div>

                <div className="text-muted mb-2">Length: {durationText}</div>

                {winner && (
                  <div className="mb-2">
                    <strong>Winner:</strong>{" "}
                    <span className="text-success fw-bold">{winner}</span>
                  </div>
                )}

                {game.Scores?.length > 0 && (
                  <div className="mb-2">
                    <strong>Scores:</strong>
                    <ul className="list-unstyled ms-3">
                      {game.Scores.map((s) => (
                        <li key={s.player_id}>
                          {s.player_name}: {s.points} pts
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                {game.GamePlayers && (
                  <div className="mb-2">
                    <strong>Players:</strong>
                    <ul className="list-unstyled ms-3">
                      {game.GamePlayers.map((gp) => (
                        <li key={gp.ID}>
                          {gp.Player?.Name || "Unnamed"} <em>({gp.Faction})</em>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}

export default GameList;
