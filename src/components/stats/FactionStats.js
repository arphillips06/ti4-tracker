import React, { useState } from "react";
import FactionPlayWinChart from "../graphs/FactionPlayWinChart";
import FactionWinRateChart from "../graphs/FactionWinRateChart";


export default function FactionStats({ stats }) {
    const [showRaw, setShowRaw] = useState(false);
    const [showFactionData, setShowFactionData] = useState(false);
    if (!stats.factionPlayWinDistribution) return <div>No faction data available.</div>;


    const entries = Object.entries(stats.factionPlayWinDistribution);
    const sorted = [...entries].sort((a, b) => b[1].winRate - a[1].winRate); // sort by win rate descending

    return (
        <div className="stats-section">
            <h3>Faction Win Rates</h3>
            <FactionWinRateChart dataMap={stats.winRateByFaction} />

            <button
                onClick={() => setShowFactionData(!showFactionData)}
                className="btn btn-sm btn-outline-secondary mt-2"
            >
                {showFactionData ? "Hide Raw Data" : "Show Raw Data"}
            </button>

            {showFactionData && (
                <ul className="mt-2 mb-3">
                    {Object.entries(stats.winRateByFaction).map(([faction, rate]) => (
                        <li key={faction}>
                            {faction}: {rate.toFixed(2)}%
                        </li>
                    ))}
                </ul>
            )}



            <h3>Faction Play vs Win Rates</h3>
            <FactionPlayWinChart data={stats.factionPlayWinDistribution} />


        </div>

    );
}
