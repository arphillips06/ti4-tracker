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

export default function ObjectiveFrequencyChart({ frequency }) {
  const [showAll, setShowAll] = useState(false);
  const [showTable, setShowTable] = useState(false);

  if (!frequency || typeof frequency !== "object") {
    return <div>No frequency data available.</div>;
  }

  const sortedObjectives = Object.entries(frequency).sort((a, b) => b[1] - a[1]);
  const displayData = showAll ? sortedObjectives : sortedObjectives.slice(0, 10);

  const chartData = {
    labels: displayData.map(([name]) => name),
    datasets: [
      {
        label: "Times Scored",
        data: displayData.map(([, count]) => count),
        backgroundColor: "rgba(54, 162, 235, 0.7)",
      },
    ],
  };

  const options = {
    indexAxis: "y",
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: { mode: "index", intersect: false },
    },
    scales: {
      x: {
        beginAtZero: true,
        ticks: { precision: 0 },
      },
      y: {
        ticks: { font: { size: 10 } },
      },
    },
  };

  return (
    <div style={{ marginBottom: "2rem" }}>
      <h3>Objective Scoring Frequency</h3>

      <div
        style={{
          height: showAll ? "600px" : "250px",
          overflowY: showAll ? "scroll" : "hidden",
          transition: "height 0.3s ease",
        }}
      >
        <Bar data={chartData} options={options} />
      </div>

      <div style={{ marginTop: "0.5rem", display: "flex", gap: "0.5rem" }}>
        <button className="btn btn-sm btn-outline-primary" onClick={() => setShowAll(!showAll)}>
          {showAll ? "Show Top 10 Only" : "Show All"}
        </button>
        <button className="btn btn-sm btn-outline-secondary" onClick={() => setShowTable(!showTable)}>
          {showTable ? "Hide Raw Data" : "Show Raw Data"}
        </button>
      </div>

      {showTable && (
        <table className="table table-sm table-bordered mt-3">
          <thead>
            <tr>
              <th>Objective</th>
              <th>Times Scored</th>
            </tr>
          </thead>
          <tbody>
            {sortedObjectives.map(([name, count]) => (
              <tr key={name}>
                <td>{name}</td>
                <td>{count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
