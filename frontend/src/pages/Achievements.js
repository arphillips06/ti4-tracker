import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import API_BASE_URL from "../config";
import "./stats.css";
import "./achievements.css";

const badgeDescriptions = {
  record_comeback_kid:
    "Largest mid-game deficit (points behind the leader) overcome by the eventual winner.",
  record_largest_margin:
    "Largest final victory margin (winner vs runner-up).",
  most_points_in_round:
    "Most points scored by a player in a single round.",
};

function formatBadgeValue(key, value) {
  if (value == null) return "";
  switch (key) {
    case "record_comeback_kid":
      return `${value} point${value === 1 ? "" : "s"} behind`;
    case "record_largest_margin":
      return `${value} point${value === 1 ? "" : "s"}`;
    case "most_points_in_round":
      return `${value} point${value === 1 ? "" : "s"}`;
    default:
      return String(value);
  }
}

function roundHeaderFor(key) {
  switch (key) {
    case "record_comeback_kid":
    case "most_points_in_round":
      return "Round (in-game #)";
    default:
      return "Round";
  }
}

export default function Achievements() {
  const [items, setItems] = useState([]);
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState();
  const [query, setQuery] = useState("");
  const [sortBy, setSortBy] = useState("label");
  const [active, setActive] = useState(null);

  useEffect(() => {
    const ctrl = new AbortController();
    (async () => {
      try {
        setLoading(true);
        setError(undefined);
        const res = await fetch(`${API_BASE_URL}/achievements`, { signal: ctrl.signal });
        if (!res.ok) throw new Error(`GET /achievements ${res.status}`);
        const json = await res.json();
        setItems(Array.isArray(json?.value) ? json.value : []);

        // Optional: get names for holders
        try {
          const pr = await fetch(`${API_BASE_URL}/players`, { signal: ctrl.signal });
          if (pr.ok) {
            const pj = await pr.json();
            const arr = Array.isArray(pj) ? pj : (Array.isArray(pj?.players) ? pj.players : []);
            setPlayers(arr);
          }
        } catch (_) { }
      } catch (e) {
        if (e.name !== "AbortError") setError(e?.message ?? "Failed to load achievements");
      } finally {
        setLoading(false);
      }
    })();
    return () => ctrl.abort();
  }, []);

  const playerNameById = useMemo(() => {
    const m = new Map();
    for (const p of players) {
      const id = p.id ?? p.ID ?? p.player_id ?? p.PlayerID;
      const name = p.name ?? p.Name ?? p.player_name ?? (id != null ? `Player #${id}` : "Unknown");
      if (id != null) m.set(Number(id), String(name));
    }
    return (id) => m.get(Number(id)) ?? `Player #${id}`;
  }, [players]);

  const filtered = useMemo(() => {
    let list = items.filter(a =>
      query ? (a.label + a.key + a.status).toLowerCase().includes(query.toLowerCase()) : true
    );
    if (sortBy === "label") list.sort((a, b) => a.label.localeCompare(b.label));
    if (sortBy === "value") list.sort((a, b) => (b.value ?? 0) - (a.value ?? 0));
    if (sortBy === "holders") list.sort((a, b) => (b.holders?.length ?? 0) - (a.holders?.length ?? 0));
    return list;
  }, [items, query, sortBy]);

  return (
    <div className="p-4">
      <h1 className="mb-4">Twilight Imperium Stats</h1>

      <div className="stats-nav mb-4">
        <Link to="/" className="nav-btn">Home</Link>
        <button className="active" disabled>Achievements</button>
      </div>

      <h2 className="section-title mb-3">Achievements</h2>

      <div className="ach-toolbar mb-3">
        <input
          className="ach-input"
          placeholder="Search achievements…"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <label className="ach-field">
          <span className="ach-label">Sort</span>
          <select className="ach-select" value={sortBy} onChange={(e) => setSortBy(e.target.value)}>
            <option value="label">By name</option>
            <option value="value">By value</option>
            <option value="holders">By # holders</option>
          </select>
        </label>
      </div>

      {loading && <div className="ach-skel" />}
      {error && <div className="ach-error alert alert-danger mt-3">Error: {String(error)}</div>}

      {!loading && !error && (
        <section className="ach-grid">
          {filtered.map((a) => (
            <article className="ach-card" key={a.key}>
              <div className="ach-top">
                <div>
                  <h3 className="ach-card__title">{a.label}</h3>
                  <div className="ach-tags">
                    <span className={`ach-tag ${a.status === 'record' ? 'ach-tag--ok' : ''}`}>
                      {(a.status || 'STATUS').toUpperCase()}
                    </span>
                    <span className="ach-tag ach-tag--warn">
                      Value: {formatBadgeValue(a.key, a.value)}
                    </span>
                  </div>
                </div>
              </div>
              <div className="ach-foot">
                <small className="ach-dim">Holders: {a.holders?.length ?? 0}</small>
                <button className="nav-btn" onClick={() => setActive(a)}>Details</button>
              </div>
            </article>
          ))}
        </section>
      )}

      {active && (
        <div className="ach-modal__backdrop" onClick={() => setActive(null)}>
          <div className="ach-modal" onClick={(e) => e.stopPropagation()}>
            <div className="d-flex align-items-center justify-content-between mb-3">
              <div>
                <h3 className="ach-modal__title m-0">{active.label}</h3>
                <div className="ach-modal__sub">{(active.status || '').toUpperCase()}</div>
                {badgeDescriptions[active.key] && (
                  <p className="text-muted mt-1 mb-0" style={{ maxWidth: 520 }}>
                    {badgeDescriptions[active.key]}
                  </p>
                )}
              </div>
              <button className="nav-btn" onClick={() => setActive(null)}>Close</button>
            </div>

            <div className="mb-3">
              <div className="text-gold fw-semibold">Current Value</div>
              <div className="ach-kpi">{formatBadgeValue(active.key, active.value)}</div>
            </div>

            <div className="text-gold fw-semibold mb-2">Holders</div>
            <div className="table-responsive">
              <table className="table table-dark table-striped table-bordered align-middle mb-0">
                <thead>
                  <tr>
                    <th style={{ width: '40%' }}>Player</th>
                    <th style={{ width: '30%' }}>Game</th>
                    <th style={{ width: '30%' }}>{roundHeaderFor(active.key)}</th>
                  </tr>
                </thead>
                <tbody>
                  {(active.holders ?? []).map((h, idx) => (
                    <tr key={idx}>
                      <td className="fw-semibold">{playerNameById(h.player_id)}</td>
                      <td>
                        <Link to={`/games/${h.game_id}`}>Game #{h.game_id}</Link>
                      </td>
                      <td>{h.round_id ? `Round ${h.round_id}` : "-"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
