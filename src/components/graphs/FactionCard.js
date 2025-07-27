import React from "react";
import "./shared/FactionCard.css";
import VictoryPointBarChart from "../VictoryPointBarChart";

export default function FactionCard({ data }) {
  const {
    faction,
    iconUrl,
    avgPoints,
    winRate,
    totalPlays,
    playersPlayed,
    playersWon,
    vpHistogram,
  } = data;

  return (
    <div className="faction-card-container">
      <div className="faction-card-left">
        <img src={iconUrl} alt={faction} className="faction-icon" />
        <div className="faction-name">{faction}</div>
      </div>

      <div className="faction-card-middle">
        <div><strong>Win Rate:</strong> {winRate.toFixed(2)}% ({data.totalWins} of {totalPlays})</div>
        <div><strong>Avg Points:</strong> {avgPoints.toFixed(2)}</div>
        <div><strong>Players Played:</strong> {Object.entries(playersPlayed).map(([p, c]) => `${p} (${c})`).join(", ") || "-"}</div>
        <div><strong>Players Won With:</strong> {Object.entries(playersWon).map(([p, c]) => `${p} (${c})`).join(", ") || "-"}</div>
      </div>

      <div className="faction-card-right">
        <VictoryPointBarChart vpHistogram={vpHistogram} />
      </div>
    </div>
  );
}
