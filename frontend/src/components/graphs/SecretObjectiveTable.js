import React, { useState } from "react";
import "./shared/graphs.css";
import "../../pages/stats.css";

export default function SecretObjectiveTable({ secrets }) {
  const [sortKey, setSortKey] = useState("timesScored");
  const [sortOrder, setSortOrder] = useState("desc");

  if (!Array.isArray(secrets) || secrets.length === 0) {
    return <div>No secret objective data available.</div>;
  }

  const setSort = (key) => {
    if (key === sortKey) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortOrder("desc");
    }
  };

  const sorted = [...secrets]
    .filter((obj) => obj.type === "Secret")
    .sort((a, b) => {
      const valA = a[sortKey] ?? 0;
      const valB = b[sortKey] ?? 0;
      return sortOrder === "asc" ? valA - valB : valB - valA;
    });

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Secret Objective Scoring</h3>
      <div className="stats-section">
        <table className="stats-table">
          <thead>
            <tr>
              <th onClick={() => setSort("name")}>
                Objective {sortKey === "name" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("timesScored")}>
                Scored {sortKey === "timesScored" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("averageRound")}>
                Avg Round {sortKey === "averageRound" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((obj) => (
              <tr key={obj.name}>
                <td>{obj.name}</td>
                <td style={{ textAlign: "center" }}>{obj.timesScored}</td>
                <td style={{ textAlign: "center" }}>{obj.averageRound.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
