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
import '../../pages/stats.css';

ChartJS.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend);

export default function ObjectiveStatsChart({ stats }) {
  const totalPublic = stats.objectiveStats.publicScored || 0;
  const cdlPromoted = stats.objectiveStats.cdlPromoted || 0;
  const normalPublic = totalPublic - cdlPromoted;

const data = {
  labels: ["Public", "Secret", "Stage I", "Stage II"],
  datasets: [
    {
      label: "Public Objectives",
      data: [stats.objectiveStats.publicScored, 0, 0, 0],
      backgroundColor: "#36A2EB",
      stack: "objectiveStack",
      barPercentage: 1.0,
      categoryPercentage: 0.5,
    },
    {
      label: "Secrets Made Public",
      data: [cdlPromoted, 0, 0, 0],
      backgroundColor: "#B28DFF",
      stack: "objectiveStack",
      barPercentage: 1.0,
      categoryPercentage: 0.5,
    },
    {
      label: "Secret Objectives",
      data: [0, stats.objectiveStats.secretScored, 0, 0],
      backgroundColor: "#FF6384",
      stack: "objectiveStack",
      barPercentage: 1.0,
      categoryPercentage: 0.5,
    },
    {
      label: "Stage I",
      data: [0, 0, stats.objectiveStats.stage1Scored, 0],
      backgroundColor: "#FFCE56",
      stack: "objectiveStack",
      barPercentage: 1.0,
      categoryPercentage: 0.5,
    },
    {
      label: "Stage II",
      data: [0, 0, 0, stats.objectiveStats.stage2Scored],
      backgroundColor: "#4BC0C0",
      stack: "objectiveStack",
      barPercentage: 1.0,
      categoryPercentage: 0.5,
    },
  ],
};

const options = {
  responsive: true,
  plugins: {
    legend: { position: "top" },
    tooltip: { mode: "nearest", intersect: true },
  },
  scales: {
    x: {
      stacked: true,
    },
    y: {
      stacked: true,
      beginAtZero: true,
      ticks: { precision: 0 },
    },
  },
};



  return (
    <div className="raw-data-list mt-2 mb-2">
      <h4>Objective Types Scored</h4>
      <Bar data={data} options={options} />
    </div>
  );
}
