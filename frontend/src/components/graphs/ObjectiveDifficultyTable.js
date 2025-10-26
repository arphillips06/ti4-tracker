import React, { useEffect, useMemo, useState } from "react";
import API_BASE_URL from "../../config";
import "./shared/graphs.css";
import "../../pages/stats.css";

export default function ObjectiveDifficultyTable() {
    const [rows, setRows] = useState([]);
    const [loading, setLoading] = useState(true);
    const [err, setErr] = useState(null);

    // controls to match your style but keep it simple
    const [stage, setStage] = useState("I");           // "I" | "II" | "Secret" | "all"
    const [minAppearances, setMinAppearances] = useState(3);

    // sorting like your SecretObjectiveTable
    const [sortKey, setSortKey] = useState("difficulty"); // hardest first by default
    const [sortOrder, setSortOrder] = useState("desc");

    useEffect(() => {
        let ignore = false;
        async function load() {
            setLoading(true);
            setErr(null);
            try {
                const url = `${API_BASE_URL}/stats/objectives/difficulty?stage=${encodeURIComponent(stage)}&minAppearances=${minAppearances}`;
                const res = await fetch(url, { credentials: "include" });
                if (!res.ok) throw new Error(`HTTP ${res.status}`);
                const data = await res.json();
                if (!ignore) setRows(Array.isArray(data?.rows) ? data.rows : []);
            } catch (e) {
                if (!ignore) setErr(e);
            } finally {
                if (!ignore) setLoading(false);
            }
        }
        load();
        return () => { ignore = true; };
    }, [stage, minAppearances]);

    const setSort = (key) => {
        if (key === sortKey) {
            setSortOrder(sortOrder === "asc" ? "desc" : "asc");
        } else {
            setSortKey(key);
            setSortOrder("desc");
        }
    };

    const sorted = useMemo(() => {
        const copy = [...rows];
        copy.sort((a, b) => {
            const A = a?.[sortKey] ?? 0;
            const B = b?.[sortKey] ?? 0;
            const cmp = A < B ? -1 : A > B ? 1 : 0;
            return sortOrder === "asc" ? cmp : -cmp;
        });
        return copy;
    }, [rows, sortKey, sortOrder]);

    const pct = (x) => `${((x ?? 0) * 100).toFixed(1)}%`;
    const num = (x) => (x ?? 0);
    const roundFmt = (x) => (x && x > 0 ? x.toFixed(1) : "—");
    const ci = (lo, hi) => `${((lo ?? 0) * 100).toFixed(1)}–${((hi ?? 0) * 100).toFixed(1)}%`;

    return (
        <div className="graph-container">
            <h3 className="chart-section-title">Objective Difficulty</h3>

            {/* simple controls styled like your pages */}
            <div className="stats-controls mb-2 d-flex gap-2 align-items-center">
                <label className="form-label mb-0">Stage:</label>
                <select
                    className="form-select form-select-sm"
                    style={{ width: 140 }}
                    value={stage}
                    onChange={(e) => setStage(e.target.value)}
                >
                    <option value="I">Stage I</option>
                    <option value="II">Stage II</option>
                    <option value="Secret">Secret</option>
                    <option value="all">All Public (I+II)</option>
                </select>

                <label className="form-label mb-0">Min Appearances:</label>
                <input
                    type="number"
                    className="form-control form-control-sm"
                    style={{ width: 100 }}
                    min="0"
                    value={minAppearances}
                    onChange={(e) => setMinAppearances(parseInt(e.target.value || "0", 10))}
                />
            </div>

            <div className="stats-section">
                {loading && <div>Loading…</div>}
                {err && <div className="text-danger">Error: {String(err)}</div>}

                {!loading && !err && (
                    <div className="table-scroll">
                        <table className="stats-table">
                            <thead>
                                <tr>
                                    <th onClick={() => setSort("name")} title="Objective card name">
                                        Objective {sortKey === "name" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-stage" onClick={() => setSort("stage")} title="I / II / Secret">
                                        Stage {sortKey === "stage" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-games" onClick={() => setSort("appearances")} title="# of games where this objective was revealed">
                                        Games Seen {sortKey === "appearances" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-chances" onClick={() => setSort("opportunities")} title="Sum of players in those games (one chance per player per game)">
                                        Player Chances {sortKey === "opportunities" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-scores" onClick={() => setSort("scores")} title="Unique player–game scores for this objective">
                                        Times Scored {sortKey === "scores" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-adj" onClick={() => setSort("adjRate")} title="Smoothed success rate = (Scores+2)/(Chances+7)">
                                        Adj. Success {sortKey === "adjRate" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-diff" onClick={() => setSort("difficulty")} title="Difficulty = 1 − Adj. Success (higher = harder)">
                                        Difficulty Score {sortKey === "difficulty" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-avg col--avgRound" onClick={() => setSort("avgRound")} title="Average first round when it was scored">
                                        Avg First-Round {sortKey === "avgRound" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-median col--medianRound" onClick={() => setSort("medianRound")} title="Median first round when it was scored">
                                        Median First-Round {sortKey === "medianRound" && (sortOrder === "asc" ? "▲" : "▼")}
                                    </th>
                                    <th className="num col-ci" title="95% confidence interval for the raw success rate">
                                        95% CI (Success)
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                {sorted.map((r) => (
                                    <tr key={`${r.objectiveId}-${r.name}`}>
                                        <td>{r.name}</td>
                                        <td className="num">{r.stage}</td>
                                        <td className="num">{r.appearances ?? 0}</td>
                                        <td className="num">{r.opportunities ?? 0}</td>
                                        <td className="num">{r.scores ?? 0}</td>
                                        <td className="num" title={`Raw: ${(r.rawRate * 100).toFixed(1)}%`}>
                                            {(r.adjRate * 100).toFixed(1)}%
                                        </td>
                                        <td className="num">{(r.difficulty * 100).toFixed(1)}%</td>
                                        <td className="num col--avgRound">
                                            {r.avgRound && r.avgRound > 0 ? r.avgRound.toFixed(1) : "—"}
                                        </td>
                                        <td className="num col--medianRound">
                                            {r.medianRound && r.medianRound > 0 ? r.medianRound.toFixed(1) : "—"}
                                        </td>
                                        <td className="num col-ci">
                                            {(r.wilsonLo * 100).toFixed(1)}–{(r.wilsonHi * 100).toFixed(1)}%
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
                {!loading && !err && rows.length === 0 && (
                    <div className="mt-2">No difficulty data available.</div>
                )}
            </div>
        </div>
    );
}
