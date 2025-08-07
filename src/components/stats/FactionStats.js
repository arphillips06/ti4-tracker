import React, { useState } from "react";
import FactionPlayWinChart from "../graphs/FactionPlayWinChart";
import FactionWinRateChart from "../graphs/FactionWinRateChart";
import '../../pages/stats.css';
import FactionCardList from "../graphs/FactionCardList";
import '../graphs/shared/modal.css';


// Stage classification sets
const STAGE1 = new Set([
  "Amass wealth", "Explore Deep Space", "Found Research Outposts", "Improve Infrastructure",
  "Intimidate Council", "Negotiate Trade Routes", "Push Boundaries", "Raise a Fleet",
  "Discover Lost Outposts", "Erect a Monument", "Expand Borders", "Diversify Research",
  "Lead from the Front", "Master the Sciences", "Populate the Outer Rim", "Protect the Border",
  "Revolutionize Warfare", "Rule Distant Lands", "Centralize Galactic Trade", "Corner the Market",
  "Form Galactic Brain Trust", "Make History", "Patrol Vast Territories", "Reclaim Ancient Monuments",
]);

const STAGE2 = new Set([
  "Achieve Supremacy", "Become a Legend", "Build Defenses", "Command an Armada",
  "Construct Massive Cities", "Control the Borderlands", "Found a Golden Age",
  "Galvanize the People", "Manipulate Galactic Law", "Subdue the Galaxy",
  "Unify the Colonies",
]);

export default function FactionStats({ stats }) {
  const [showFactionData, setShowFactionData] = useState(false);
  const [selectedFaction, setSelectedFaction] = useState(null);
  const [tab, setTab] = useState("stage1");

  if (!stats.factionPlayWinDistribution) {
    return <div>No faction data available.</div>;
  }

  const factionData = stats.factionPlayerStats || [];
  const aggregateStats = stats.factionAggregateStats || [];

  const handleMoreStatsClick = (faction) => {
    console.log(`Opening More Stats for ${faction}`);
    setSelectedFaction(faction);
    setTab("stage1");
  };

  const closeModal = () => setSelectedFaction(null);

  const factionObjStats = stats.factionObjectiveStats || {};
  const current = selectedFaction ? factionObjStats[selectedFaction] : null;

  const filteredEntries = () => {
    if (!current) return [];
    return Object.entries(current)
      .filter(([name]) => name && name.trim().length > 0)
      .filter(([_, obj]) => {
        const t = (obj.type || "").toLowerCase();
        if (tab === "secret") return t === "secret";
        if (t !== "public") return false;
        if (tab === "stage1") return STAGE1.has(_);
        if (tab === "stage2") return STAGE2.has(_);
        return false;
      })
      .sort((a, b) => (b[1].scoredCount || 0) - (a[1].scoredCount || 0));
  };

  return (
    <>
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
        <FactionCardList
          data={factionData}
          aggregates={aggregateStats}
          factionObjectiveStats={stats.factionObjectiveStats} // ✅ pass here
          onMoreStatsClick={handleMoreStatsClick}
        />
      </div>

      {/* Modal */}
      {selectedFaction && (
        <div className="stats-modal-overlay">
          <div className="stats-modal">
            <div className="d-flex justify-content-between align-items-center mb-2">
              <h4 className="m-0">{selectedFaction} — Objectives</h4>
              <button className="btn btn-sm btn-outline-secondary" onClick={closeModal}>
                Close
              </button>
            </div>

            <div className="btn-group mb-3">
              <button className={`btn btn-sm ${tab === "stage1" ? "btn-primary" : "btn-outline-primary"}`} onClick={() => setTab("stage1")}>Stage I</button>
              <button className={`btn btn-sm ${tab === "stage2" ? "btn-primary" : "btn-outline-primary"}`} onClick={() => setTab("stage2")}>Stage II</button>
              <button className={`btn btn-sm ${tab === "secret" ? "btn-primary" : "btn-outline-primary"}`} onClick={() => setTab("secret")}>Secrets</button>
            </div>

            {!current ? (
              <div className="text-muted">No objective stats for this faction yet.</div>
            ) : (
              <div className="overflow-auto" style={{ maxHeight: 420 }}>
                {tab === "secret" ? (
                  <table className="table table-sm table-dark table-striped">
                    <thead>
                      <tr>
                        <th>Objective</th>
                        <th className="text-end">Scored</th>
                      </tr>
                    </thead>
                    <tbody>
                      {filteredEntries().map(([name, obj]) => (
                        <tr key={name}>
                          <td>{name}</td>
                          <td className="text-end">{obj.scoredCount ?? 0}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                ) : (
                  <table className="table table-sm table-dark table-striped">
                    <thead>
                      <tr>
                        <th>Objective</th>
                        <th className="text-end">Appeared</th>
                        <th className="text-end">Scored</th>
                        <th className="text-end">% When Appeared</th>
                      </tr>
                    </thead>
                    <tbody>
                      {filteredEntries().map(([name, obj]) => (
                        <tr key={name}>
                          <td>{name}</td>
                          <td className="text-end">{obj.appearedCount ?? 0}</td>
                          <td className="text-end">{obj.scoredCount ?? 0}</td>
                          <td className="text-end">
                            {((obj.scoredWhenAppearedRate ?? 0) * 100).toFixed(1)}%
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </>
  );
}
