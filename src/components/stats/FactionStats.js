import React, { useState } from "react";
import FactionPlayWinChart from "../graphs/FactionPlayWinChart";
import FactionWinRateChart from "../graphs/FactionWinRateChart";
import '../../pages/stats.css';
import FactionUsageTable from "../graphs/FactionUsageTable";


export default function FactionStats({ stats }) {
  const [showRaw, setShowRaw] = useState(false);
  const [showFactionData, setShowFactionData] = useState(false);
  if (!stats.factionPlayWinDistribution) return <div>No faction data available.</div>;


  const entries = Object.entries(stats.factionPlayWinDistribution);
  const sorted = [...entries].sort((a, b) => b[1].winRate - a[1].winRate); // sort by win rate descending

  return (
    <>
      {/* Faction Win Rate Chart */}
      <div className="chart-glass-container">
        <FactionWinRateChart dataMap={stats.winRateByFaction} />
        <button
          onClick={() => setShowFactionData(!showFactionData)}
          className="btn btn-sm btn-outline-secondary mt-2"
        >
          {showFactionData ? "Hide Raw Data" : "Show Raw Data"}
        </button>
        {showFactionData && (
          <ul className="raw-data-list mt-2 mb-2">
            {Object.entries(stats.winRateByFaction).map(([faction, rate]) => (
              <li key={faction}>
                {faction}: {rate.toFixed(2)}%
              </li>
            ))}
          </ul>
        )}
      </div>
      <div className="chart-glass-container">
        <h3 className="chart-title">Faction Play vs Win Rates</h3>
        <FactionPlayWinChart data={stats.factionPlayWinDistribution} />
      </div>
      <div className="chart-glass-container">
        <h3 className="chart-title">Faction Usage</h3>
        <FactionUsageTable
          data={stats.factionPlayerStats}
          aggregates={stats.factionAggregateStats} />
      </div>
    </>
  );
}
