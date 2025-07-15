import React, { useState } from "react";
import { Bar } from "react-chartjs-2";
import "./shared/graphs.css";
import { sortData } from "./shared/chartUtils";

export default function ObjectiveFrequencyChart({ frequency, secretPublic = 0 }) {
  const [showAll, setShowAll] = useState(false);
  const [showTable, setShowTable] = useState(false);

  if (!frequency || typeof frequency !== "object") {
    return <div>No frequency data available.</div>;
  }

  const secretPublicBreakdown = {};

  const objectiveData = Object.entries(frequency).map(([name, count]) => ({
    name,
    publicCount: count - (secretPublicBreakdown[name] || 0),
    total: count,
  }));

  const sorted = sortData(objectiveData, "total", "desc");
  const displayData = showAll ? sorted : sorted.slice(0, 10);

  const chartData = {
    labels: displayData.map((o) => o.name),
    datasets: [
      {
        label: "Public Objectives",
        data: displayData.map((o) => o.publicCount),
        backgroundColor: "rgba(54, 162, 235, 0.7)",
        stack: "stack1",
      },
    ],
  };

  const options = {
    indexAxis: "y",
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      tooltip: { mode: "nearest", intersect: true },
      legend: { position: "top" },
    },
    layout: {
      padding: { left: 0 }, 
    },
    scales: {
      x: {
        beginAtZero: true,
        stacked: true,
        ticks: { precision: 0 },
      },
      y: {
        stacked: true,
        ticks: {
          font: { size: 11 },
          autoSkip: false,
          padding: 8,
        },
      },
    },
  };

  return (
    <div className="graph-container">
      <h2 className="chart-section-title">Objective Appearance Frequency</h2>

      <div
        className="graph-bar-container"
        style={{ height: `${displayData.length * 35}px` }}
      >
        <Bar data={chartData} options={options} />
      </div>

      <div className="graph-button-row">
        <button
          onClick={() => setShowAll(!showAll)}
          className="btn btn-sm btn-outline-secondary"
        >
          {showAll ? "Show Top 10" : "Show All"}
        </button>
        <button
          onClick={() => setShowTable(!showTable)}
          className="btn btn-sm btn-outline-secondary"
        >
          {showTable ? "Hide Raw Data" : "Show Raw Data"}
        </button>
      </div>

      {showTable && (
        <table className="table table-sm table-bordered mt-2">
          <thead>
            <tr>
              <th>Objective</th>
              <th>Times Appeared</th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((row) => (
              <tr key={row.name}>
                <td>{row.name}</td>
                <td>{row.total}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
