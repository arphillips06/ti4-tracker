import React, { useState } from "react";
import FactionWinRateChart from "../graphs/FactionWinRateChart";
import ObjectiveStatsChart from "../graphs/ObjectiveStatsChart";
import ObjectiveFrequencyChart from "../graphs/ObjectiveFrequencyChart";
import ObjectiveAppearanceChart from "../graphs/ObjectiveAppearanceChart";

export default function OverviewStats({ stats }) {
  const [showFactionData, setShowFactionData] = useState(false);
  const [showObjectiveData, setShowObjectiveData] = useState(false);

  return (
    <div className="stats-section">
      <h2>Total Games: {stats.totalGames}</h2>

      <ObjectiveStatsChart stats={stats} />

      <button onClick={() => setShowObjectiveData(!showObjectiveData)} className="btn btn-sm btn-outline-secondary mb-2">
        {showObjectiveData ? "Hide Raw Data" : "Show Raw Data"}
      </button>

      {showObjectiveData && (
        <ul>
          <li>Public Objectives: {stats.objectiveStats.publicScored}</li>
          <li>Secret Objectives: {stats.objectiveStats.secretScored}</li>
          <li>Stage I: {stats.objectiveStats.stage1Scored}</li>
          <li>Stage II: {stats.objectiveStats.stage2Scored}</li>
        </ul>
      )}
      <ObjectiveFrequencyChart frequency={stats.objectiveFrequency} />

      <h3>Faction Win Rates</h3>
      <FactionWinRateChart dataMap={stats.winRateByFaction} />

      <button onClick={() => setShowFactionData(!showFactionData)} className="btn btn-sm btn-outline-secondary mt-2">
        {showFactionData ? "Hide Raw Data" : "Show Raw Data"}
      </button>
      {showFactionData && (
        <ul className="mt-2">
          {Object.entries(stats.winRateByFaction).map(([faction, rate]) => (
            <li key={faction}>
              {faction}: {rate.toFixed(2)}%
            </li>
          ))}
        </ul>
      )}
      <ObjectiveAppearanceChart stats={stats} />

    </div>
  );
}
