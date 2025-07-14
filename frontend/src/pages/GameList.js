import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import API_BASE_URL from "../config";
import './gamelist.css';

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
            const scores = await fetchScoresForGame(game.id);
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
      <div className="d-flex justify-content-end gap-3 mb-4">
        <Link to="/" className="btn btn-outline-light">Home</Link>
        <Link to="/stats" className="btn btn-outline-light">Stats</Link>
      </div>

      {games.length === 0 ? (
        <p>No games found.</p>
      ) : (
        <ul className="list-group">
          {games.map((game) => {
            const start = new Date(game.created_at);
            const end = game.finished_at ? new Date(game.finished_at) : null;
            const durationText = formatDuration(start, end);
            const winner = game.winner?.Name;

            return (
              <li className="card-glass mb-4 p-3" key={game.id}>
                <div className="mb-2">
                  <strong>Game #{game.id}</strong> –{" "}
                  <span>{game.players?.length || 0} players</span> –{" "}
                  <span>Round {game.current_round ?? "?"}</span>
                </div>

                <div className="mb-2"><strong>Length:</strong> {durationText}</div>


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

                {game.players && (
                  <div className="mb-2">
                    <strong>Players:</strong>
                    <ul className="list-unstyled ms-3">
                      {game.players.map((gp) => (
                        <li key={gp.ID}>
                          {gp.Player?.Name || "Unnamed"} <em>({gp.Faction})</em>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                <div className="text-end">
                  <Link to={`/games/${game.id}`} className="btn btn-sm btn-primary">
                    View Game
                  </Link>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}

export default GameList;
