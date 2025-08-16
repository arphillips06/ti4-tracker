import React, { useEffect, useState, useRef, useCallback } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import API_BASE_URL from "../config";
import './gamelist.css';

function useDebouncedValue(value, delay = 300) {
  const [v, setV] = useState(value);
  useEffect(() => {
    const id = setTimeout(() => setV(value), delay);
    return () => clearTimeout(id);
  }, [value, delay]);
  return v;
}

export default function GameList() {
  const [games, setGames] = useState([]);
  const [query, setQuery] = useState('');
  const debounced = useDebouncedValue(query, 300);
  const cancelRef = useRef(null);

  const formatDuration = (start, end) => {
    if (!end) return "Ongoing";
    const ms = end - start;
    const totalMinutes = Math.floor(ms / (1000 * 60));
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;
    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  };

  const fetchScoresForGame = async (id, cancelToken) => {
    try {
      const res = await axios.get(`${API_BASE_URL}/games/${id}/score-summary`, { cancelToken });
      return res.data;
    } catch (e) {
      if (axios.isCancel(e)) return [];
      console.error("Error fetching scores for game", id, e);
      return [];
    }
  };

  const loadGames = useCallback(async (term) => {
    if (cancelRef.current) cancelRef.current.cancel('New request starting');
    cancelRef.current = axios.CancelToken.source();

    try {
      const url = new URL(`${API_BASE_URL}/games`);
      if (term.trim()) url.searchParams.set('search', term.trim());
      const res = await axios.get(url.toString(), { cancelToken: cancelRef.current.token });
      const gamesData = res.data || [];

      const withScores = await Promise.all(
        gamesData.map(async (game) => {
          const scores = await fetchScoresForGame(game.id, cancelRef.current.token);
          return { ...game, Scores: scores };
        })
      );

      setGames(withScores);
    } catch (err) {
      if (!axios.isCancel(err)) console.error('Error fetching games:', err);
    }
  }, []);

  useEffect(() => {
    loadGames(debounced);
    return () => {
      if (cancelRef.current) cancelRef.current.cancel('Unmount/dependency change');
    };
  }, [debounced, loadGames]);

  const onKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      loadGames(query); // immediate fetch on Enter
    }
  };

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (query) params.set("search", query); else params.delete("search");
    const next = `${window.location.pathname}?${params.toString()}`;
    window.history.replaceState(null, "", next);
  }, [query]);

  return (
    <div className="container">
      <h1 className="mb-4">TI4 Games</h1>

      <div className="d-flex justify-content-between align-items-center gap-3 mb-4">
        <div className="search-wrap d-flex align-items-center" style={{ width: '100%', maxWidth: 720 }}>
          <input
            id="games-search"
            className="form-control search-input"
            placeholder="Search… (hover for help)"
            aria-describedby="games-search-help"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={onKeyDown}
          />
          <div id="games-search-help" className="help-popover" role="tooltip">
            <div className="mb-1"><strong>Supported:</strong></div>
            <div><code>winner:ross</code> · <code>player:alex</code> · <code>faction:nekro</code></div>
            <div><code>rounds&gt;=8</code> · <code>after:2025-07-01</code> · <code>before:2025-08-01</code></div>
            <div><code>agenda:"Seed of an Empire"</code> · <code>relic:shard</code> · <code>custodians:true</code></div>
            <div>Free text matches title/notes.</div>
            <div className="mt-1"><em>Press Enter to search immediately.</em></div>
          </div>
        </div>

        <div className="d-flex gap-3">
          <button
            className="nav-btn"
            type="button"
            onClick={() => setQuery('')}
            disabled={!query}
            title="Clear search"
          >
            Clear
          </button>
          <Link to="/" className="nav-btn">Home</Link>
          <Link to="/stats" className="nav-btn">Stats</Link>
        </div>
      </div>

      {games.length === 0 ? (
        <p>No games found.</p>
      ) : (
        <ul className="list-group list-unstyled">
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
                  <div className="stats-nav mb-4">
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
                  <Link to={`/games/${game.id}`} className="nav-btn">
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
