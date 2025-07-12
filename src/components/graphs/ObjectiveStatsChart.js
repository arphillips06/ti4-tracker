// src/components/stats/ObjectiveStatsChart.js
import React from "react";
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

export default function ObjectiveStatsChart({ stats }) {
  const data = {
    labels: ["Public", "Secret", "Stage I", "Stage II"],
    datasets: [
      {
        label: "Objectives Scored",
        data: [
          stats.objectiveStats.publicScored,
          stats.objectiveStats.secretScored,
          stats.objectiveStats.stage1Scored,
          stats.objectiveStats.stage2Scored,
        ],
        backgroundColor: ["#36A2EB", "#FF6384", "#FFCE56", "#4BC0C0"],
      },
    ],
  };

  const options = {
    responsive: true,
    plugins: {
      legend: { position: "top" },
      tooltip: { mode: "index", intersect: false },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: { precision: 0 },
      },
    },
  };

  return (
    <div className="mb-4">
      <h4>Objective Types Scored</h4>
      <Bar data={data} options={options} />
    </div>
  );
}
