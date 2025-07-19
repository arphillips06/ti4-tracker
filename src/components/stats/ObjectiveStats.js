import React, { useState } from "react";
import ObjectiveStatsChart from "../graphs/ObjectiveStatsChart";
import ObjectiveFrequencyChart from "../graphs/ObjectiveFrequencyChart";
import ObjectiveAppearanceChart from "../graphs/ObjectiveAppearanceChart";
import ObjectiveMetaTable from "../graphs/ObjectiveMetaTable";
import SecretObjectiveTable from "../graphs/SecretOnjectiveTable";
import '../../pages/stats.css';

export default function ObjectiveStats({ stats }) {
  const [showObjectiveData, setShowObjectiveData] = useState(false);

  return (
    <div className="stats-section">
      <h2 className="chart-header">Objective Statistics</h2>
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
      <div className="chart-glass-container">
        <ObjectiveFrequencyChart
          data={stats.objectiveAppearanceStats} />
      </div>
      <div className="chart-glass-container">
        <ObjectiveAppearanceChart data={stats.objectiveAppearanceStats} />
      </div>
      <div className="chart-glass-container">
        <ObjectiveMetaTable metas={stats.objectiveMetaStats} />
        <SecretObjectiveTable secrets={stats.objectiveMetaStats} />
      </div>

    </div>
  );

}
