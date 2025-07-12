import React, { useState } from "react";
import FactionWinRateChart from "../graphs/FactionWinRateChart";

export default function OverviewStats({ stats }) {
  const averageRounds = stats.averageGameRounds?.toFixed(1) || "—";
const averagePoints = stats.averagePlayerPoints?.toFixed(2) || "—";
const mostPlayedFaction = stats.mostPlayedFaction || "—";
const mostVictoriousFaction = stats.mostVictoriousFaction || "—";
const totalPlayers = stats.totalUniquePlayers || "—";

  return (
<div className="d-flex flex-wrap gap-4 mb-4">
  <div className="border rounded p-3 text-center">
    <h5 className="mb-1">Total Games</h5>
    <div className="fs-4">{stats.totalGames}</div>
  </div>
  <div className="border rounded p-3 text-center">
    <h5 className="mb-1">Average Rounds per Game</h5>
    <div className="fs-4">{averageRounds}</div>
  </div>
  <div className="border rounded p-3 text-center">
    <h5 className="mb-1">Average Player Score</h5>
    <div className="fs-4">{averagePoints}</div>
  </div>
  <div className="border rounded p-3 text-center">
    <h5 className="mb-1">Unique Players</h5>
    <div className="fs-4">{totalPlayers}</div>
  </div>
  <div className="border rounded p-3 text-center">
    <h5 className="mb-1">Most Played Faction</h5>
    <div className="fs-5">{mostPlayedFaction}</div>
  </div>
  <div className="border rounded p-3 text-center">
    <h5 className="mb-1">Most Victorious Faction</h5>
    <div className="fs-5">{mostVictoriousFaction}</div>
  </div>
</div>

  );
}
