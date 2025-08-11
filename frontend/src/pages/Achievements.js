
import { useEffect, useMemo, useState } from "react";
import API_BASE_URL from "../config";
import "./achievements.css";

/**
 * Renders records from GET /achievements
 * Also (best-effort) GET /players to resolve holder names.
 */
export default function Achievements() {
  const [items, setItems] = useState([]);
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState();
  const [query, setQuery] = useState("");
  const [sortBy, setSortBy] = useState("label"); // label | value | holders
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

        // players is optional
        try {
          const pr = await fetch(`${API_BASE_URL}/players`, { signal: ctrl.signal });
          if (pr.ok) {
            const pj = await pr.json();
            const arr = Array.isArray(pj) ? pj : (Array.isArray(pj?.players) ? pj.players : []);
            setPlayers(arr);
          }
        } catch (_) {}
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
    if (sortBy === "label") list.sort((a,b)=>a.label.localeCompare(b.label));
    if (sortBy === "value") list.sort((a,b)=>(b.value ?? 0) - (a.value ?? 0));
    if (sortBy === "holders") list.sort((a,b)=>(b.holders?.length ?? 0) - (a.holders?.length ?? 0));
    return list;
  }, [items, query, sortBy]);

  return (
    <div className="ach-page">
      <div className="ach-page__inner">
        <header className="ach-head">
          <div>
            <h1 className="ach-title">Achievements</h1>
            <p className="ach-sub">Records and milestones pulled from <code>/achievements</code>.</p>
          </div>
          <div className="ach-ctrls">
            <input
              className="ach-input"
              placeholder="Search achievements…"
              value={query}
              onChange={(e)=>setQuery(e.target.value)}
            />
            <label className="ach-field">
              <span className="ach-label">Sort</span>
              <select className="ach-select" value={sortBy} onChange={(e)=>setSortBy(e.target.value)}>
                <option value="label">By name</option>
                <option value="value">By value</option>
                <option value="holders">By # holders</option>
              </select>
            </label>
          </div>
        </header>

        {loading && <div className="ach-skel" />}
        {error && <div className="ach-error">Error: {String(error)}</div>}

        {!loading && !error && (
          <section className="ach-grid">
            {filtered.map((a) => (
              <article className="ach-card" key={a.key}>
                <div className="ach-top">
                  <div className="ach-icon" aria-hidden>{iconFor(a.key)}</div>
                  <div>
                    <h3 className="ach-card__title">{a.label}</h3>
                    <div className="ach-tags">
                      <span className={`ach-tag ${a.status === 'record' ? 'ach-tag--ok' : ''}`}>{(a.status||'STATUS').toUpperCase()}</span>
                      <span className="ach-tag ach-tag--warn">Value: {a.value}</span>
                    </div>
                  </div>
                </div>
                <div className="ach-card__body">Key: <code>{a.key}</code></div>
                <div className="ach-foot">
                  <small className="ach-dim">Holders: {a.holders?.length ?? 0}</small>
                  <button className="ach-btn" onClick={()=>setActive(a)}>Details</button>
                </div>
              </article>
            ))}
          </section>
        )}
      </div>

      {active && (
        <div className="ach-modal__backdrop" onClick={()=>setActive(null)}>
          <div className="ach-modal" onClick={(e)=>e.stopPropagation()}>
            <div className="ach-top">
              <div className="ach-icon" aria-hidden>{iconFor(active.key)}</div>
              <div>
                <h3 className="ach-card__title">{active.label}</h3>
                <div className="ach-sub">{(active.status||'').toUpperCase()}</div>
              </div>
              <button className="ach-btn ach-btn--ghost" onClick={()=>setActive(null)}>Close</button>
            </div>
            <div className="ach-modal__body">
              <h4>Current Value</h4>
              <p className="ach-kpi">{active.value}</p>
              <h4>Holders</h4>
              <ul className="ach-list">
                {(active.holders ?? []).map((h, idx) => (
                  <li key={idx} className="ach-list__item">
                    <span className="ach-list__title">{playerNameById(h.player_id)}</span>
                    <span className="ach-list__meta">Game #{h.game_id}{h.round_id ? ` · Round #${h.round_id}` : ''}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function iconFor(key) {
  if (/fast|speed/i.test(key)) return "⚡";
  if (/round/i.test(key)) return "🌀";
  if (/point/i.test(key)) return "⭐";
  if (/win/i.test(key)) return "🏆";
  return "🎯";
}
