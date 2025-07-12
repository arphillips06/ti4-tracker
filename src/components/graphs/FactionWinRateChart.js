import React from "react";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

// Register Chart.js components
ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

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
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
        title: { display: true, text: "Win Rate (%)" },
      },
    },
  };

  return <Bar data={data} options={options} />;
}
