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
export const horizontalBarOptions = () => ({
  indexAxis: 'y',
  responsive: true,
  plugins: {
    legend: {
      labels: {
        color: '#ffffff'  // White legend text
      }
    },
    title: {
      display: true,
      text: 'Faction Play vs Win Rates',
      color: '#ffffff',  // White title text
      font: {
        size: 20
      }
    },
    tooltip: {
      bodyColor: '#ffffff',
      titleColor: '#ffffff',
    }
  },
  scales: {
    x: {
      ticks: {
        color: '#ffffff'  // White x-axis labels
      },
      grid: {
        color: 'rgba(255,255,255,0.1)' // Optional: faint white gridlines
      }
    },
    y: {
      ticks: {
        color: '#ffffff'  // White y-axis labels (faction names)
      },
      grid: {
        color: 'rgba(255,255,255,0.1)'
      }
    }
  }
});

