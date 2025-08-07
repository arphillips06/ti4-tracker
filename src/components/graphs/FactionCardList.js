import React from "react";
import "./shared/FactionCard.css";
import FactionCard from "./FactionCard";
import factionImageMap from "../../data/factionIcons";

export default function FactionCardList({
  data,
  aggregates = [],
  factionObjectiveStats = {}, // ✅ added here
  onMoreStatsClick
}) {
  const aggregateMap = {};
  aggregates.forEach((agg) => {
    aggregateMap[agg.faction.toLowerCase()] = agg;
  });

  const grouped = {};
  data.forEach(({ faction, player, playedCount, wonCount }) => {
    if (!grouped[faction]) {
      grouped[faction] = {
        faction,
        totalPlays: 0,
        totalWins: 0,
        playersPlayed: {},
        playersWon: {},
      };
    }
    grouped[faction].totalPlays += playedCount;
    grouped[faction].totalWins += wonCount;
    if (playedCount > 0)
      grouped[faction].playersPlayed[player] =
        (grouped[faction].playersPlayed[player] || 0) + playedCount;
    if (wonCount > 0)
      grouped[faction].playersWon[player] =
        (grouped[faction].playersWon[player] || 0) + wonCount;
  });

  const tableData = Object.values(grouped).map((row) => {
    const iconUrl = `/faction-icons/${factionImageMap[row.faction] || "default.webp"}`;
    const agg = aggregateMap[row.faction.toLowerCase()];
    const avgPoints =
      agg && agg.totalPlays > 0 ? agg.totalPointsScored / agg.totalPlays : 0;

    return {
      ...row,
      avgPoints,
      winRate: row.totalPlays > 0 ? (row.totalWins / row.totalPlays) * 100 : 0,
      iconUrl,
      vpHistogram: (agg && agg.vpHistogram) || [],
      objectiveStats: factionObjectiveStats?.[row.faction] || {} // ✅ now passed down
    };
  });

  return (
    <div className="faction-card-grid">
      {tableData.map((f) => (
        <FactionCard
          key={f.faction}
          data={f}
          onMoreStatsClick={onMoreStatsClick} // ✅ pass callback
        />
      ))}
    </div>
  );
}
