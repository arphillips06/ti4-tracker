import { color } from "chart.js/helpers";
import React from "react";
import { Bar } from "react-chartjs-2";

export default function GameLengthChart({ lengthData }) {
  const labels = Object.keys(lengthData).sort((a, b) => parseInt(a) - parseInt(b));
  const data = labels.map(key => lengthData[key]);

  const chartData = {
    labels,
    datasets: [
      {
        label: "Games Ending on Round",
        data,
        backgroundColor: "rgba(75, 192, 192, 0.7)",
      },
    ],
  };

  const options = {
    responsive: true,
    plugins: {
      legend: { display: false },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: { stepSize: 1 },
        ticks: {
            color: 'white'
        }
      },
      x: {
        ticks: {
            color: 'white'
        }
      }
    },
  };

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Game Length Distribution</h3>
      <Bar data={chartData} options={options} />
    </div>
  );
}
