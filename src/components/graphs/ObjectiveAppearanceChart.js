import React, { useState } from "react";
import { Bar } from "react-chartjs-2";
import "./shared/graphs.css";
import { sortData } from "./shared/chartUtils";
import { color } from "chart.js/helpers";

export default function ObjectiveAppearanceChart({ data  }) {
  const [showAll, setShowAll] = useState(false);
  const [showTable, setShowTable] = useState(false);
  const [sortKey, setSortKey] = useState("appearPct");
  const [sortOrder, setSortOrder] = useState("desc");

  const setSort = (key) => {
    if (key === sortKey) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortOrder("desc");
    }
  };

  const appearanceStats = data  || {};

  const merged = Object.entries(appearanceStats)
    .filter(([_, data]) => data.type !== "Secret")
    .map(([name, data]) => ({
      name,
      appearPct: data.appearanceRate || 0,
      scorePct: data.scoredWhenAppearedRate || 0,
      appeared: data.appearedCount || 0,
      scored: data.scoredCount || 0,
    }));
  const mergedMap = Object.fromEntries(merged.map(d => [d.name, d]));

  const sorted = sortData(merged, sortKey, sortOrder);
  const displayData = showAll ? sorted : sorted.slice(0, 10);

  const chartData = {
    labels: displayData.map((d) => d.name), // ✅ based on filtered `displayData`
    datasets: [
      {
        label: "% of Games Appeared",
        data: displayData.map((d) => d.appearPct),
        backgroundColor: "rgba(75, 192, 192, 0.6)",
      },
      {
        label: "% Scored When Appeared",
        data: displayData.map((d) => d.scorePct),
        backgroundColor: "rgba(255, 159, 64, 0.6)",
      },
    ],
  };


  const options = {
    indexAxis: "y",
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      tooltip: {
        mode: "nearest",
        intersect: true,
        callbacks: {
          afterLabel: (context) => {
            const stat = mergedMap[context.label]; // Only public objectives
            if (!stat) return "";
            if (context.dataset.label === "% of Games Appeared") {
              return `Appeared in ${stat.appeared} games`;
            } else if (context.dataset.label === "% Scored When Appeared") {
              return `Scored in ${stat.scored} of ${stat.appeared} games`;
            }
            return "";
          },
        },
      },
      legend: { position: "top" },
    },
    scales: {
      x: {
        beginAtZero: true,
        max: 100,
        ticks: {
          callback: (v) => `${v}%`,
          color: 'white',
        },
      },
      y: {
        ticks: {
          color: 'white',
        },
      }
    },
  };

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Objective Appearance vs Scoring</h3>

      <div
        className="graph-bar-container"
        style={{ height: `${displayData.length * 35}px` }}
      >
        <Bar data={chartData} options={options} />
      </div>

      <div className="graph-button-row">
        <button className="btn btn-sm btn-outline-primary" onClick={() => setShowAll(!showAll)}>
          {showAll ? "Show Top 10 Only" : "Show All"}
        </button>
        <button className="btn btn-sm btn-outline-secondary" onClick={() => setShowTable(!showTable)}>
          {showTable ? "Hide Objective Raw Data" : "Show Objective Raw Data"}
        </button>
      </div>

      {showTable && (
        <table className="table table-sm table-bordered mt-2 graph-table">
          <thead>
            <tr>
              <th onClick={() => setSort("name")}>
                Objective {sortKey === "name" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("appearPct")}>
                % Appeared {sortKey === "appearPct" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
              <th onClick={() => setSort("scorePct")}>
                % Scored When Appeared {sortKey === "scorePct" && (sortOrder === "asc" ? "▲" : "▼")}
              </th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((row) => (
              <tr key={row.name}>
                <td>{row.name}</td>
                <td style={{ textAlign: "center" }}>{row.appearPct.toFixed(1)}%</td>
                <td style={{ textAlign: "center" }}>{row.scorePct.toFixed(1)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}