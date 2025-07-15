// /components/graphs/shared/chartUtils.js
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

// Ensure all charts work with common components registered once
ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

// Standard sort helper for table and chart data
export const sortData = (data, sortKey, sortOrder = "desc") => {
  return [...data].sort((a, b) => {
    const valA = a[sortKey];
    const valB = b[sortKey];
    if (typeof valA === "string") {
      return sortOrder === "asc"
        ? valA.localeCompare(valB)
        : valB.localeCompare(valA);
    }
    return sortOrder === "asc" ? valA - valB : valB - valA;
  });
};

// Horizontal bar chart config with % axis
export const horizontalBarOptions = (xMax = 100) => ({
  indexAxis: "y",
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: "top" },
    tooltip: { mode: "nearest", intersect: true },
  },
  scales: {
    x: {
      beginAtZero: true,
      max: xMax,
      ticks: { callback: (v) => `${v}%` },
    },
  },
});
