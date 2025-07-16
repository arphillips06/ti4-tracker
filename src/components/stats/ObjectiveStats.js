import React, { useState } from "react";
import ObjectiveStatsChart from "../graphs/ObjectiveStatsChart";
import ObjectiveFrequencyChart from "../graphs/ObjectiveFrequencyChart";
import ObjectiveAppearanceChart from "../graphs/ObjectiveAppearanceChart";
import '../../pages/stats.css';

export default function ObjectiveStats({ stats }) {
  const [showObjectiveData, setShowObjectiveData] = useState(false);

  return (
    <div className="stats-section">
      <h2 className="chart-header">Objective Statistics</h2>

      {/* Objective Points Chart */}
      <div className="chart-glass-container">
        <ObjectiveStatsChart stats={stats} />
        <button
          onClick={() => setShowObjectiveData(!showObjectiveData)}
          className="btn btn-sm btn-outline-secondary mt-2"
        >
          {showObjectiveData ? "Hide Raw Data" : "Show Raw Data"}
        </button>

        {showObjectiveData && (
          <ul className="raw-data-list mt-2 mb-2">
            <li>Public Objectives: {stats.objectiveStats.publicScored}</li>
            <li>Secret Objectives: {stats.objectiveStats.secretScored}</li>
            <li>Stage I: {stats.objectiveStats.stage1Scored}</li>
            <li>Stage II: {stats.objectiveStats.stage2Scored}</li>
          </ul>
        )}
      </div>

      {/* Frequency Chart */}
      <div className="chart-glass-container">
        <ObjectiveFrequencyChart
          frequency={stats.objectiveFrequency}
          secretPublic={stats.objectiveStats.secretPublic}
        />
      </div>

      {/* Appearance Chart */}
      <div className="chart-glass-container">
        <ObjectiveAppearanceChart stats={stats} />
      </div>
    </div>
  );

}
