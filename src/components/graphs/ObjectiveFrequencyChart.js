import { Bar } from "react-chartjs-2";
import React, { useState } from "react";

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

export default function ObjectiveFrequencyChart({ frequency, secretPublic = 0 }) {
  const [showAll, setShowAll] = useState(false);
  const [showTable, setShowTable] = useState(false);

  if (!frequency || typeof frequency !== "object") {
    return <div>No frequency data available.</div>;
  }

  const secretPublicBreakdown = {}; // For future use

  const sortedObjectives = Object.entries(frequency).sort((a, b) => b[1] - a[1]);
  const displayData = showAll ? sortedObjectives : sortedObjectives.slice(0, 10);

  const chartData = {
    labels: displayData.map(([name]) => name),
    datasets: [
      {
        label: "Public Objectives",
        data: displayData.map(([name]) => {
          const count = frequency[name];
          const secretAmount = secretPublicBreakdown?.[name] || 0;
          return count - secretAmount;
        }),
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
      padding: {
        left: 0, // Gives room for long Y-axis labels
      },
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
    <div className="p-3">
      <h2 className="mb-2" style={{ color: "#f4d35e" }}>Objective Appearance Frequency</h2>
      <div style={{ height: `${displayData.length * 35}px` }}>
        <Bar data={chartData} options={options} />
      </div>

      <div className="mt-3 d-flex gap-2">
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
    </div>
  );
}
