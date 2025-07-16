import React, { useState } from "react";
import { Bar } from "react-chartjs-2";
import "./shared/graphs.css"
import { sortData, horizontalBarOptions } from "./shared/chartUtils";

export default function FactionPlayWinChart({ data }) {
  const [showTable, setShowTable] = useState(false);
  const [sortKey, setSortKey] = useState("playRate");
  const [sortOrder, setSortOrder] = useState("desc");

  const setSort = (key) => {
    if (key === sortKey) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortOrder("desc");
    }
  };

  const processed = Object.entries(data).map(([name, val]) => ({
    name,
    ...val,
  }));

  const sorted = sortData(processed, sortKey, sortOrder);

  const chartData = {
    labels: sorted.map((f) => f.name),
    datasets: [
      {
        label: "% Played",
        data: sorted.map((f) => f.playRate),
        backgroundColor: "rgba(54, 162, 235, 0.7)",
      },
      {
        label: "% Won",
        data: sorted.map((f) => f.winRate),
        backgroundColor: "rgba(255, 99, 132, 0.7)",
      },
    ],
  };

  const options = horizontalBarOptions();

  return (
    <div className="graph-container">

      <div className="graph-bar-container-large">
        <Bar data={chartData} options={options} />
      </div>
      <div className="graph-toggle-buttons">
        <button
          className="btn btn-sm btn-outline-secondary graph-button-sm"
          onClick={() => setShowTable(!showTable)}
        >
          {showTable ? "Hide Raw Data" : "Show Raw Data"}
        </button>
      </div>
      {showTable && (
        <table className="table table-sm table-bordered graph-table">
          <thead>
            <tr>
              <th onClick={() => setSort("name")}>
                Faction {sortKey === "name" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("playedCount")}>
                Played {sortKey === "playedCount" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("winCount")}>
                Won {sortKey === "winCount" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("playRate")}>
                % Played {sortKey === "playRate" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("winRate")}>
                % Won {sortKey === "winRate" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((faction) => (
              <tr key={faction.name}>
                <td>{faction.name}</td>
                <td>{faction.playedCount}</td>
                <td>{faction.winCount}</td>
                <td>{faction.playRate.toFixed(1)}%</td>
                <td>{faction.winRate.toFixed(1)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
