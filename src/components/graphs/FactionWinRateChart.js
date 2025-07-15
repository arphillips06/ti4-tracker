import React from "react";
import { Bar } from "react-chartjs-2";
import "./shared/graphs.css";

export default function FactionWinRateChart({ dataMap }) {
  const factions = Object.keys(dataMap || {});
  const rates = Object.values(dataMap || {});

  const data = {
    labels: factions,
    datasets: [
      {
        label: "Win Rate (%)",
        data: rates,
        backgroundColor: "rgba(75,192,192,0.6)",
        borderColor: "rgba(75,192,192,1)",
        borderWidth: 1,
      },
    ],
  };

  const options = {
    responsive: true,
    plugins: {
      legend: { position: "top" },
      title: { display: true, text: "Faction Win Rate" },
      tooltip: { mode: "nearest", intersect: true },
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
        title: { display: true, text: "Win Rate (%)" },
      },
    },
  };

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Faction Win Rate</h3>
      <div className="graph-bar-container-medium">
        <Bar data={data} options={options} />
      </div>
    </div>
  );
}