import React, { useState } from "react";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend);

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

  const sorted = [...processed].sort((a, b) => {
    const valA = a[sortKey];
    const valB = b[sortKey];
    if (typeof valA === "string") {
      return sortOrder === "asc"
        ? valA.localeCompare(valB)
        : valB.localeCompare(valA);
    }
    return sortOrder === "asc" ? valA - valB : valB - valA;
  });

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

  const options = {
    indexAxis: "y",
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: "top" },
      tooltip: { mode: "nearest", intersect: true },
    },
    scales: {
      x: {
        beginAtZero: true,
        max: 100,
        ticks: { callback: (v) => `${v}%` },
      },
    },
  };

  return (
    <div className="mt-4">
      <div style={{ height: "500px" }}>
        <Bar data={chartData} options={options} />
      </div>

      <div className="mt-3">
        <button
          className="btn btn-sm btn-outline-secondary"
          onClick={() => setShowTable(!showTable)}
        >
          {showTable ? "Hide Raw Data" : "Show Raw Data"}
        </button>
      </div>

      {showTable && (
        <table className="table table-sm table-bordered mt-3">
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
                <td style={{ textAlign: "center" }}>{faction.playedCount}</td>
                <td style={{ textAlign: "center" }}>{faction.winCount}</td>
                <td style={{ textAlign: "center" }}>{faction.playRate.toFixed(1)}%</td>
                <td style={{ textAlign: "center" }}>{faction.winRate.toFixed(1)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
