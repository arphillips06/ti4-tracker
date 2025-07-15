import React from "react";
import { Bar } from "react-chartjs-2";
import "./shared/graphs.css"; 

export default function ObjectiveStatsChart({ stats }) {
  const totalPublic = stats.objectiveStats.publicScored || 0;
  const cdlPromoted = stats.objectiveStats.cdlPromoted || 0;
  const normalPublic = totalPublic - cdlPromoted;

  const data = {
    labels: ["Public", "Secret", "Stage I", "Stage II"],
    datasets: [
      {
        label: "Public Objectives",
        data: [normalPublic, 0, 0, 0],
        backgroundColor: "#36A2EB",
        stack: "objectiveStack",
      },
      {
        label: "Secrets Made Public",
        data: [cdlPromoted, 0, 0, 0],
        backgroundColor: "#B28DFF",
        stack: "objectiveStack",
      },
      {
        label: "Secret Objectives",
        data: [0, stats.objectiveStats.secretScored, 0, 0],
        backgroundColor: "#FF6384",
        stack: "objectiveStack",
      },
      {
        label: "Stage I",
        data: [0, 0, stats.objectiveStats.stage1Scored, 0],
        backgroundColor: "#FFCE56",
        stack: "objectiveStack",
      },
      {
        label: "Stage II",
        data: [0, 0, 0, stats.objectiveStats.stage2Scored],
        backgroundColor: "#4BC0C0",
        stack: "objectiveStack",
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: "top" },
      tooltip: { mode: "nearest", intersect: true },
    },
    scales: {
      x: { stacked: true },
      y: {
        stacked: true,
        beginAtZero: true,
        ticks: { precision: 0 },
      },
    },
  };

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Objective Types Scored</h3>
      <div className="graph-bar-container" style={{ height: "300px" }}>
        <Bar data={data} options={options} />
      </div>
    </div>
  );
}
