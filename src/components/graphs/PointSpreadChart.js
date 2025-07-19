import React from "react";
import { Bar } from "react-chartjs-2";
import { Chart, BarElement, CategoryScale, LinearScale, Tooltip, Legend, TimeScale } from "chart.js";
import { color } from "chart.js/helpers";

Chart.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend);

export default function PointSpreadChart({ spreadData }) {
  const labels = Object.keys(spreadData).sort((a, b) => parseInt(a) - parseInt(b));
  const data = labels.map(key => spreadData[key]);

  const chartData = {
    labels,
    datasets: [
      {
        label: "Games by Point Spread",
        data,
        backgroundColor: "rgba(255, 99, 132, 0.7)",
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
        ticks: { 
            stepSize: 1,
            color: 'white',
         },
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
      <h3 className="chart-section-title">Victory Point Spread</h3>
      <Bar data={chartData} options={options} />
    </div>
  );
}
