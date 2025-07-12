import React, { useState } from "react";
import ObjectiveStatsChart from "../graphs/ObjectiveStatsChart";
import ObjectiveFrequencyChart from "../graphs/ObjectiveFrequencyChart";
import ObjectiveAppearanceChart from "../graphs/ObjectiveAppearanceChart";

export default function ObjectiveStats({ stats }) {
  const [showObjectiveData, setShowObjectiveData] = useState(false);

  return (
    <div className="stats-section">
      <h2>Objective Statistics</h2>

      <ObjectiveStatsChart stats={stats} />

      <button
        onClick={() => setShowObjectiveData(!showObjectiveData)}
        className="btn btn-sm btn-outline-secondary mb-2"
      >
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
      <ObjectiveAppearanceChart stats={stats} />
    </div>
  );
}
