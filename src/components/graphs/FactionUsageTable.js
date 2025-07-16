import React, { useState } from "react";
import "../../pages/stats.css";

export default function FactionUsageTable({ data }) {
  const [sortKey, setSortKey] = useState("faction");
  const [sortOrder, setSortOrder] = useState("asc");

  // Aggregate rows by faction
  const grouped = {};

  data.forEach((entry) => {
    const { faction, player, playedCount, wonCount } = entry;
    if (!grouped[faction]) {
      grouped[faction] = {
        faction,
        totalPlays: 0,
        playersPlayed: {},
        playersWon: {},
      };
    }

    grouped[faction].totalPlays += playedCount;
    if (playedCount > 0) {
      grouped[faction].playersPlayed[player] =
        (grouped[faction].playersPlayed[player] || 0) + playedCount;
    }
    if (wonCount > 0) {
      grouped[faction].playersWon[player] =
        (grouped[faction].playersWon[player] || 0) + wonCount;
    }
  });

  // Turn grouped data into array
  const tableData = Object.values(grouped);

  // Sort
  const sorted = [...tableData].sort((a, b) => {
    const valA = a[sortKey];
    const valB = b[sortKey];
    if (typeof valA === "string") {
      return sortOrder === "asc"
        ? valA.localeCompare(valB)
        : valB.localeCompare(valA);
    }
    return sortOrder === "asc" ? valA - valB : valB - valA;
  });

  const toggleSort = (key) => {
    if (key === sortKey) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortOrder("asc");
    }
  };

  return (
    <div className="stats-section">
      <h2>Faction Usage</h2>
      <table className="stats-table">
        <thead>
          <tr>
            <th onClick={() => toggleSort("faction")}>
              Faction {sortKey === "faction" && (sortOrder === "asc" ? "▲" : "▼")}
            </th>
            <th onClick={() => toggleSort("totalPlays")}>
              Total Plays {sortKey === "totalPlays" && (sortOrder === "asc" ? "▲" : "▼")}
            </th>
            <th>Players Played</th>
            <th>Players Won With</th>
          </tr>
        </thead>
        <tbody>
          {sorted.map((row) => (
            <tr key={row.faction}>
              <td>{row.faction}</td>
              <td>{row.totalPlays || "-"}</td>
              <td>
                {Object.keys(row.playersPlayed).length > 0
                  ? Object.entries(row.playersPlayed)
                      .map(([player, count]) => `${player} (${count})`)
                      .join(", ")
                  : "-"}
              </td>
              <td>
                {Object.keys(row.playersWon).length > 0
                  ? Object.entries(row.playersWon)
                      .map(([player, count]) => `${player} (${count})`)
                      .join(", ")
                  : "-"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
