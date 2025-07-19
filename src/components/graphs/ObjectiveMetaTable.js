import React, { useState } from "react";
import "./shared/graphs.css";
import "../../pages/stats.css";

export default function ObjectiveMetaTable({ metas }) {
  const [sortKey, setSortKey] = useState("timesScored");
  const [sortOrder, setSortOrder] = useState("desc");
  console.log("metas received:", metas);

  if (!metas || !Array.isArray(metas) || metas.length === 0) {
    return <div>No objective metadata available.</div>;
  }

  const setSort = (key) => {
    if (key === sortKey) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortOrder("desc");
    }
  };

  const sorted = [...metas]
    .filter((obj) => obj.type !== "Secret")
    .sort((a, b) => {
      const valA = a[sortKey] ?? 0;
      const valB = b[sortKey] ?? 0;
      return sortOrder === "asc" ? valA - valB : valB - valA;
    });

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Objective Scoring Insights</h3>
      <div className="stats-section">
        <table className="stats-table">
          <thead>
            <tr>
              <th onClick={() => setSort("name")}>Objective {sortKey === "name" && (sortOrder === "asc" ? "▲" : "▼")}</th>
              <th onClick={() => setSort("type")}>Type {sortKey === "type" && (sortOrder === "asc" ? "▲" : "▼")}</th>
              <th onClick={() => setSort("timesAppeared")}>Appeared {sortKey === "timesAppeared" && (sortOrder === "asc" ? "▲" : "▼")}</th>
              <th onClick={() => setSort("timesScored")}>Scored {sortKey === "timesScored" && (sortOrder === "asc" ? "▲" : "▼")}</th>
              <th onClick={() => setSort("scoredPercent")}>Scored % {sortKey === "scoredPercent" && (sortOrder === "asc" ? "▲" : "▼")}</th>
              <th onClick={() => setSort("averageRound")}>Avg Round {sortKey === "averageRound" && (sortOrder === "asc" ? "▲" : "▼")}</th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((obj) => (
              <tr key={obj.name}>
                <td>{obj.name}</td>
                <td>{obj.type}</td>
                <td style={{ textAlign: "center" }}>{obj.timesAppeared}</td>
                <td style={{ textAlign: "center" }}>{obj.timesScored}</td>
                <td style={{ textAlign: "center" }}>{obj.scoredPercent.toFixed(2)}%</td>
                <td style={{ textAlign: "center" }}>{obj.averageRound.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
